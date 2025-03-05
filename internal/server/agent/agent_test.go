package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Result struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}

func TestComputeAddition(t *testing.T) {
	task := Task{
		Arg1:          2,
		Arg2:          3,
		Operation:     "+",
		OperationTime: 0,
	}
	result := compute(task)
	if result != 5 {
		t.Errorf("expected 5, got %v", result)
	}
}

func TestComputeSubtraction(t *testing.T) {
	task := Task{
		Arg1:          5,
		Arg2:          2,
		Operation:     "-",
		OperationTime: 0,
	}
	result := compute(task)
	if result != 3 {
		t.Errorf("expected 3, got %v", result)
	}
}

func TestComputeMultiplication(t *testing.T) {
	task := Task{
		Arg1:          4,
		Arg2:          3,
		Operation:     "*",
		OperationTime: 0,
	}
	result := compute(task)
	if result != 12 {
		t.Errorf("expected 12, got %v", result)
	}
}

func TestComputeDivision(t *testing.T) {
	task := Task{
		Arg1:          6,
		Arg2:          2,
		Operation:     "/",
		OperationTime: 0,
	}
	result := compute(task)
	if result != 3 {
		t.Errorf("expected 3, got %v", result)
	}
}

func TestComputeUnknownOperation(t *testing.T) {
	task := Task{
		Arg1:          2,
		Arg2:          3,
		Operation:     "^",
		OperationTime: 0,
	}
	result := compute(task)
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestComputeWithDelay(t *testing.T) {
	task := Task{
		Arg1:          2,
		Arg2:          3,
		Operation:     "+",
		OperationTime: 100,
	}
	start := time.Now()
	result := compute(task)
	duration := time.Since(start)
	if result != 5 {
		t.Errorf("expected 5, got %v", result)
	}
	if duration < 100*time.Millisecond || duration > 150*time.Millisecond {
		t.Errorf("expected delay around 100ms, got %v", duration)
	}
}

func TestWorkerIntegration(t *testing.T) {
	tasks := make(chan Task, 1)
	results := make(chan Result, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/internal/task" {
			select {
			case task := <-tasks:
				json.NewEncoder(w).Encode(map[string]Task{"task": task})
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		} else if r.Method == "POST" && r.URL.Path == "/internal/task" {
			var result Result
			if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
				t.Errorf("Failed to decode result: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			results <- result
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	go func() {
		client := &http.Client{}
		for {
			resp, err := client.Get(server.URL + "/internal/task")
			if err != nil || resp.StatusCode == http.StatusNotFound {
				time.Sleep(10 * time.Millisecond) // Уменьшено для теста
				continue
			}

			var data struct {
				Task Task `json:"task"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				t.Error("Error decoding task:", err)
				resp.Body.Close()
				continue
			}
			resp.Body.Close()

			resultVal := compute(data.Task)

			reqBody, _ := json.Marshal(Result{
				ID:     data.Task.ID,
				Result: resultVal,
			})
			resp, err = client.Post(server.URL+"/internal/task", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Error("Error sending result:", err)
				continue
			}
			resp.Body.Close()
		}
	}()
	task := Task{
		ID:            1,
		Arg1:          2,
		Arg2:          3,
		Operation:     "+",
		OperationTime: 50,
	}
	tasks <- task
	select {
	case result := <-results:
		if result.ID != 1 || result.Result != 5 {
			t.Errorf("expected ID=1 and Result=5, got ID=%d and Result=%v", result.ID, result.Result)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for worker result")
	}
}
