# Nisl-sdk-blcok

## 简介

通过fabric-sdk-go方式获取区块链网络的概要信息，高度，哈希等，由于go语言方式的sdk提供的接口和java/node方式的接口有较大的差异; 这里简单记录下

## QuickStart

```bash
# 启动Fabric网络
./startFabric.sh

# 编译并运行go
go build
go run nisl-sdk-block

# 关闭Fabric网络
./networkDown.sh
```

基本的功能, 都在main.go中
```go
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
```


## 项目说明

- `block`：与区块链区块查询相关, 主要提供了一些区块查询的操作;
    - QueryBlock: 查询单个区块
    - QueryAllBlock: 查询所有区块
- `client`: 建立sdk-client, 使用sdk查询区块链区块之前, 需要先建立连接client
- `config`: 与区块链建立连接的sdk
- `db`: 根据区块链的区块建立数据库, 提供了根据Block生成对应的DB的方法
  - BlockchainDB: 区块链数据库
  - TransactionDB: 交易数据库
- `script`: 通过shell的形式, 获取区块链的.pb和.json文件, 直接运行getBlockPayload.sh即可
- `submit`: 提供了提交一些交易, 更新区块的方法, submitTransactions方法
- `wallet`: 不含代码, 只是网络ca信息

## 区块链数据库设计

区块数据库, 记录区块的信息
```go
type BlockchainDB struct {
    Block_number  int      //区块号
    Previous_hash string   //上一个区块哈希
    Data_hash     string   //当前区块数据哈希
    Tx_id_list    []string //当前区块包含的tx_id列表
}
```

交易数据库，记录交易的状态信息
```go
type TransactionEntry struct {
    Tx_id                   string      //交易id
    Timestamp               time.Time   //时间戳
    Chaincode_function_name string      //链码名
    Tx_content              []string    //交易内容载荷
    Status                  interface{} //交易状态
    Block_number            int         //区块号
}

type TransactionDB []TransactionEntry
```

## TODO
