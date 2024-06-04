package main

import (
	"fmt"
	"nisl-sdk-block/block"
	client2 "nisl-sdk-block/client"
	"nisl-sdk-block/db"
)

func main() {
	// 可以通过SubmitTransactions方法模拟提交一些交易
	//submit.SubmitTransactions()

	// 1. 连接 sdk client
	client, sdk := client2.CreateSdkClient()
	defer client2.CloseSdkClient(sdk)

	// 2. 查询所有区块
	blocks := block.QueryAllBlock(client)
	fmt.Println("所有区块:", blocks)

	// 3. 查询单个区块
	info, _ := client.QueryInfo() //info为区块链区块基本信息
	blockNumber := info.BCI.Height - 1
	latestBlock := block.QueryBlock(client, blockNumber)
	fmt.Println("最新区块:", latestBlock)

	// 4. 根据单个区块生成两种db
	txDB := db.GenerateTxDBFromBlock(latestBlock)
	bcDB := db.GenerateBCDBFromBlock(latestBlock)
	fmt.Println("TransactionDB: ", txDB)
	fmt.Println("---------------------------")
	fmt.Println("BlockChainDB: ", bcDB)
}
