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
	walletURL   = "/wallet/bank" // bank wallet for regtest use
	nmcPort     = 18443
	btcPort     = 18444
)

var (
	nmcParams chaincfg.Params
	btcParms  chaincfg.Params
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
	http.HandleFunc("/nmc/loadHomePage", nmcLoadHomeReq)

	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
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

		if len(tx.Vin) == 0 { //Block Reward Tx
			reward += float32(voutVal)
			value += float32(voutVal) // rewards don't have vin or fee but do contribute to block tx value
		} else { //Regular Transaction
			for _, vin := range tx.Vin {
				temp, _ := getTx(vin.TxID, port)
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
