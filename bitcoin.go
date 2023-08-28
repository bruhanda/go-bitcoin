package bitcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
)

type Config struct {
	Host string
	User string
	Pass string
}

type Bitcoin struct {
	*Config
}

func Connect(host string, user string, pass string) *Bitcoin {
	bitcoin := &Bitcoin{
		Config: &Config{
			Host: host,
			User: user,
			Pass: pass,
		},
	}

	if _, err := bitcoin.GetBlockCount(); err != nil {
		log.Fatal(err)
	}
	return bitcoin
}

func (b Bitcoin) Call(params map[string]interface{}) (gjson.Result, error) {
	data, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", b.Host, bytes.NewBuffer(data))
	if err != nil {
		return gjson.Result{}, err
	}
	client := &http.Client{}

	req.SetBasicAuth(b.User, b.Pass)
	res, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return gjson.Result{}, err
	}

	rjson := gjson.ParseBytes(body)
	if rjson.String() == "" {
		err = fmt.Errorf("Username / password invalid.")
		return gjson.Result{}, err
	}

	if rjson.Get("error").String() != "" {
		err = fmt.Errorf(rjson.Get("error").Get("message").String())
		return gjson.Result{}, err
	}
	return rjson.Get("result"), nil
}

func (b Bitcoin) GetBlockchainInfo() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getblockchaininfo",
		"params": []interface{}{},
	}
	return b.Call(data)
}

func (b Bitcoin) GetBlockCount() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getblockcount",
		"params": []interface{}{},
	}
	return b.Call(data)
}

func (b Bitcoin) GetBalance() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getbalance",
		"params": []interface{}{},
	}
	return b.Call(data)
}

// Returns an object with all balances in BTC.
func (b Bitcoin) GetBalances() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getbalances",
		"params": []interface{}{},
	}
	return b.Call(data)
}

// Returns details on the active state of the TX memory pool.
func (b Bitcoin) GetMempoolinfo() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getmempoolinfo",
		"params": []interface{}{},
	}
	return b.Call(data)
}

// Get detailed information about in-wallet transaction
func (b Bitcoin) GetTransaction(txid string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "gettransaction",
		"params": [1]string{txid},
	}
	return b.Call(data)
}

// The getrawtransaction RPC returns the raw transaction data.
func (b Bitcoin) GetRawTransaction(txid string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getrawtransaction",
		"params": []interface{}{txid, true},
	}
	return b.Call(data)
}

// The listunspent RPC returns array of unspent transaction outputs with between minconf and maxconf (inclusive)
// confirmations. Optionally filter to only include txouts paid to specified addresses.
func (b Bitcoin) ListUnspent(minConf, maxConf int, addressesFilter []string, includeUnsafe bool) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "listunspent",
		"params": []interface{}{minConf, maxConf, addressesFilter, includeUnsafe},
	}
	return b.Call(data)
}

// The listtransactions RPC returns up to 'count' most recent transactions skipping the first 'from' transactions.
func (b Bitcoin) ListTransactions(count, skip int) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "listtransactions",
		"params": []interface{}{"*", count, skip},
	}
	return b.Call(data)
}

// The abandontransaction RPC marks an in-wallet transaction and all its in-wallet descendants as abandoned.
// This allows their inputs to be respent.
func (b Bitcoin) Abandontransaction(txid string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "abandontransaction",
		"params": [1]string{txid},
	}
	return b.Call(data)
}

// The importaddress RPC adds an address or script (in hex) that can be watched as if it were in your wallet
// but cannot be used to spend. Requires a new wallet backup.
func (b Bitcoin) ImportAddress(address string, label string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "importaddress",
		"params": [2]string{address, label},
	}
	return b.Call(data)
}

/*
 1. "address"            (string, required) The bitcoin address to send to.
 2. "amount"             (numeric or string, required) The amount in BTC to send. eg 0.1
 3. "comment"            (string, optional) A comment used to store what the transaction is for.
    This is not part of the transaction, just kept in your wallet.
 4. "comment_to"         (string, optional) A comment to store the name of the person or organization
    to which you're sending the transaction. This is not part of the
    transaction, just kept in your wallet.
 5. subtractfeefromamount  (boolean, optional, default=false) The fee will be deducted from the amount being sent.
    The recipient will receive less bitcoins than you enter in the amount field.
 6. replaceable            (boolean, optional) Allow this transaction to be replaced by a transaction with higher fees via BIP 125
 7. conf_target            (numeric, optional) Confirmation target (in blocks)
 8. "estimate_mode"      (string, optional, default=UNSET) The fee estimate mode, must be one of:
    "UNSET"
    "ECONOMICAL"
    "CONSERVATIVE"
*/
func (b Bitcoin) SendToAddress(address string, amount float64, comment, commentTo string, subtractfeefromamount, replaceable bool, confTarget int, estimateMode string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "sendtoaddress",
		"params": []interface{}{address, amount, comment, commentTo, subtractfeefromamount, replaceable, confTarget, estimateMode},
	}
	return b.Call(data)
}

func (b Bitcoin) ImportPrivkey(privkey string, label string, rescan bool) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "importaddress",
		"params": []interface{}{privkey, label, rescan},
	}
	return b.Call(data)
}

func (b Bitcoin) GetNewAddress(label string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getnewaddress",
		"params": [1]string{label},
	}
	return b.Call(data)
}

