#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -ex
NETWORK_PATH="/home/zhangyiguang/fabric/fabric-samples-2.2"

# Bring the test network down
pushd ${NETWORK_PATH}/test-network
./network.sh down
popd

# clean out any old identites in the wallets
rm -rf ./wallet/*

