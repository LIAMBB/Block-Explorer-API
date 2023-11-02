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
	electrumURL = "127.0.0.1:50001" // for localhost testing Electrum
	// electrumURL = "192.168.1.2:50001" // NMC Endpoint on docker image for deployment
	coreURL   = "http://localhost"
	walletURL = "/wallet/bank" // bank wallet for regtest use
	nmcPort   = 18443
	btcPort   = 18444
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
	nmcParams = nmcMainnetChainParams
	// spew.Dump(block.Tx)
	http.HandleFunc("/template", templateEndpoint)
	http.HandleFunc("/nmc/loadhomepage", nmcLoadHomeReq)
	http.HandleFunc("/nmc/address", nmcAddressReq)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
	r := mux.NewRouter()

	// Apply CORS middleware globally for all routes
	r.Use(corsMiddleware)

	// Define your endpoints and handlers
	r.HandleFunc("/nmc/loadhomepage", nmcLoadHomeReq)

	// Start the server
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow access from anywhere
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle the preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
	spew.Dump(req)
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
		var err error
		scriptHash, err = ElectrumScripthash(addr, &nmcParams)
		fmt.Println("Error: ", err)
	} else if chain == "btc" {
		scriptHash, _ = ElectrumScripthash(addr, &chaincfg.MainNetParams)
	}
	fmt.Println(chain, ": ", scriptHash)

	histTxs := getAddressHist(scriptHash)
	addrBal := getAddressBal(scriptHash)
	spew.Dump(histTxs, addrBal)
	// getFullHistTx()
}

type FullHistTransaction struct {
	TxID          string
	Confirmations int
	Size          int
	Vin           []FullVin
	Vout          []FullVout
}

type FullVin struct {
	TxID    string
	Amount  float64
	Index   int
	Address string
}
type FullVout struct {
	Amount  float64
	Index   int
	Address string
}

func getFullHistTx(histTx HistoryTransaction) {
	tx, _ := getTx(histTx.TxHash, nmcPort)

	// Get Txs for all Vins
	// Extract their amount, address, txid and index and add to struct as FullVin
	// Go through al vouts and populate the FullVout array

	// calculate confirmations
	// get size
	// add Txid

}

type Response struct {
	JSONRPC string               `json:"jsonrpc"`
	Result  []HistoryTransaction `json:"result"`
	ID      int                  `json:"id"`
}

type HistoryTransaction struct {
	TxHash string `json:"tx_hash"`
	Height int    `json:"height"`
}

func getAddressHist(scriptHash string) []HistoryTransaction {
	params := []any{scriptHash}
	reqJSON := createElectrumRequest("blockchain.scripthash.get_history", params)
	elecRes := sendElectrumRequest(reqJSON)
	var response Response

	// Unmarshal JSON data into the struct
	if err := json.Unmarshal([]byte(elecRes), &response); err != nil {
		fmt.Println("Error:", err)
		return []HistoryTransaction{}
	}
	return response.Result
}

type BalanceResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	Result  AddrBal `json:"result"`
	ID      int     `json:"id"`
}

type AddrBal struct {
	Confirmed   int64 `json:"confirmed"`
	Unconfirmed int64 `json:"unconfirmed"`
}

func getAddressBal(scriptHash string) AddrBal {
	params := []any{scriptHash}
	reqJSON := createElectrumRequest("blockchain.scripthash.get_balance", params)
	elecRes := sendElectrumRequest(reqJSON)
	fmt.Println(elecRes)
	var response BalanceResponse

	// Unmarshal JSON data into the struct
	if err := json.Unmarshal([]byte(elecRes), &response); err != nil {
		fmt.Println("Error:", err)
		return AddrBal{}
	}
	return response.Result
}
