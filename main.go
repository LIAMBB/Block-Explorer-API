package main

import (
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
	electrumURL = "127.0.0.1:50002"
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

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
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
	getFullHistTxs(histTxs)

}

type FullHistTransaction struct {
	TxID          string
	Confirmations int
	Height        int
	Size          int
	VSize         int
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

func getFullHistTxs(histTxs []HistoryTransaction, addr string) {
	tx, _ := getTx(histTx.TxHash, nmcPort)

	// Get Txs for all Vins
	for _, t := range histTxs {

	}
	// Extract their amount, address, txid and index and add to struct as FullVin
	// Go through al vouts and populate the FullVout array

	// calculate confirmations
	// get size
	// add Txid

}
func getFullHistTx(histTx HistoryTransaction, addr string) {

	tx, _ := getTx(histTx.TxHash, 50001)

	currentHeight, _ := getBlockHeight(18443)

	var fullTx FullHistTransaction
	fullTx.TxID = tx.TxID
	fullTx.Confirmations = currentHeight - histTx.Height
	fullTx.Height = histTx.Height
	// Loop over transaction OUTPUTS
	for _, vout := range tx.Vout {
		if vout.Value > 0 {
			// Retrieve address associated with output
			address := vout.ScriptPubKey.Address
			// Add to fullVout struct and append to vout array in fullTx
		}
	}

	// Loop over transaction INPUTS
	for _, vin := range tx.Vin {
		// Block Rewards won't have a TxId
		if vin.TxId != "" {
			// Get transaction associated with this inputs tx id
			vinTx, err := GetTransaction(vin.TxId, electrumURL)
			if err != nil {
				// fmt.Println("286: ", err)
				return 0.0, err
			}

			// Assign address associated with specific output of this tx
			address := vinTx.Vout[vin.Vout].ScriptPubKey.Address

			// Is tx address in wallet? Update tx balance change
			for _, walletAddr := range addresses {
				if walletAddr == address {
					addressBalances[address] -= vinTx.Vout[vin.Vout].Value
				}
			}
		}
	}
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
