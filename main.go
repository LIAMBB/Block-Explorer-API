package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/davecgh/go-spew/spew"
)

const (
	rpcUser     = "rpc"
	rpcPass     = "rpc"
	electrumURL = "127.0.0.1:50001"
	coreURL     = "http://localhost"
	walletURL   = "/wallet/bank" // bank wallet for regtest use
	nmcPort     = 18443
	btcPort     = 18444
)

var (
	nmcParams chaincfg.Params
	btcParms  chaincfg.Params
)

type AuxPow struct {
	Tx                map[string]interface{} `json:"tx"`
	ChainIndex        float64                `json:"chainindex"`
	MerkleBranch      []interface{}          `json:"merklebranch"`
	ChainMerkleBranch []interface{}          `json:"chainmerklebranch"`
	ParentBlock       struct {
		Version    float64 `json:"version"`
		VersionHex string  `json:"versionHex"`
		MerkleRoot string  `json:"merkleroot"`
		Time       float64 `json:"time"`
		Nonce      float64 `json:"nonce"`
		Bits       string  `json:"bits"`
		Difficulty float64 `json:"difficulty"`
		Hash       string  `json:"hash"`
	} `json:"parentblock"`
}

type Block struct {
	Version           float64  `json:"version"`
	MedianTime        float64  `json:"mediantime"`
	ChainWork         string   `json:"chainwork"`
	Tx                []string `json:"tx"`
	Weight            float64  `json:"weight"`
	MerkleRoot        string   `json:"merkleroot"`
	Time              float64  `json:"time"`
	Nonce             float64  `json:"nonce"`
	Bits              string   `json:"bits"`
	Confirmations     float64  `json:"confirmations"`
	NextBlockHash     string   `json:"nextblockhash"`
	StrippedSize      float64  `json:"strippedsize"`
	Hash              string   `json:"hash"`
	VersionHex        string   `json:"versionHex"`
	PreviousBlockHash string   `json:"previousblockhash"`
	Difficulty        float64  `json:"difficulty"`
	Height            float64  `json:"height"`
	NTx               float64  `json:"nTx"`
	Size              float64  `json:"size"`
	AuxPow            AuxPow   `json:"auxpow"`
}

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

	loadHome("nmc")

	http.HandleFunc("/template", templateEndpoint)
	http.HandleFunc("/nmc/loadHomePage", nmcLoadHomeReq)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func loadHome(coin string) {
	port := 0
	if coin == "nmc" {
		port = nmcPort
	} else if coin == "btc" {
		port = btcPort
	} else {
		return // Error
	}

	// Get BlockCount
	blockHeight, err := getBlockHeight(port)

	if err != nil {
		fmt.Println("Error getting current blockheight")
	}

	fmt.Println("Blockheight: ", blockHeight)

	// Get 10 Latest Blocks

	for i := 0; i < 10; i++ {
		blockHash, _ := getBlockHash((blockHeight - i), nmcPort)
		block, _ := getBlock(blockHash, nmcPort)
		spew.Dump(block)
	}

	// Get Trends

}

func getBlock(hash string, portNum int) (interface{}, error) {

	method := "getblock"
	params := []interface{}{hash}

	result, err := makeRPCRequest(method, params, nmcPort)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	// Convert the provided data to the MyStruct type
	var myStruct Block
	dataJSON, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	err = json.Unmarshal(dataJSON, &myStruct)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return myStruct, nil
}

func getBlockHash(height int, portNum int) (string, error) {

	method := "getblockhash"
	params := []interface{}{height}

	result, err := makeRPCRequest(method, params, nmcPort)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return fmt.Sprint(result), nil
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
