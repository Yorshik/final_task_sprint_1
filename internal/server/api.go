package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yorshik/final_task_sprint_1/internal/calc"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func ApiCalcHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandleApiCalcGet(w, r)
	case http.MethodPost:
		HandleApiCalcPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleApiCalcPost(w http.ResponseWriter, r *http.Request) {
	var req CalcRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Incorrect JSON format", http.StatusBadRequest)
		return
	}
	var res CalcResponse
	defer func() {
		if rec := recover(); rec != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res = CalcResponse{
				Error: "Internal server error",
			}
			json.NewEncoder(w).Encode(res)
		}
	}()
	calc_result, calc_err := calc.Calc(req.Expression)
	if calc_err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		res = CalcResponse{
			Error: "Expression is not valid",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	res = CalcResponse{
		Result: calc_result,
	}
	json.NewEncoder(w).Encode(res)
}

func HandleApiCalcGet(w http.ResponseWriter, r *http.Request) {
	message := `Allowed only POST requests.\n
Parameters: 'expression': string\n
Response:\n
200 - 'result': float\n
422 - 'error': 'Expression is not valid'\n
500 - 'error': 'internal server error'
	`
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, message)
}
