package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestApiCalcHandler тестирует эндпоинт /api/v1/calculate для различных методов и данных.
func TestApiCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{} // string for GET/no body, struct or string for POST
		contentType    string
		expectedStatus int
		expectedBody   interface{} // string for GET/error messages, map for POST success
	}{
		// 1. POST-запрос с корректным JSON
		{
			name:           "POST_Valid_JSON",
			method:         http.MethodPost,
			body:           CalcRequest{Expression: "2+2*2"},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]float64{"result": 6},
		},
		// 2. POST-запрос с некорректным JSON (невалидный формат)
		{
			name:           "POST_Invalid_JSON_Format",
			method:         http.MethodPost,
			body:           "{expression:2+2*2}",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Incorrect JSON format\n",
		},
		// 3. POST-запрос с невалидным выражением
		{
			name:           "POST_Invalid_Expression",
			method:         http.MethodPost,
			body:           CalcRequest{Expression: "2++2"},
			contentType:    "application/json",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": "Expression is invalid"},
		},
		// 4. POST-запрос с незаполненным полем
		{
			name:           "POST_Empty_Expression",
			method:         http.MethodPost,
			body:           CalcRequest{Expression: ""},
			contentType:    "application/json",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": "Expression is invalid"},
		},
		// 8. POST-запрос с неверными ключами в JSON
		{
			name:           "POST_Invalid_JSON_Keys",
			method:         http.MethodPost,
			body:           `{"expr":"2+2*2"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   map[string]string{"error": "Expression is invalid"},
		},
		// 9. POST-запрос с пустым телом
		{
			name:           "POST_Empty_Body",
			method:         http.MethodPost,
			body:           "",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Incorrect JSON format\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			switch v := tt.body.(type) {
			case string:
				reqBody = bytes.NewBufferString(v)
			case CalcRequest:
				jsonData, err := json.Marshal(v)
				if err != nil {
					t.Fatalf("Failed to marshal body: %v", err)
				}
				reqBody = bytes.NewBuffer(jsonData)
			default:
				reqBody = nil
			}

			req := httptest.NewRequest(tt.method, "/api/v1/calculate/", reqBody)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rr := httptest.NewRecorder()
			ApiCalcHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			responseBody := rr.Body.String()

			switch expected := tt.expectedBody.(type) {
			case string:
				if responseBody != expected {
					t.Errorf("expected body '%s', got '%s'", expected, responseBody)
				}
			case map[string]float64:
				var response map[string]float64
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
					return
				}
				for key, val := range expected {
					if response[key] != val {
						t.Errorf("for key '%s', expected %v, got %v", key, val, response[key])
					}
				}
			}
		})
	}
}
