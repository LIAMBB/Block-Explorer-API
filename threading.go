package main

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type ParseBlockTxChannel struct {
	Reward float64
	Fee    float64
	Value  float64
	Error  error
}

// TODO: Error Handling
func parseBlockTxWorker(tx TxData, port int, channel chan ParseBlockTxChannel, wg *sync.WaitGroup) {
	defer (*wg).Done()
	var reward float64 = 0.0
	var fees float64 = 0.0
	var value float64 = 0.0

	vinVal := 0.0
	voutVal := 0.0

	for _, vout := range tx.Vout {
		voutVal += vout.Value
	}

	if len(tx.Vin) == 1 && tx.Vin[0].TxID == "" { //Block Reward Tx
		reward += voutVal
		value += voutVal // rewards don't have vin or fee but do contribute to block tx value
	} else { //Regular Transaction
		for _, vin := range tx.Vin {
			temp, _ := getTx(vin.TxID, port)
			spew.Dump(tx.Vin)
			fmt.Println("float: ", vin.Vout, " int: ", int(vin.Vout))
			vinVal += temp.Vout[int(vin.Vout)].Value
		}
		fees += vinVal - voutVal
		value += vinVal
	}
	res := ParseBlockTxChannel{Reward: reward, Fee: fees, Value: value, Error: nil}
	channel <- res
}

// TODO: implement go channels for multi-threading the vin process (requires a lot of electrum requests)
// Current implementation will be pretty slow due to single threaded iteration
func multiParseBlockTxs(txs []TxData, port int) (float64 /*reward*/, float64 /*fees*/, float64 /*value*/, error /*error*/) {
	var reward float64 = 0.0
	var fees float64 = 0.0
	var value float64 = 0.0

	c := make(chan ParseBlockTxChannel)
	var wg sync.WaitGroup

	for _, tx := range txs {
		wg.Add(1)
		go parseBlockTxWorker(tx, port, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for pTx := range c {
		if pTx.Error != nil {
			return 0, 0, 0, pTx.Error
		}
		reward += pTx.Reward
		fees += pTx.Fee
		value += pTx.Value
	}

	// reward, fee, value, err
	return reward, fees, value, nil
}
