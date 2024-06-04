package client

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

/**
 * @Author: ygzhang
 * @Date: 2024/6/4 13:20
 * @Func: 连接sdk-client 客户端
 **/

const (
	channelID   = "mychannel"
	chaincodeID = "fabcar"
	orgID       = "Org1"
	userID      = "Admin"
)

func CreateSdkClient() (*ledger.Client, *fabsdk.FabricSDK) {
	configOpt := config.FromFile("./config/sdk-config.yaml")
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		log.Fatalf("创建新的SDK失败: %v\n", err)
		return nil, nil
	}
	log.Printf("---> 创建SDK成功\n")

	var options_user fabsdk.ContextOption
	var options_org fabsdk.ContextOption

	options_user = fabsdk.WithUser(userID)
	options_org = fabsdk.WithOrg(orgID)

	clientChannelContext := sdk.ChannelContext(channelID, options_user, options_org)
	client, err := ledger.New(clientChannelContext)
	if err != nil {
		log.Fatalf("创建sdk客户端失败: %v\n", err)
		return nil, nil
	}
	return client, sdk
}

func CloseSdkClient(sdk *fabsdk.FabricSDK) {
	sdk.Close()
}
