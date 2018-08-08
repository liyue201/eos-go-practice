package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/token"
)

const (
	bpUrl          = "https://mainnet.eoscanada.com"
	MainNetChainID = "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906"
	AccountFrom    = ""
	PrivateKey     = ""
	AccountTo      = "gooooooooooe"
)

func printJson(v interface{}) {
	data, err := json.Marshal(v)
	if err == nil {
		fmt.Printf("%v\n", string(data))
	}
}

func testTransfer() {
	api := eos.New(bpUrl)

	info, err := api.GetInfo()
	if err != nil {
		fmt.Printf("err=%#v\n", err)
		return
	}

	fmt.Printf("chan_info:\n")
	printJson(info)

	if hex.EncodeToString(info.ChainID) != MainNetChainID {
		fmt.Println("%v not mainnet chain", hex.EncodeToString(info.ChainID))
		return
	}

	keybag := eos.NewKeyBag()
	keybag.ImportPrivateKey(PrivateKey)
	api.SetSigner(keybag)

	assert, _ := eos.NewAsset("0.0001 EOS")
	action := token.NewTransfer(AccountFrom, AccountTo, assert, "test")
	actions := []*eos.Action{}
	actions = append(actions, action)

	opts := eos.TxOptions{}
	if err := opts.FillFromChain(api); err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	tx := eos.NewTransaction(actions, &opts)

	fmt.Printf("tx=%v\n", tx.ID())

	signedTx, packedTx, err := api.SignTransaction(tx, info.ChainID, eos.CompressionNone)
	if err != nil {
		fmt.Printf("SignTransaction error: %v\n", err.Error())
		return
	}

	fmt.Printf("signed tx: %v\n", signedTx.String())
	fmt.Printf("packed tx: %v\n", packedTx.ID())

	resp, err := api.PushTransaction(packedTx)
	if err != nil {
		fmt.Printf("PushTransaction: %v", err)
		return
	}
	fmt.Printf("PushTransaction resp:\n")
	printJson(resp)
}

func main() {
	testTransfer()
}
