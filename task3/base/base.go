package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	solana_url := os.Getenv("solana_url")

	address := os.Getenv("address") // 钱包
	fmt.Printf("solana url: %s\n", solana_url)

	c := client.NewClient(solana_url)
	resp, err := c.GetVersion(context.TODO())
	balance, err := c.GetBalance(context.Background(), address)
	fmt.Println(balance)

	if err != nil {
		log.Fatalf("GetVersion: %v", err)
	}

	log.Println("GetVersion", resp.SolanaCore)

	info, err := c.GetAccountInfo(context.Background(), address)
	if err != nil {
		log.Fatalf("GetAccountInfo: %v", err)
	}
	log.Printf("GetAccountInfo%v", info)

	balance1, err1 := c.GetBalanceWithConfig(context.TODO(), address, client.GetBalanceConfig{
		Commitment: rpc.CommitmentFinalized,
	})
	if err1 != nil {
		log.Fatalf("GetBalanceWithConfig: %v", err1)
	}
	log.Printf("GetBalanceWithConfig %v", balance1)
}
