package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
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

func nmcBlockReq(w http.ResponseWriter, r *http.Request) {
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
		BlockHash   string `json:"blockhash"`
		BlockHeight int    `json:"blockheight"`
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

	if req.BlockHash == "" && req.BlockHeight == 0 {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if req.BlockHash == "" {
		req.BlockHash, _ = getBlockHash(req.BlockHeight, 18443)
	}

	block := getBlockData(req.BlockHash, "nmc")
	//================================================================================//
	//================================================================================//
	//================================================================================//

	// // Marshal the struct into JSON
	resJSON, err := json.Marshal(block)
	if err != nil {
		http.Error(w, "Error marshaling data", http.StatusInternalServerError)
		return
	}

	// Set headers and write JSON to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

type FullBlock struct {
	Weight            float64           `json:"weight"`
	Bits              string            `json:"bits"`
	Confirmations     float64           `json:"confirmations"`
	MedianTime        float64           `json:"mediantime"`
	NTx               float64           `json:"nTx"`
	MerkleRoot        string            `json:"merkleroot"`
	Time              float64           `json:"time"`
	Nonce             float64           `json:"nonce"`
	Difficulty        float64           `json:"difficulty"`
	Hash              string            `json:"hash"`
	VersionHex        string            `json:"versionHex"`
	ChainWork         string            `json:"chainwork"`
	Tx                []FullTransaction `json:"tx"`
	AuxPow            AuxPowData        `json:"auxpow"`
	Version           float64           `json:"version"`
	PreviousBlockHash string            `json:"previousblockhash"`
	Height            float64           `json:"height"`
	StrippedSize      float64           `json:"strippedsize"`
}

func getBlockData(blockHash string, chain string) FullBlock {
	port := 0
	if chain == "nmc" {
		port = 18443
	} else if chain == "btc" {
		port = 0
	}

	block, _ := getBlock(blockHash, port)
	var fullBlock FullBlock
	fullBlock.Weight = block.Weight
	fullBlock.Bits = block.Bits
	fullBlock.Confirmations = block.Confirmations
	fullBlock.MedianTime = block.MedianTime
	fullBlock.NTx = block.NTx
	fullBlock.MerkleRoot = block.MerkleRoot
	fullBlock.Time = block.Time
	fullBlock.Nonce = block.Nonce
	fullBlock.Difficulty = block.Difficulty
	fullBlock.Hash = block.Hash
	fullBlock.VersionHex = block.VersionHex
	fullBlock.ChainWork = block.ChainWork
	fullBlock.AuxPow = block.AuxPow
	fullBlock.Version = block.Version
	fullBlock.PreviousBlockHash = block.PreviousBlockHash
	fullBlock.Height = block.Height
	fullBlock.StrippedSize = block.StrippedSize

	for _, tx := range block.Tx {
		fullTx := getFullTx(tx.TxID)
		fullBlock.Tx = append(fullBlock.Tx, fullTx)
	}

	return fullBlock
}

type FullTransaction struct {
	TxID   string     `json:"txid"`
	Hex    string     `json:"hex"`
	Height int        `json:"height"`
	Size   int        `json:"size"`
	VSize  int        `json:"vsize"`
	Vin    []FullVin  `json:"vins"`
	Vout   []FullVout `json:"vouts"`
}

func getFullTx(txid string) FullTransaction {

	tx, _ := getTx(txid, 50001)

	var fullTx FullTransaction
	fullTx.TxID = tx.TxID
	// fullTx.Height = histTx.Height
	block, _ := getBlock(tx.BlockHash, 18443)

	fullTx.Height = int(block.Height)
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