func (b Bitcoin) CreateRawTransaction(inputs []map[string]interface{}, outputs []map[string]interface{}) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "createrawtransaction",
		"params": []interface{}{inputs, outputs},
	}
	return b.Call(data)
}

func (b Bitcoin) FundRawTransaction(hexstring string, feerate float64) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "fundrawtransaction",
		"params": []interface{}{hexstring, map[string]interface{}{
			"feeRate":     feerate,
			"replaceable": false,
		}},
	}
	return b.Call(data)
}

// get fee from vout
func (b Bitcoin) FundRawTransactionVout(hexstring string, vout int) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "fundrawtransaction",
		"params": []interface{}{hexstring, map[string]interface{}{
			"subtractFeeFromOutputs": []int{vout},
			"replaceable":            false,
		}},
	}
	return b.Call(data)
}

// get fee from vout
func (b Bitcoin) FundRawTransactionVoutFee(hexstring string, vout int, feerate float64) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "fundrawtransaction",
		"params": []interface{}{hexstring, map[string]interface{}{
			"subtractFeeFromOutputs": []int{vout},
			"feeRate":                feerate,
			"replaceable":            false,
		}},
	}
	return b.Call(data)
}

// The fundrawtransaction RPC adds inputs to a transaction until it has enough in value to meet its out value.
func (b Bitcoin) FundRawTransactionWithoutParams(hexstring string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "fundrawtransaction",
		"params": []interface{}{hexstring, map[string]interface{}{
			"replaceable": false,
		}},
	}
	return b.Call(data)
}

func (b Bitcoin) SignRawTransactionWithWallet(rawtx string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "signrawtransactionwithwallet",
		"params": [1]string{rawtx},
	}
	return b.Call(data)
}

func (b Bitcoin) DumpPrivKey(wallet string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "dumpprivkey",
		"params": [1]string{wallet},
	}
	return b.Call(data)
}

// Sign inputs for raw transaction (serialized, hex-encoded).
// The second optional argument (may be null) is an array of previous transaction outputs that
// this transaction depends on but may not yet be in the block chain.
// The third optional argument (may be null) is an array of base58-encoded private
// keys that, if given, will be the only keys used to sign the transaction.
func (b Bitcoin) SignRawTransaction(ewTx string, inputs []map[string]interface{}, privKeys []string) (gjson.Result, error) {
	data := map[string]interface{}{
		"jsonrpc": "1.0",
		"id":      "signrawtransaction",
		"method":  "signrawtransaction",
		"params":  []interface{}{ewTx, inputs, privKeys},
	}

	resp, err := b.Call(data)
	if err != nil {
		return gjson.Result{}, err
	}

	return resp, nil
}

func (b Bitcoin) SignRawTransactionAgain(ewTx string, inputs []map[string]interface{}) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "signrawtransactionwithwallet",
		"params": []interface{}{ewTx, inputs},
	}
	return b.Call(data)
}

func (b Bitcoin) SignRawTransactionWithKey(rawtx string, privkey string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "signrawtransactionwithkey",
		"params": []interface{}{rawtx, [1]string{
			privkey,
		}},
	}
	return b.Call(data)
}

func (b Bitcoin) ValidateAddress(address string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "validateaddress",
		"params": [1]string{address},
	}
	return b.Call(data)
}

func (b Bitcoin) DecodeRawTransaction(rawtx string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "decoderawtransaction",
		"params": [1]string{rawtx},
	}
	return b.Call(data)
}

// Submit a raw transaction (serialized, hex-encoded) to local node and network.
// Note that the transaction will be sent unconditionally to all peers,
// so using this for manual rebroadcast may degrade privacy by leaking the transaction’s origin,
// as nodes will normally not rebroadcast non-wallet transactions already in their mempool.
// Also see createrawtransaction and signrawtransactionwithkey calls.
func (b Bitcoin) SendRawTransaction(rawtx string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "sendrawtransaction",
		"params": []interface{}{rawtx},
	}
	return b.Call(data)
}

func (b Bitcoin) GetAddressInfo(address string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getaddressinfo",
		"params": [1]string{address},
	}
	return b.Call(data)
}

// The testmempoolaccept RPC tests acceptance of a transaction to the mempool without adding it.
func (b Bitcoin) Testmempoolaccept(rawtxs []string, allowhighfees bool) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getaddressinfo",
		"params": []interface{}{rawtxs, allowhighfees},
	}
	return b.Call(data)
}

func (b Bitcoin) Importpubkey(pubkey, label string, rescan bool) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "importpubkey",
		"params": []interface{}{pubkey, label, rescan},
	}
	return b.Call(data)
}

func (b Bitcoin) BumpFee(txid string, fee_rate int64) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "bumpfee",
		"params": []interface{}{
			txid, map[string]interface{}{
				"fee_rate":    fee_rate,
				"replaceable": true,
			},
		},
	}
	return b.Call(data)
}

// EstimateSmartFee stimates the approximate fee per kilobyte needed for a transaction..
// https://bitcoincore.org/en/doc/0.16.0/rpc/util/estimatesmartfee/
func (b *Bitcoin) EstimateSmartFeeWithMode(minconf int, mode string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "estimatesmartfee",
		"params": []interface{}{minconf, mode},
	}
	return b.Call(data)
}
