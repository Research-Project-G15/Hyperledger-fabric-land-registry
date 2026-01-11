#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}===> Starting network setup...${NC}"

# Check if cryptogen exists
if ! command -v cryptogen &> /dev/null; then
    echo -e "${RED}ERROR: cryptogen not found. Install Fabric binaries.${NC}"
    exit 1
fi

if ! command -v configtxgen &> /dev/null; then
    echo -e "${RED}ERROR: configtxgen not found. Install Fabric binaries.${NC}"
    exit 1
fi

echo -e "${GREEN}===> Generating crypto materials...${NC}"
cryptogen generate --config=./crypto-config.yaml --output="organizations"

echo -e "${GREEN}===> Generating genesis block...${NC}"
export FABRIC_CFG_PATH=${PWD}
configtxgen -profile TwoOrgsApplicationGenesis -channelID system-channel -outputBlock ./channel-artifacts/genesis.block

echo -e "${GREEN}===> Generating channel configuration...${NC}"
export CHANNEL_NAME=landchannel
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME

echo -e "${GREEN}===> Generating anchor peer updates...${NC}"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/LandRegistryMSPanchors.tx -channelID $CHANNEL_NAME -asOrg LandRegistryMSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/BankMSPanchors.tx -channelID $CHANNEL_NAME -asOrg BankMSP

echo -e "${GREEN}===> Setup completed!${NC}"
echo ""
echo "Generated:"
echo "  - organizations/"
echo "  - channel-artifacts/"
echo ""
echo "Next: docker-compose up -d"
