#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

NETWORK_PATH="/home/zhangyiguang/fabric/fabric-samples-2.2"

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE=${1:-"go"}
CC_SRC_LANGUAGE=`echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:]`

if [ "$CC_SRC_LANGUAGE" = "go" -o "$CC_SRC_LANGUAGE" = "golang" ] ; then
	CC_SRC_PATH=${NETWORK_PATH}"/chaincode/fabcar/go/"
elif [ "$CC_SRC_LANGUAGE" = "javascript" ]; then
	CC_SRC_PATH="../chaincode/fabcar/javascript/"
elif [ "$CC_SRC_LANGUAGE" = "java" ]; then
	CC_SRC_PATH="../chaincode/fabcar/java"
elif [ "$CC_SRC_LANGUAGE" = "typescript" ]; then
	CC_SRC_PATH="../chaincode/fabcar/typescript/"
else
	echo The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script
	echo Supported chaincode languages are: go, java, javascript, and typescript
	exit 1
fi

# clean out any old identites in the wallets
rm -rf ./wallet/*

# launch network; create channel and join peer to channel
pushd ${NETWORK_PATH}/test-network
./network.sh down
./network.sh up createChannel -ca -s couchdb
./network.sh deployCC -ccn fabcar -ccv 1 -cci initLedger -ccl ${CC_SRC_LANGUAGE} -ccp ${CC_SRC_PATH}
popd

# 修正sdk-config.yaml 配置文件
YAML_FILE="/home/zhangyiguang/fabric/nisl-sdk-block/config/sdk-config.yaml" # 替换为实际的YAML文件路径
# 定义Org1和Org2的keystore目录路径
KEY_DIR_ORG1="/home/zhangyiguang/fabric/fabric-samples-2.2/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore"
KEY_DIR_ORG2="/home/zhangyiguang/fabric/fabric-samples-2.2/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/keystore"

# 查找keystore目录下的私钥文件
KEY_FILE_ORG1=$(ls ${KEY_DIR_ORG1}/*_sk)
KEY_FILE_ORG2=$(ls ${KEY_DIR_ORG2}/*_sk)

# 确保找到一个私钥文件
if [[ -z "${KEY_FILE_ORG1}" ]]; then
  echo "未找到Org1的私钥文件"
  exit 1
fi

if [[ -z "${KEY_FILE_ORG2}" ]]; then
  echo "未找到Org2的私钥文件"
  exit 1
fi

# 使用yq工具更新YAML文件
yq eval ".organizations.Org1.users.Admin.key.path = \"${KEY_FILE_ORG1}\"" -i "${YAML_FILE}"
yq eval ".organizations.Org2.users.Admin.key.path = \"${KEY_FILE_ORG2}\"" -i "${YAML_FILE}"


cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

Next, use the FabCar applications to interact with the deployed FabCar contract.
The FabCar applications are available in multiple programming languages.
Follow the instructions for the programming language of your choice:

Go:

  Start by changing into the "go" directory:
    cd go

  Then, install dependencies and run the test using:
    go run fabcar.go

  The test will invoke the sample client app which perform the following:
    - Import user credentials into the wallet (if they don't already exist there)
    - Submit a transaction to create a new car
    - Evaluate a transaction (query) to return details of this car
    - Submit a transaction to change the owner of this car
    - Evaluate a transaction (query) to return the updated details of this car

EOF
