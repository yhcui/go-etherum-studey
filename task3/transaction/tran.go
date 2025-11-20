package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/joho/godotenv"
)
import (
	"log"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/types"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	solana_url := os.Getenv("solana_url")
	c := client.NewClient(solana_url)

	// to fetch recent blockhash
	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	feePayerAccount := os.Getenv("feePayer")
	var feePayer, _ = types.AccountFromBase58(feePayerAccount)

	aliceAccount := os.Getenv("alice")
	var alice, _ = types.AccountFromBase58(aliceAccount)

	// create a transfer tx
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer, alice},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.Transfer(system.TransferParam{
					From:   alice.PublicKey,
					To:     common.PublicKeyFromString("HMj9X4VdQXUCkeZAJZpoxCPzjXz4nXL4Pvtb67uG88i9"),
					Amount: 1e8, // 0.1 SOL
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	// send tx
	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	log.Println("txhash:", txhash)
}
