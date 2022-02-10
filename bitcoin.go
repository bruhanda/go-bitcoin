package bitcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

	body, err := ioutil.ReadAll(res.Body)
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

func (b Bitcoin) GetBalances() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "getbalances",
		"params": []interface{}{},
	}
	return b.Call(data)
}

func (b Bitcoin) GetTransaction(txid string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "gettransaction",
		"params": [1]string{txid},
	}
	return b.Call(data)
}

func (b Bitcoin) ListUnspent() (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "listunspent",
		"params": []interface{}{},
	}
	return b.Call(data)
}

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
		"params": [3]string{privkey, label, strconv.FormatBool(rescan)},
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

func (b Bitcoin) FundRawTransaction(rawtx string, feerate int64) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "fundrawtransaction",
		"params": []interface{}{rawtx, map[string]interface{}{
			"fee_rate": feerate,
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

func (b Bitcoin) SendRawTransaction(rawtx string) (gjson.Result, error) {
	data := map[string]interface{}{
		"method": "sendrawtransaction",
		"params": [1]string{rawtx},
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
