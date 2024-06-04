package db

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"log"
	"strings"
	"time"
)

/**
 * @Author: ygzhang
 * @Date: 2024/6/4 14:18
 * @Func:
 **/

func GenerateTxDBFromBlock(block *common.Block) *TransactionDB {
	txIDList := []string{}
	var transactionDatabase TransactionDB
	for _, data := range block.Data.Data {
		envelope := &common.Envelope{}
		if err := proto.Unmarshal(data, envelope); err != nil {
			log.Fatalf("解析Envelope失败: %v\n", err)
		}

		payload := &common.Payload{}
		if err := proto.Unmarshal(envelope.Payload, payload); err != nil {
			log.Fatalf("解析Payload失败: %v\n", err)
		}

		channelHeader := &common.ChannelHeader{}
		if err := proto.Unmarshal(payload.Header.ChannelHeader, channelHeader); err != nil {
			log.Fatalf("解析ChannelHeader失败: %v\n", err)
		}
		// 获取tx id list
		txID := channelHeader.TxId
		txIDList = append(txIDList, txID)

		// 解析交易数据
		tx, err := getTransaction(payload.Data)
		if err != nil {
			log.Fatalf("解析交易数据失败: %v\n", err)
		}

		for _, action := range tx.Actions {
			payload := &peer.ChaincodeActionPayload{}
			err := proto.Unmarshal(action.Payload, payload)
			if err != nil {
				log.Fatalf("解析ChaincodeActionPayload失败: %v\n", err)
			}

			// 解析ProposalResponsePayload
			proposalResponsePayload := &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(payload.Action.ProposalResponsePayload, proposalResponsePayload)
			if err != nil {
				log.Fatalf("解析ProposalResponsePayload失败: %v\n", err)
			}

			// 解析ChaincodeAction
			ccAction := &peer.ChaincodeAction{}
			err = proto.Unmarshal(proposalResponsePayload.Extension, ccAction)
			if err != nil {
				log.Fatalf("解析ChaincodeAction失败: %v\n", err)
			}

			// 解析ChaincodeProposalPayload
			ccProposalPayload := &peer.ChaincodeProposalPayload{}
			err = proto.Unmarshal(payload.ChaincodeProposalPayload, ccProposalPayload)
			if err != nil {
				log.Fatalf("解析ChaincodeProposalPayload失败: %v\n", err)
			}

			// 解析ChaincodeInvocationSpec
			chaincodeInvocationSpec := &peer.ChaincodeInvocationSpec{}
			err = proto.Unmarshal(ccProposalPayload.Input, chaincodeInvocationSpec)
			if err != nil {
				log.Fatalf("解析ChaincodeInvocationSpec失败: %v\n", err)
			}

			chaincodeName := chaincodeInvocationSpec.ChaincodeSpec.ChaincodeId.Name
			functionName := chaincodeInvocationSpec.ChaincodeSpec.Input.Args[0]
			function := fmt.Sprintf("%s.%s", chaincodeName, string(functionName))

			timestamp := time.Unix(channelHeader.Timestamp.Seconds, int64(channelHeader.Timestamp.Nanos))

			// 将 txContent 转换为 []string 类型并去掉第一个元素
			txContentArray := strings.Split(string(ccProposalPayload.Input), "\n")
			txContent := txContentArray[1:]

			status := ccAction.Response.Status
			txEntry := TransactionEntry{
				Tx_id:                   txID,
				Timestamp:               timestamp,
				Chaincode_function_name: function,
				Tx_content:              txContent,
				Status:                  status,
				Block_number:            int(block.Header.Number),
			}

			transactionDatabase = append(transactionDatabase, txEntry)
		}
	}
	return &transactionDatabase
}

func GenerateBCDBFromBlock(block *common.Block) *BlockchainDB {
	txIDList := []string{}

	for _, data := range block.Data.Data {
		envelope := &common.Envelope{}
		if err := proto.Unmarshal(data, envelope); err != nil {
			log.Fatalf("解析Envelope失败: %v\n", err)
		}

		payload := &common.Payload{}
		if err := proto.Unmarshal(envelope.Payload, payload); err != nil {
			log.Fatalf("解析Payload失败: %v\n", err)
		}

		channelHeader := &common.ChannelHeader{}
		if err := proto.Unmarshal(payload.Header.ChannelHeader, channelHeader); err != nil {
			log.Fatalf("解析ChannelHeader失败: %v\n", err)
		}

		txID := channelHeader.TxId
		txIDList = append(txIDList, txID)
	}

	return &BlockchainDB{
		Block_number:  int(block.Header.Number),
		Previous_hash: hex.EncodeToString(block.Header.PreviousHash),
		Data_hash:     hex.EncodeToString(block.Header.DataHash),
		Tx_id_list:    txIDList,
	}
}

func getTransaction(data []byte) (*peer.Transaction, error) {
	tx := &peer.Transaction{}
	err := proto.Unmarshal(data, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
