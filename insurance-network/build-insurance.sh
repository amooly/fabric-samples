#!/usr/bin/env bash
export FABRIC_CFG_PATH=$PWD

## 生成证书
cryptogen generate --config=./insurance-crypto-config.yaml

## 生成创世区块
configtxgen -profile InsuranceOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

## 生成保险业务渠道
configtxgen -profile InsuranceOrgsChannel -outputCreateChannelTx ./channel-artifacts/insuranceChannel.tx -channelID insurance

## 生成基础信息渠道
configtxgen -profile FundamentalOrgsChannel -outputCreateChannelTx ./channel-artifacts/fundamentalChannel.tx -channelID fundamental

# 启动
docker-compose -f docker-compose-insurance.yaml up -d