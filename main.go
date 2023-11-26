package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
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
	nmcParams chaincfg.Params = nmcMainnetChainParams
	btcParams chaincfg.Params = chaincfg.MainNetParams
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

	// Create a new router
	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/template", templateEndpoint)
	router.HandleFunc("/nmc/loadhomepage", nmcLoadHomeReq)
	router.HandleFunc("/nmc/address", nmcAddressReq)
	router.HandleFunc("/nmc/block", nmcBlockReq)
	router.HandleFunc("/nmc/tx", nmcTxReq)
	router.HandleFunc("/nmc/name", nmcNameReq)

	// Set up a handler function to handle CORS headers
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers to allow any origin
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Check if the request is an OPTIONS request (preflight)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Continue processing the request
			next.ServeHTTP(w, r)
		})
	}

	// Use the CORS handler for all routes
	http.Handle("/", corsHandler(router))

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}

// postHandler is a dedicated function to handle POST requests to "/post".
func nmcNameReq(w http.ResponseWriter, r *http.Request) {
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
		Name string `json:"name"`
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
	getName(req.Name)
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

func getName(name string) {

	// Name_history
	method := "name_history"
	params := []interface{}{name} // verbosity = 2 includes all transactions in block

	result, err := makeRPCRequest(method, params, nmcPort)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the provided data to the MyStruct type
	var nameHistory []NameHistoryRes
	dataJSON, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// fmt.Println("=============================================================")
	// fmt.Println(string(dataJSON))
	// fmt.Println("=============================================================")
	err = json.Unmarshal(dataJSON, &nameHistory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	spew.Dump(nameHistory[0])
}

type NameHistoryRes struct {
	Address       string `json:"address"`
	Expired       bool   `json:"expired"`
	ExpiresIn     int    `json:"expires_in"`
	Height        int    `json:"height"`
	IsMine        bool   `json:"ismine"`
	Name          string `json:"name"`
	NameEncoding  string `json:"name_encoding"`
	TxID          string `json:"txid"`
	Value         string `json:"value"`
	ValueEncoding string `json:"value_encoding"`
	Vout          int    `json:"vout"`
}
