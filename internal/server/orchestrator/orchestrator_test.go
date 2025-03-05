package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Yorshik/final_task_sprint_1/internal/ast"
	"github.com/gorilla/mux"
)

func setupEnv() {
	os.Setenv("TIME_ADDITION_MS", "100")
	os.Setenv("TIME_SUBTRACTION_MS", "100")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "100")
	os.Setenv("TIME_DIVISIONS_MS", "100")
}

func TestAddExpression(t *testing.T) {
	setupEnv()
	o := NewOrchestrator()
	r := httptest.NewServer(http.HandlerFunc(o.AddExpression))
	defer r.Close()

	tests := []struct {
		name         string
		expression   string
		expectedCode int
	}{
		{"ValidExpression", "2 + 3", http.StatusCreated},
		{"InvalidExpression", "2 + * 3", http.StatusUnprocessableEntity},
		{"EmptyExpression", "", http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(map[string]string{"expression": tt.expression})
			resp, err := http.Post(r.URL, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}
			if tt.expectedCode == http.StatusCreated {
				var respData struct {
					ID string `json:"id"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if respData.ID == "" {
					t.Error("Expected non-empty ID")
				}
			}
		})
	}
}

func TestGetExpressions(t *testing.T) {
	setupEnv()
	o := NewOrchestrator()
	node, _ := ast.Parse("2 + 3")
	o.mu.Lock()
	o.expressions["1"] = &Expression{ID: "1", Status: "pending", Node: node}
	o.mu.Unlock()
	r := httptest.NewServer(http.HandlerFunc(o.GetExpressions))
	defer r.Close()
	resp, err := http.Get(r.URL)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	var respData struct {
		Expressions []*Expression `json:"expressions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(respData.Expressions) != 1 {
		t.Errorf("Expected 1 expression, got %d", len(respData.Expressions))
	}
	if respData.Expressions[0].ID != "1" {
		t.Errorf("Expected ID '1', got %s", respData.Expressions[0].ID)
	}
}

func TestGetExpression(t *testing.T) {
	setupEnv()
	o := NewOrchestrator()
	node, _ := ast.Parse("2 + 3")
	result := 5.0
	o.mu.Lock()
	o.expressions["1"] = &Expression{ID: "1", Status: "completed", Result: &result, Node: node}
	o.mu.Unlock()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/expressions/{id}", o.GetExpression).Methods("GET")
	s := httptest.NewServer(r)
	defer s.Close()
	tests := []struct {
		name         string
		id           string
		expectedCode int
		expectedID   string
	}{
		{"ExistingExpression", "1", http.StatusOK, "1"},
		{"NonExistingExpression", "2", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(s.URL + "/api/v1/expressions/" + tt.id)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}
			if tt.expectedCode == http.StatusOK {
				var respData struct {
					Expression *Expression `json:"expression"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if respData.Expression.ID != tt.expectedID {
					t.Errorf("Expected ID %s, got %s", tt.expectedID, respData.Expression.ID)
				}
				if *respData.Expression.Result != 5 {
					t.Errorf("Expected result 5, got %v", *respData.Expression.Result)
				}
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	setupEnv()
	o := NewOrchestrator()
	task := Task{ID: 1, Arg1: 2, Arg2: 3, Operation: "+", OperationTime: 100}
	o.tasks <- task
	r := httptest.NewServer(http.HandlerFunc(o.GetTask))
	defer r.Close()
	resp, err := http.Get(r.URL)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	var respData struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if respData.Task.ID != 1 || respData.Task.Arg1 != 2 || respData.Task.Arg2 != 3 {
		t.Errorf("Expected task ID=1, Arg1=2, Arg2=3, got %+v", respData.Task)
	}
}

func TestReceiveResult(t *testing.T) {
	setupEnv()
	o := NewOrchestrator()
	r := httptest.NewServer(http.HandlerFunc(o.ReceiveResult))
	defer r.Close()
	reqBody, _ := json.Marshal(struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}{ID: 1, Result: 5})
	resp, err := http.Post(r.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	o.mu.Lock()
	result, ok := o.results[1]
	o.mu.Unlock()
	if !ok || result != 5 {
		t.Errorf("Expected result 5 for ID 1, got %v (found: %v)", result, ok)
	}
}

func TestFullWorkflow(t *testing.T) {
	setupEnv()
	o := NewOrchestrator()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", o.AddExpression).Methods("POST")
	r.HandleFunc("/api/v1/expressions/{id}", o.GetExpression).Methods("GET")
	r.HandleFunc("/internal/task", o.GetTask).Methods("GET")
	r.HandleFunc("/internal/task", o.ReceiveResult).Methods("POST")
	s := httptest.NewServer(r)
	defer s.Close()
	reqBody, _ := json.Marshal(map[string]string{"expression": "2 + 3"})
	resp, err := http.Post(s.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send calculate request: %v", err)
	}
	defer resp.Body.Close()
	var calcResp struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&calcResp); err != nil {
		t.Fatalf("Failed to decode calculate response: %v", err)
	}
	resp, err = http.Get(s.URL + "/internal/task")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	defer resp.Body.Close()

	var taskResp struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		t.Fatalf("Failed to decode task: %v", err)
	}
	result := taskResp.Task.Arg1 + taskResp.Task.Arg2
	reqBody, _ = json.Marshal(struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}{ID: taskResp.Task.ID, Result: result})
	resp, err = http.Post(s.URL+"/internal/task", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send result: %v", err)
	}
	defer resp.Body.Close()
	timeout := time.After(1 * time.Second)
	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for expression to complete")
		case <-ticker:
			resp, err := http.Get(s.URL + "/api/v1/expressions/" + calcResp.ID)
			if err != nil {
				t.Fatalf("Failed to get expression: %v", err)
			}
			defer resp.Body.Close()

			var exprResp struct {
				Expression *Expression `json:"expression"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&exprResp); err != nil {
				t.Fatalf("Failed to decode expression: %v", err)
			}

			if exprResp.Expression.Status == "completed" {
				if *exprResp.Expression.Result != 5 {
					t.Errorf("Expected result 5, got %v", *exprResp.Expression.Result)
				}
				return
			}
		}
	}
}
