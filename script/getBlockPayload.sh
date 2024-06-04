#!/bin/bash

# 设置全局变量
export PATH=/home/zhangyiguang/fabric/fabric-samples-2.2/bin:$PATH

# 定义变量
CONTAINER_NAME="peer0.org1.example.com"
CHANNEL_NAME="mychannel"
PB_FILE_PATH_IN_CONTAINER="/root"
RESULT_DIR="./result"
BLOCK_DIR="${RESULT_DIR}/block"
TEMP_DIR="${RESULT_DIR}/block/temp"
WORK_PATH="/home/zhangyiguang/fabric/nisl-go-sdk-block/script"

# 进入工作目录
cd ${WORK_PATH} 

# 确保结果目录和区块目录存在
mkdir -p ${RESULT_DIR}
mkdir -p ${BLOCK_DIR}
mkdir -p ${TEMP_DIR}

# 获取最新的区块号
LATEST_BLOCK_NUMBER=$(docker exec ${CONTAINER_NAME} peer channel fetch newest -c ${CHANNEL_NAME} /dev/null 2>&1 | grep "Received block" | awk '{print $NF}')
LATEST_BLOCK_NUMBER=$(($LATEST_BLOCK_NUMBER - 1))

# 遍历所有区块并获取数据
for (( i=0; i<=LATEST_BLOCK_NUMBER; i++ ))
do
    PB_FILE_NAME="block${i}.pb"
    JSON_FILE_NAME="block${i}.json"

    # 1. 进入容器并获取区块信息，保存到.pb文件
    docker exec ${CONTAINER_NAME} peer channel fetch $i -c ${CHANNEL_NAME} ${PB_FILE_PATH_IN_CONTAINER}/${PB_FILE_NAME} >/dev/null 2>&1

    # 2. 将.pb文件从容器复制到宿主机
    docker cp ${CONTAINER_NAME}:${PB_FILE_PATH_IN_CONTAINER}/${PB_FILE_NAME} ${BLOCK_DIR}/ >/dev/null 2>&1

    # 3. 使用configtxlator工具将.pb文件转换为JSON格式
    configtxlator proto_decode --input ${BLOCK_DIR}/${PB_FILE_NAME} --type common.Block | jq . > ${BLOCK_DIR}/${JSON_FILE_NAME}

    echo "Block $i fetched and decoded to ${JSON_FILE_NAME}"
    # 4. 将临时的pb文件存起来
    mv ${BLOCK_DIR}/${PB_FILE_NAME} ${TEMP_DIR}
done
