package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
)

const (
	rpcUser     = "rpc"
	rpcPass     = "rpc"
	electrumURL = "192.168.1.2:50001"
	coreURL     = "http://localhost"
	walletURL   = "/wallet/bank" // bank wallet for regtest use
	nmcPort     = 18443
	btcPort     = 18444
)

var (
	nmcParams chaincfg.Params
	btcParams chaincfg.Params
)

// TODO: replace the implementations pf this with an interface{}["result"] instead to save on repetitive ElectrumResponse structs
type ElectrumTransactionResponse struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  ElectrumTransaction `json:"result"`
	ID      int                 `json:"id"`
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

	// loadHome("nmc")
	// block, _ := getBlock("4ddbe4874f32ad83727a9dafbf177394d9da3e1311c361e5fb27aa139f2a2103", nmcPort)

	// spew.Dump(block.Tx)
	http.HandleFunc("/template", templateEndpoint)
	http.HandleFunc("/nmc/loadhomepage", nmcLoadHomeReq)
	http.HandleFunc("/nmc/address", nmcAddressReq)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

// postHandler is a dedicated function to handle POST requests to "/post".
func nmcAddressReq(w http.ResponseWriter, r *http.Request) {
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
		Address string `json:"address"` //struct fields here
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
	getAddress(req.Address, "nmc")
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

func getAddress(addr string, chain string) {
	var scriptHash string
	if chain == "nmc" {
		scriptHash, _ = ElectrumScripthash("mqC6EWespCSjGPXZtz8VCxRSNtrep7FJDA", &nmcParams)
	} else if chain == "btc" {
		scriptHash, _ = ElectrumScripthash("mqC6EWespCSjGPXZtz8VCxRSNtrep7FJDA", &chaincfg.MainNetParams)
	}
	params := []any{scriptHash}
	reqJSON := createElectrumRequest("blockchain.scripthash.get_balance", params)
	fmt.Println(reqJSON)
	fmt.Println()
	elecRes := sendElectrumRequest(reqJSON)
	fmt.Println(elecRes)
}
