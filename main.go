package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
)

const (
	rpcUser     = "rpc"
	rpcPass     = "rpc"
	electrumURL = "127.0.0.1:50001"
	coreURL     = "http://localhost"
	walletURL   = "/wallet/bank"
	nmcPort     = 18443
	btcPort     = 18444
)

var (
	nmcParams chaincfg.Params
	btcParms  chaincfg.Params
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

func exampleElectrum() {
	scriptHash, _ := ElectrumScripthash("mqC6EWespCSjGPXZtz8VCxRSNtrep7FJDA", &nmcParams)
	params := []any{scriptHash}
	reqJSON := createElectrumRequest("blockchain.scripthash.get_balance", params)
	fmt.Println(reqJSON)
	fmt.Println()
	sendElectrumRequest(reqJSON)
}

func main() {

	height, err := getBlockHeight(nmcPort)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("height: ", height)

	http.HandleFunc("/template", templateEndpoint)
	http.HandleFunc("/nmc/loadHomePage", nmcLoadHomeReq)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func loadHome(coin string) {
	// Get BlockCount

	// Get Trends

	// Get 10 Latest Blocks

}

func getBlockHeight(portNum int) (int, error) {

	method := "getblockcount"
	params := []interface{}{}

	result, err := makeRPCRequest(method, params, nmcPort)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	if val, ok := result.(float64); ok {
		return int(val), nil
	} else {
		return 0, fmt.Errorf("int conversion failed")
	}
}

func nmcLoadHomeReq(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// // Read the request body
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Error reading request body", http.StatusBadRequest)
	// 	return
	// }

	// // Define a struct to unmarshal the JSON data
	// var req struct {
	// 	//struct fields here
	// }

	// // Unmarshal the JSON data
	// err = json.Unmarshal(body, &req)
	// if err != nil {
	// 	http.Error(w, "Error unmarshaling JSON data", http.StatusBadRequest)
	// 	return
	// }

	// response, err := loadHome("nmc")
	//================================================================================//
	//============================== Code Goes Here ==================================//
	//================================================================================//

	//================================================================================//
	//================================================================================//
	//================================================================================//

	// type res struct {
	// 	//struct fields here
	// }

	// // response := res{
	// // 	//fill fields
	// // }
	// // // Marshal the struct into JSON
	// resJSON, err := json.Marshal(response)
	// if err != nil {
	// 	http.Error(w, "Error marshaling data", http.StatusInternalServerError)
	// 	return
	// }

	// // Set headers and write JSON to response
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(resJSON)
}

// method := "getblockhash"
// params := []interface{}{250, []interface{}{"minfeerate", "avgfeerate"}}

// method := "getblockhash"
// params := []interface{}{250}

// result, err := makeRPCRequest(method, params)
// if err != nil {
// 	fmt.Println("Error:", err)
// 	return
// }

// fmt.Println("Result:", result)
// Sends a RPC request to the local core
func makeRPCRequest(method string, params []interface{}, portNum int) (interface{}, error) {
	// Set the RPC credentials
	username := "rpc"
	password := "rpc"

	// Create the RPC request JSON
	requestBody := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "curltest",
		"method":  method,
		"params":  params,
	}

	// Marshal the request JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Encode the credentials for Basic Authentication
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", coreURL+":"+fmt.Sprint(portNum), bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", "Basic "+auth)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response JSON
	var response map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	// Check for RPC error
	if response["error"] != nil {
		return nil, fmt.Errorf("RPC error: %v", response["error"])
	}

	// Extract the result
	result := response["result"]

	return result, nil
}
