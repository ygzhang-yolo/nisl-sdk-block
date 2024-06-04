package block

/**
 * @Author: ygzhang
 * @Date: 2024/6/3 21:12
 * @Func:
 **/
import (
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"log"
)

// QueryAllBlock 查询区块链中的所有区块
func QueryAllBlock(client *ledger.Client) []*common.Block {
	var blocks []*common.Block
	info, err := client.QueryInfo()
	if err != nil {
		log.Fatalf("查询区块链概况失败: %v\n", err)
		return nil
	}

	blockNumber := int64(info.BCI.Height - 1)
	if blockNumber < 0 {
		log.Println("区块链中没有区块")
		return nil
	}

	for i := blockNumber; i >= 0; i-- {
		block, err := client.QueryBlock(uint64(i))
		if err != nil {
			fmt.Printf("查询区块 %d 失败: %v\n", i, err)
			return blocks
		}
		blocks = append(blocks, block)
	}
	return blocks
}

func QueryBlock(client *ledger.Client, number uint64) *common.Block {
	// -------------------- 根据块号查询区块 ----------------
	block, err := client.QueryBlock(number)
	if err != nil {
		fmt.Println("client.QueryBlock() error", err)
		return nil
	}
	return block
}
