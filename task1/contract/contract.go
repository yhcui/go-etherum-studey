package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	counter "github.com/local/go-etherum-studey/task1"
)

const (
	contractAddr = "0x18872b70e7D201B55fD864e897B87F2D4549D9e5"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal(err)
	}
	counterContract, err := counter.NewCounter(common.HexToAddress(contractAddr), client)

	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.ChainID(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := counterContract.Incre(opt)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		log.Fatal("incorrect receipt status")
	}
	fmt.Printf("交易已确认，所在区块: %d\n", receipt.BlockNumber.Uint64())
	callOpts := &bind.CallOpts{Context: context.Background()}
	cc, err := counterContract.Cc(callOpts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("value is ", cc.Int64())
}
