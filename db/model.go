package db

import "time"

/**
 * @Author: ygzhang
 * @Date: 2024/6/4 14:17
 * @Func:
 **/

type BlockchainDB struct {
	Block_number  int      //区块号
	Previous_hash string   //上一个区块哈希
	Data_hash     string   //当前区块数据哈希
	Tx_id_list    []string //当前区块包含的tx_id列表
}

type TransactionEntry struct {
	Tx_id                   string      //交易id
	Timestamp               time.Time   //时间戳
	Chaincode_function_name string      //链码名
	Tx_content              []string    //交易内容载荷
	Status                  interface{} //交易状态
	Block_number            int         //区块号
}

type TransactionDB []TransactionEntry
