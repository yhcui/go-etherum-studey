package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal(err)
	}
	blockNumber := big.NewInt(9583171)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.Hash())
	fmt.Println(block.Nonce())
	fmt.Println(block.Number())
	fmt.Println(block.Time())
	fmt.Println(block.ReceivedAt)
	fmt.Println(block.ReceivedFrom)
	fmt.Println(block.Difficulty())
	fmt.Println(block.BaseFee())
	fmt.Println(block.GasLimit())
	fmt.Println(len(block.Transactions()))

	chainID, err := client.ChainID(context.Background())
	transactions := block.Transactions()
	for _, tx := range transactions {
		fmt.Println("hash:", tx.Hash())
		fmt.Println(tx.Nonce())
		fmt.Println("type:", tx.Type())
		fmt.Println(tx.Value())
		fmt.Println("gas:", tx.Gas())
		fmt.Println(tx.GasPrice())
		fmt.Println("to:", tx.To())

		sender, err := types.Sender(types.LatestSignerForChainID(chainID), tx)
		if err != nil {
			log.Fatal("111:", err)
		} else {
			fmt.Println("sender:", sender.Hex())
		}
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(receipt.Status)
		fmt.Println(receipt.Logs)
		break

	}

	fmt.Println("开始发送交易操作")

	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(100000000000000000) //0.1eth
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xd919c3534e587a08b520133877a818f0e149ef50")
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx sent:", signTx.Hash().Hex())

}
