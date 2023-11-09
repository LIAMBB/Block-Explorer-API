package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/davecgh/go-spew/spew"
)

// TODO: replace the implementations pf this with an interface{}["result"] instead to save on repetitive ElectrumResponse structs
type ElectrumTransactionResponse struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  ElectrumTransaction `json:"result"`
	ID      int                 `json:"id"`
}

type AddrBalHistory struct {
	Block   int     `json:"block"`
	Balance float64 `json:"balance"`
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

type BlockData struct {
	Weight            float64    `json:"weight"`
	Bits              string     `json:"bits"`
	Confirmations     float64    `json:"confirmations"`
	MedianTime        float64    `json:"mediantime"`
	NTx               float64    `json:"nTx"`
	MerkleRoot        string     `json:"merkleroot"`
	Time              float64    `json:"time"`
	Nonce             float64    `json:"nonce"`
	Difficulty        float64    `json:"difficulty"`
	Hash              string     `json:"hash"`
	VersionHex        string     `json:"versionHex"`
	ChainWork         string     `json:"chainwork"`
	Tx                []TxData   `json:"tx"`
	AuxPow            AuxPowData `json:"auxpow"`
	Version           float64    `json:"version"`
	PreviousBlockHash string     `json:"previousblockhash"`
	Height            float64    `json:"height"`
	StrippedSize      float64    `json:"strippedsize"`
}

type TxData struct {
	Locktime float64    `json:"locktime"`
	Vout     []VoutData `json:"vout"`
	Hex      string     `json:"hex"`
	Version  float64    `json:"version"`
	Weight   float64    `json:"weight"`
	Size     float64    `json:"size"`
	Vsize    float64    `json:"vsize"`
	Vin      []VinData  `json:"vin"`
	TxID     string     `json:"txid"`
	Hash     string     `json:"hash"`
}

type VoutData struct {
	Value        float64          `json:"value"`
	N            float64          `json:"n"`
	ScriptPubKey ScriptPubKeyData `json:"scriptPubKey"`
}

type ScriptPubKeyData struct {
	Asm     string `json:"asm"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type VinData struct {
	Sequence  float64       `json:"sequence"`
	TxID      string        `json:"txid"`
	Vout      float64       `json:"vout"`
	ScriptSig ScriptSigData `json:"scriptSig"`
	Coinbase  string        `json:"coinbase"`
}

type ScriptSigData struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type AuxPowData struct {
	Tx                TxData          `json:"tx"`
	ChainIndex        float64         `json:"chainindex"`
	MerkleBranch      []interface{}   `json:"merklebranch"`
	ChainMerkleBranch []interface{}   `json:"chainmerklebranch"`
	ParentBlock       ParentBlockData `json:"parentblock"`
}

type ParentBlockData struct {
	Difficulty float64 `json:"difficulty"`
	Hash       string  `json:"hash"`
	Version    float64 `json:"version"`
	VersionHex string  `json:"versionHex"`
	MerkleRoot string  `json:"merkleroot"`
	Time       float64 `json:"time"`
	Nonce      float64 `json:"nonce"`
	Bits       string  `json:"bits"`
}

type HomeBlock struct {
	Height      int     `json:"height"`
	Hash        string  `json:"hash"`
	Fees        float32 `json:"fees"`
	BlockReward float32 `json:"blockreward"`
	Size        float32 `json:"size"`
	BlockTime   int32   `json:"blocktime"`
	TxCount     int     `json:"txcount"`
	BlockValue  float32 `json:"blockvalue"`
}

type ElectrumTransaction struct {
	TxID          string             `json:"txid"`
	Hash          string             `json:"hash"`
	Version       int                `json:"version"`
	Size          int                `json:"size"`
	Vsize         int                `json:"vsize"`
	Weight        int                `json:"weight"`
	Locktime      int                `json:"locktime"`
	Vin           []ElectrumVinData  `json:"vin"`
	Vout          []ElectrumVoutData `json:"vout"`
	Hex           string             `json:"hex"`
	BlockHash     string             `json:"blockhash"`
	Confirmations int                `json:"confirmations"`
	Time          int64              `json:"time"`
	BlockTime     int64              `json:"blocktime"`
}

type ElectrumVinData struct {
	TxID      string                `json:"txid"`
	Vout      int                   `json:"vout"`
	ScriptSig ElectrumScriptSigData `json:"scriptSig"`
	Sequence  int                   `json:"sequence"`
}

type ElectrumScriptSigData struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type ElectrumVoutData struct {
	Value        float64                  `json:"value"`
	N            int                      `json:"n"`
	ScriptPubKey ElectrumScriptPubKeyData `json:"scriptPubKey"`
}

type ElectrumScriptPubKeyData struct {
	Asm     string `json:"asm"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type HomeBlockTrend struct {
	TxCount    int     `json:"txcount"`
	BlockValue float32 `json:"blockvalue"`
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

type FullTransaction struct {
	TxID   string     `json:"txid"`
	Hex    string     `json:"hex"`
	Height int        `json:"height"`
	Size   int        `json:"size"`
	VSize  int        `json:"vsize"`
	Vin    []FullVin  `json:"vins"`
	Vout   []FullVout `json:"vouts"`
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

// TODO: add error handling
func sendElectrumRequest(jsonRequest string) string {
	// Connect to the server
	conn, err := net.Dial("tcp", electrumURL)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return ""
	}
	defer conn.Close()

	// Send the JSON-RPC request to the server
	_, err = fmt.Fprintf(conn, jsonRequest)
	if err != nil {
		fmt.Println("Error sending JSON-RPC request:", err)
		return ""
	}

	// Read the server's response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return ""
	}

	return response
}

// Based on the following discussion of electrum scripthashing:
// https://github.com/bitcoinjs/bitcoinjs-lib/issues/990
func ElectrumScripthash(addressStr string, chainParams *chaincfg.Params) (string, error) {
	addr, err := btcutil.DecodeAddress(addressStr, chainParams)

	if err != nil {
		fmt.Println("210:", err)
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

// TODO: implement go channels for multi-threading the vin process (requires a lot of electrum requests)
// Current implementation will be pretty slow due to single threaded iteration
func parseBlockTxs(txs []TxData, port int) (float32 /*reward*/, float32 /*fees*/, float32 /*value*/, error /*error*/) {
	var reward float32 = 0.0
	var fees float32 = 0.0
	var value float32 = 0.0

	for _, tx := range txs {
		vinVal := 0.0
		voutVal := 0.0

		for _, vout := range tx.Vout {
			voutVal += vout.Value
		}

		if len(tx.Vin) == 1 && tx.Vin[0].TxID == "" { //Block Reward Tx
			reward += float32(voutVal)
			value += float32(voutVal) // rewards don't have vin or fee but do contribute to block tx value
		} else { //Regular Transaction
			for _, vin := range tx.Vin {
				temp, _ := getTx(vin.TxID, port)
				spew.Dump(tx.Vin)
				fmt.Println("float: ", vin.Vout, " int: ", int(vin.Vout))
				vinVal += temp.Vout[int(vin.Vout)].Value
			}
			fees += float32(vinVal - voutVal)
			value += float32(vinVal)
		}
	}

	// reward, fee, value, err
	return reward, fees, value, nil
}

func getTx(txid string, port int) (ElectrumTransaction, error) {
	params := []any{txid, true} // false=rawTx, true=verboseTx
	reqJSON := createElectrumRequest("blockchain.transaction.get", params)
	res := sendElectrumRequest(reqJSON)

	var response ElectrumTransactionResponse
	err := json.Unmarshal([]byte(res), &response)
	if err != nil {
		return ElectrumTransaction{}, err
	}
	return response.Result, nil
}

func loadHome(coin string) ([]HomeBlock, []HomeBlockTrend, error) {
	port := 0
	if coin == "nmc" {
		port = nmcPort
	} else if coin == "btc" {
		port = btcPort
	} else {
		return []HomeBlock{}, []HomeBlockTrend{}, fmt.Errorf("invalid coin") // Error
	}

	// Get BlockCount
	blockHeight, err := getBlockHeight(port)

	if err != nil {
		fmt.Println("Error getting current blockheight")
	}

	fmt.Println("Blockheight: ", blockHeight)

	var newestBlocks []HomeBlock
	var homeTrends []HomeBlockTrend
	// Get 10 Latest Blocks
	for i := 0; i < 10; i++ {
		blockHash, _ := getBlockHash((blockHeight - i), nmcPort)
		fmt.Println(blockHash)
		block, _ := getBlock(blockHash, nmcPort)
		fmt.Println(block.Height)
		r, f, v, _ := parseBlockTxs(block.Tx, port)
		// Add block to block list
		temp := HomeBlock{
			Height:      int(block.Height),
			Hash:        block.Hash,
			Fees:        f,
			BlockReward: r,
			BlockValue:  v,
			Size:        float32(block.Weight),
			BlockTime:   int32(block.MedianTime),
			TxCount:     int(block.NTx),
		}
		newestBlocks = append(newestBlocks, temp)
		homeTrends = append(homeTrends, HomeBlockTrend{TxCount: int(block.NTx), BlockValue: v})
	}

	return newestBlocks, homeTrends, nil
}

func getBlock(hash string, portNum int) (BlockData, error) {

	method := "getblock"
	params := []interface{}{hash, 2} // verbosity = 2 includes all transactions in block

	result, err := makeRPCRequest(method, params, nmcPort)
	if err != nil {
		fmt.Println("Error:", err)
		return BlockData{}, err
	}

	// Convert the provided data to the MyStruct type
	var myStruct BlockData
	dataJSON, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return BlockData{}, err
	}

	err = json.Unmarshal(dataJSON, &myStruct)
	if err != nil {
		fmt.Println("Error:", err)
		return BlockData{}, err
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blocks, trends, _ := loadHome("nmc")

	var res struct {
		Blocks []HomeBlock      `json:"blocks"`
		Trends []HomeBlockTrend `json:"trends"`
	}

	res.Blocks = blocks
	res.Trends = trends

	resJSON, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error marshaling data", http.StatusInternalServerError)
		return
	}
	fmt.Println()
	// Set headers and write JSON to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
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

	balHist = append(balHist, AddrBalHistory{fullHistTxs[0].Height - 1, 0.0})

	for i, tx := range fullHistTxs {
		balChange := getBalanceChange(tx, addr)
		balance += balChange
		fullHistTxs[i].BalanceChange = balChange
		balHist = append(balHist, AddrBalHistory{tx.Height, balance})
	}
	currentHeight, _ := getBlockHeight(18443)
	if balHist[len(balHist)-1].Block != currentHeight {
		balHist = append(balHist, AddrBalHistory{currentHeight, balHist[len(balHist)-1].Balance})
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