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

func exampleElectrum() {
	scriptHash, _ := ElectrumScripthash("mqC6EWespCSjGPXZtz8VCxRSNtrep7FJDA", &nmcParams)
	params := []any{scriptHash}
	reqJSON := createElectrumRequest("blockchain.scripthash.get_balance", params)
	fmt.Println(reqJSON)
	fmt.Println()
	sendElectrumRequest(reqJSON)
}

func createElectrumRequest(method string, params []interface{}) string {
	// Create a map for the JSON-RPC request
	request := map[string]interface{}{
		"id":      3,
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	// Marshal the map into a JSON string
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error creating JSON-RPC request:", err)
		return ""
	}

	return string(jsonRequest) + "\n" // Ensure the request ends with a newline character

}

func sendElectrumRequest(jsonRequest string) {
	// Connect to the server
	conn, err := net.Dial("tcp", electrumURL)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	// Send the JSON-RPC request to the server
	_, err = fmt.Fprintf(conn, jsonRequest)
	if err != nil {
		fmt.Println("Error sending JSON-RPC request:", err)
		return
	}

	// Read the server's response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return
	}

	fmt.Println("Server Response:", response)
}

// Based on the following discussion of electrum scripthashing:
// https://github.com/bitcoinjs/bitcoinjs-lib/issues/990
func ElectrumScripthash(addressStr string, chainParams *chaincfg.Params) (string, error) {
	addr, err := btcutil.DecodeAddress(addressStr, chainParams)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	sum := sha256.Sum256(script)
	length := len(sum)
	for i := 0; i < length/2; i++ {
		// Swap arr[i] with arr[length-i-1]
		sum[i], sum[length-i-1] = sum[length-i-1], sum[i]
	}
	return hex.EncodeToString(sum[:]), nil
}
