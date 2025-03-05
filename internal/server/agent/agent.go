package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func compute(task Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}

func worker() {
	client := &http.Client{}
	for {
		resp, err := client.Get("http://localhost:8080/internal/task")
		if err != nil || resp.StatusCode == http.StatusNotFound {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		var data struct {
			Task Task `json:"task"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		result := compute(data.Task)

		reqBody, err := json.Marshal(map[string]interface{}{
			"id":     data.Task.ID,
			"result": result,
		})
		if err != nil {
			continue
		}

		resp, err = client.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			continue
		}
		resp.Body.Close()
	}
}

func StartAgent() {
	power, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if power <= 0 {
		power = 1
	}

	for i := 0; i < power; i++ {
		go worker()
	}

	select {}
}
