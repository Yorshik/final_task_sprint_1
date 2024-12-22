package server

import (
	"fmt"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/api/v1/calculate/", ApiCalcHandler)
	fmt.Println("Listening on port 9000...")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
