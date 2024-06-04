package db

import "time"

/**
 * @Author: ygzhang
 * @Date: 2024/6/4 14:17
 * @Func:
 **/

type BlockchainDB struct {
	Block_number  int
	Previous_hash string
	Data_hash     string
	Tx_id_list    []string
}

type TransactionEntry struct {
	Tx_id                   string
	Timestamp               time.Time
	Chaincode_function_name string
	Tx_content              []string
	Status                  interface{}
	Block_number            int
}

type TransactionDB []TransactionEntry
