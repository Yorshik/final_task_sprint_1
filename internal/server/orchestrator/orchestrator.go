package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Yorshik/final_task_sprint_1/internal/ast"
	"github.com/gorilla/mux"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Expression struct {
	ID     string    `json:"id"`
	Status string    `json:"status"`
	Result *float64  `json:"result"`
	Node   *ast.Node `json:"-"`
	Tasks  []Task    `json:"-"`
}

type Orchestrator struct {
	expressions map[string]*Expression
	tasks       chan Task
	results     map[int]float64
	taskID      int
	mu          sync.Mutex
	wg          sync.WaitGroup
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		expressions: make(map[string]*Expression),
		tasks:       make(chan Task, 100),
		results:     make(map[int]float64),
	}
}

func (o *Orchestrator) AddExpression(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	node, err := ast.Parse(req.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	id := strconv.Itoa(len(o.expressions) + 1)
	expr := &Expression{ID: id, Status: "pending", Node: node}
	o.expressions[id] = expr
	o.mu.Unlock()

	go o.processExpression(expr)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (o *Orchestrator) processExpression(expr *Expression) {
	result := o.evaluateNode(expr.Node, expr)

	o.mu.Lock()
	expr.Result = &result
	expr.Status = "completed"
	o.mu.Unlock()
}

func (o *Orchestrator) evaluateNode(node *ast.Node, expr *Expression) float64 {
	if node.Operator == "" {
		return node.Value
	}

	var leftVal, rightVal float64
	if node.Left.Operator != "" {
		leftVal = o.evaluateNode(node.Left, expr)
	} else {
		leftVal = node.Left.Value
	}
	if node.Right.Operator != "" {
		rightVal = o.evaluateNode(node.Right, expr)
	} else {
		rightVal = node.Right.Value
	}

	o.mu.Lock()
	o.taskID++
	task := Task{
		ID:            o.taskID,
		Arg1:          leftVal,
		Arg2:          rightVal,
		Operation:     node.Operator,
		OperationTime: o.getOperationTime(node.Operator),
	}
	expr.Tasks = append(expr.Tasks, task)
	o.mu.Unlock()

	o.tasks <- task

	for {
		o.mu.Lock()
		if result, ok := o.results[task.ID]; ok {
			delete(o.results, task.ID)
			o.mu.Unlock()
			return result
		}
		o.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (o *Orchestrator) getOperationTime(op string) int {
	switch op {
	case "+":
		return getEnvInt("TIME_ADDITION_MS", 1000)
	case "-":
		return getEnvInt("TIME_SUBTRACTION_MS", 1000)
	case "*":
		return getEnvInt("TIME_MULTIPLICATIONS_MS", 1000)
	case "/":
		return getEnvInt("TIME_DIVISIONS_MS", 1000)
	default:
		return 1000
	}
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func (o *Orchestrator) GetExpressions(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	defer o.mu.Unlock()

	resp := struct {
		Expressions []*Expression `json:"expressions"`
	}{Expressions: make([]*Expression, 0, len(o.expressions))}
	for _, expr := range o.expressions {
		resp.Expressions = append(resp.Expressions, expr)
	}
	json.NewEncoder(w).Encode(resp)
}

func (o *Orchestrator) GetExpression(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	o.mu.Lock()
	defer o.mu.Unlock()

	expr, ok := o.expressions[id]
	if !ok {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]*Expression{"expression": expr})
}

func (o *Orchestrator) GetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-o.tasks:
		json.NewEncoder(w).Encode(map[string]Task{"task": task})
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func (o *Orchestrator) ReceiveResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	o.results[req.ID] = req.Result
	o.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (o *Orchestrator) Web(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func StartServer() {
	o := NewOrchestrator()
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", o.AddExpression).Methods("POST")
	r.HandleFunc("/api/v1/expressions", o.GetExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", o.GetExpression).Methods("GET")
	r.HandleFunc("/internal/task", o.GetTask).Methods("GET")
	r.HandleFunc("/internal/task", o.ReceiveResult).Methods("POST")
	r.HandleFunc("/", o.Web).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
