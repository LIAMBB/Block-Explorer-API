package main

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

func sendElectrumRequest(jsonRequest string) {
	// Connect to the server
	conn, err := net.Dial("tcp", electrumURL)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	// Send the JSON-RPC request to the server
	_, err = fmt.Fprintf(conn, jsonRequest)
	if err != nil {
		fmt.Println("Error sending JSON-RPC request:", err)
		return
	}

	// Read the server's response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return
	}

	fmt.Println("Server Response:", response)
}

// Based on the following discussion of electrum scripthashing:
// https://github.com/bitcoinjs/bitcoinjs-lib/issues/990
func ElectrumScripthash(addressStr string, chainParams *chaincfg.Params) (string, error) {
	addr, err := btcutil.DecodeAddress(addressStr, chainParams)

	if err != nil {
		fmt.Println(err)
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
