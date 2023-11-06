package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

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

// TODO: replace the implementations pf this with an interface{}["result"] instead to save on repetitive ElectrumResponse structs
type ElectrumTransactionResponse struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  ElectrumTransaction `json:"result"`
	ID      int                 `json:"id"`
}

type AddrBalHistory struct {
	Block   int
	Balance float64
}

type FullVin struct {
	TxID    string  `json:"txid"`
	Amount  float64 `json:"amount"`
	Index   int     `json:"index"`
	Address string  `json:"address"`
}

type FullVout struct {
	Amount  float64 `json:"amount"`
	Index   int     `json:"index"`
	Address string  `json:"address"`
}

type FullHistTransaction struct {
	TxID          string     `json:"txid"`
	Confirmations int        `json:"confirmations"`
	Hex           string     `json:"hex"`
	Height        int        `json:"height"`
	Size          int        `json:"size"`
	VSize         int        `json:"vsize"`
	BalanceChange float64    `json:"balchange"`
	Vin           []FullVin  `json:"vins"`
	Vout          []FullVout `json:"vouts"`
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

type BalanceResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	Result  AddrBal `json:"result"`
	ID      int     `json:"id"`
}

type AddrBal struct {
	Confirmed   int64 `json:"confirmed"`
	Unconfirmed int64 `json:"unconfirmed"`
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

	// Create a new router
	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/template", templateEndpoint)
	router.HandleFunc("/nmc/loadhomepage", nmcLoadHomeReq)
	router.HandleFunc("/nmc/address", nmcAddresReq)

	// Set up a handler function to handle CORS headers
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers to allow any origin
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
func nmcAddresReq(w http.ResponseWriter, r *http.Request) {
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
		Address string `json:"address"`
	}

	// Unmarshal the JSON data
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON data", http.StatusBadRequest)
		return
	}

	fmt.Println(req.Address)

	transactionHistory, balanceHistory, balance := getAddress(req.Address, "nmc")

	type res struct {
		Balance        AddrBal               `json:"balance"`
		TxHistory      []FullHistTransaction `json:"txhistory"`
		BalanceHistory []AddrBalHistory      `json:"balancehistory"`
	}

	response := res{
		Balance:        balance,
		TxHistory:      transactionHistory,
		BalanceHistory: balanceHistory,
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

func getAddress(addr string, chain string) ([]FullHistTransaction, []AddrBalHistory, AddrBal) {
	var scriptHash string
	if chain == "nmc" {
		var err error
		scriptHash, err = ElectrumScripthash(addr, &nmcMainnetChainParams)
		fmt.Println("Error: ", err)
	} else if chain == "btc" {
		scriptHash, _ = ElectrumScripthash(addr, &chaincfg.MainNetParams)
	}
	fmt.Println(chain, ": ", scriptHash)

	histTxs := getAddressHist(scriptHash)
	fmt.Println("++++++++++++++++++++++++++++++==")
	addrBal := getAddressBal(scriptHash)
	spew.Dump(histTxs, addrBal)
	// getFullHistTx()

	fullHistTxs := make([]FullHistTransaction, 0)
	for _, t := range histTxs {
		tx := getFullHistTx(t, addr)
		fullHistTxs = append(fullHistTxs, tx)
	}

	// Sort the array by blockheight in ascending order
	sort.Slice(fullHistTxs, func(i, j int) bool {
		return fullHistTxs[i].Height < fullHistTxs[j].Height
	})

	// Calculate balance history
	balance := 0.0
	balHist := make([]AddrBalHistory, 0)
	for i, tx := range fullHistTxs {
		balChange := getBalanceChange(tx, addr)
		balance += balChange
		fullHistTxs[i].BalanceChange = balChange
		balHist = append(balHist, AddrBalHistory{tx.Height, balance})
	}

	// Need ascending order for return
	length := len(fullHistTxs)
	for i := 0; i < length/2; i++ {
		fullHistTxs[i], fullHistTxs[length-i-1] = fullHistTxs[length-i-1], fullHistTxs[i]
	}

	return fullHistTxs, balHist, addrBal

}

func getFullHistTxs(histTxs []HistoryTransaction, addr string) {

	// Extract their amount, address, txid and index and add to struct as FullVin
	// Go through al vouts and populate the FullVout array

	// calculate confirmations
	// get size
	// add Txid

}

func getBalanceChange(fullTx FullHistTransaction, addr string) float64 {
	inputVal := 0.0
	outputVal := 0.0

	for _, vin := range fullTx.Vin {
		if vin.Address == addr {
			inputVal += vin.Amount
		}
	}

	for _, vout := range fullTx.Vout {
		if vout.Address == addr {
			outputVal += vout.Amount
		}
	}

	return outputVal - inputVal
}

func getFullHistTx(histTx HistoryTransaction, addr string) FullHistTransaction {

	tx, _ := getTx(histTx.TxHash, 50001)

	currentHeight, _ := getBlockHeight(18443)

	var fullTx FullHistTransaction
	fullTx.TxID = tx.TxID
	fullTx.Confirmations = currentHeight - histTx.Height
	fullTx.Height = histTx.Height
	fullTx.Size = tx.Size
	fullTx.VSize = tx.Vsize
	fullTx.Hex = tx.Hex
	// Loop over transaction OUTPUTS
	for _, vout := range tx.Vout {
		if vout.Value > 0 {
			// Add to fullVout struct and append to vout array in fullTx
			fullVin := FullVout{Amount: vout.Value, Index: vout.N, Address: vout.ScriptPubKey.Address}
			fullTx.Vout = append(fullTx.Vout, fullVin)
		}
	}

	// Loop over transaction INPUTS
	for _, vin := range tx.Vin {
		// Block Rewards won't have a TxId
		if vin.TxID != "" {
			// Get transaction associated with this inputs tx id
			vinTx, err := getTx(vin.TxID, 50001)
			if err != nil {
				// fmt.Println("286: ", err)
				// return 0.0, err
			}

			// Assign address associated with specific output of this tx
			address := vinTx.Vout[vin.Vout].ScriptPubKey.Address

			fullVin := FullVin{vin.TxID, vinTx.Vout[vin.Vout].Value, vin.Vout, address}
			fullTx.Vin = append(fullTx.Vin, fullVin)
		}
	}

	return fullTx
}

func getAddressHist(scriptHash string) []HistoryTransaction {
	params := []any{scriptHash}
	reqJSON := createElectrumRequest("blockchain.scripthash.get_history", params)
	elecRes := sendElectrumRequest(reqJSON)
	var response Response
	fmt.Println("=====================")
	fmt.Println(elecRes)

	// Unmarshal JSON data into the struct
	if err := json.Unmarshal([]byte(elecRes), &response); err != nil {
		fmt.Println("Error:", err)
		return []HistoryTransaction{}
	}
	return response.Result
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
