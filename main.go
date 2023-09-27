package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// postHandler is a dedicated function to handle POST requests to "/post".
func templateEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Define a struct to unmarshal the JSON data
	var req struct {
		//struct fields here
	}

	// Unmarshal the JSON data
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON data", http.StatusBadRequest)
		return
	}

	//================================================================================//
	//============================== Code Goes Here ==================================//
	//================================================================================//

	//================================================================================//
	//================================================================================//
	//================================================================================//

	type res struct {
		//struct fields here
	}

	response := res{
		//fill fields
	}
	// // Marshal the struct into JSON
	resJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling data", http.StatusInternalServerError)
		return
	}

	// Set headers and write JSON to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func main() {
	http.HandleFunc("/template", templateEndpoint)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
