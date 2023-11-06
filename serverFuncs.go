package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/davecgh/go-spew/spew"
)

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
