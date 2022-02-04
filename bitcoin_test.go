package bitcoin

import (
	"testing"
)

// Bitcoin RPC settings.
var (
	host = "http://127.0.0.1:18445"
	user = ""
	pass = ""
)

var bitcoin *Bitcoin

func init() {
	bitcoin = Connect(host, user, pass)
}

// Store the result of test functions to be used
// later by other functions.
var (
	address string
	txid    string
	rawtx   string
)

func TestGetNewAddress(t *testing.T) {
	new_address, err := bitcoin.GetNewAddress("")
	if err != nil {
		t.Errorf(err.Error())
	}
	address = new_address.String()
}

func TestImportAddress(t *testing.T) {
	if _, err := bitcoin.ImportAddress(address, "bitcoin"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestValidateAddress(t *testing.T) {
	address_valid, err := bitcoin.ValidateAddress(address)
	if err != nil {
		t.Errorf(err.Error())
	}

	if address_valid.Get("isvalid").Bool() == false {
		t.Errorf("Invalid address.")
	}
}

func TestGetBlockCount(t *testing.T) {
	if _, err := bitcoin.GetBlockCount(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetBalance(t *testing.T) {
	if _, err := bitcoin.GetBalance(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestListUnspent(t *testing.T) {
	list_unspent, err := bitcoin.ListUnspent()
	if err != nil || len(list_unspent.Array()) == 0 {
		t.Errorf(err.Error())
	}
	txid = list_unspent.Array()[0].Get("txid").String()
}

func TestGetTransaction(t *testing.T) {
	if _, err := bitcoin.GetTransaction(txid); err != nil {
		t.Errorf(err.Error())
	}
}

func TestCreateRawTransaction(t *testing.T) {
	inputs := []map[string]interface{}{}
	outputs := []map[string]interface{}{{address: 0.00001}}
	tx, err := bitcoin.CreateRawTransaction(inputs, outputs)
	if err != nil {
		t.Errorf(err.Error())
	}
	rawtx = tx.String()
}

func TestFundRawTransaction(t *testing.T) {
	tx, err := bitcoin.FundRawTransaction(rawtx, 0.00000001)
	if err != nil {
		t.Errorf(err.Error())
	}
	rawtx = tx.Get("hex").String()
}

func TestSignRawTransactionWithWallet(t *testing.T) {
	sign_tx, err := bitcoin.SignRawTransactionWithWallet(rawtx)
	if err != nil {
		t.Errorf(err.Error())
	}
	rawtx = sign_tx.Get("hex").String()
}

func TestDecodeRawTransaction(t *testing.T) {
	if _, err := bitcoin.DecodeRawTransaction(rawtx); err != nil {
		t.Errorf(err.Error())
	}
}

func TestSendRawTransaction(t *testing.T) {
	if _, err := bitcoin.SendRawTransaction(rawtx); err != nil {
		t.Errorf(err.Error())
	}
}
