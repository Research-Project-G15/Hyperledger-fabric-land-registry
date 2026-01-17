# Hyperledger Fabric Land Registry System

A blockchain-based land deed verification system for Sri Lankan land administration.

## 🎯 Overview

Permissioned blockchain network using Hyperledger Fabric 2.5.4 for:
- Land deed registration with cryptographic hashes
- Immutable verification and audit trails
- PDPA compliance (Personal Data Protection Act)

## 🏗️ Architecture

- **Organizations:** Land Registry, Bank
- **Network:** 1 Orderer, 2 Peers, 1 Channel
- **Chaincode:** landrecords v2.0 (Go)

## 🚀 Quick Start

See full setup instructions in the repository.

## 📊 Features

- Register deeds with SHA-256 hashes
- Verify deed authenticity
- Query deed history
- PDPA-compliant erasure

## 👥 Research Team

Group 15 - University of Sri Jayewardenepura  
BICT (Hons) Networking - 2025/2026







BEFORE YOU START

# Navigate to your project
cd C:\users\abhir\Desktop\RP\fabric

# Start the network (if not running)
docker-compose up -d

# Wait 10 seconds for network to be ready
Start-Sleep -Seconds 10

# Check all containers are running
docker ps


DEMO 1: Check Current Deeds

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"QueryAllDeeds\",\"Args\":[]}'





DEMO 2: Register New Deed (DEED006)

docker exec cli peer chaincode invoke -o orderer.landregistry.lk:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/landregistry.lk/orderers/orderer.landregistry.lk/msp/tlscacerts/tlsca.landregistry.lk-cert.pem -C landchannel -n landrecords --peerAddresses peer0.landregistry.lk:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/landregistry.lk/peers/peer0.landregistry.lk/tls/ca.crt --peerAddresses peer0.bank.lk:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.lk/peers/peer0.bank.lk/tls/ca.crt -c '{\"function\":\"RegisterDeed\",\"Args\":[\"DEED006\",\"c8e3f5a7b2d4e6f8a1c3b5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7\",\"Registrar003\",\"Galle\"]}'



DEMO 3: Query the New Deed

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"QueryDeed\",\"Args\":[\"DEED006\"]}'




DEMO 4: Verification with CORRECT Hash ✅

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"VerifyDeed\",\"Args\":[\"DEED006\",\"c8e3f5a7b2d4e6f8a1c3b5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7\"]}'




DEMO 5: Tampering Detection with WRONG Hash ❌

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"VerifyDeed\",\"Args\":[\"DEED006\",\"FAKE000000000000000000000000000000000000000000000000000000000000\"]}'




DEMO 6: Data Distribution - Query from BOTH Peers

Query from Land Registry peer:

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"QueryDeed\",\"Args\":[\"DEED006\"]}'



Query from Bank peer:

docker exec -e CORE_PEER_LOCALMSPID=BankMSP -e CORE_PEER_ADDRESS=peer0.bank.lk:9051 -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.lk/users/Admin@bank.lk/msp -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.lk/peers/peer0.bank.lk/tls/ca.crt cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"QueryDeed\",\"Args\":[\"DEED006\"]}'





DEMO 7: View All Deeds

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"QueryAllDeeds\",\"Args\":[]}'



DEMO 8: PDPA Compliance - Mark as Erased

docker exec cli peer chaincode invoke -o orderer.landregistry.lk:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/landregistry.lk/orderers/orderer.landregistry.lk/msp/tlscacerts/tlsca.landregistry.lk-cert.pem -C landchannel -n landrecords --peerAddresses peer0.landregistry.lk:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/landregistry.lk/peers/peer0.landregistry.lk/tls/ca.crt --peerAddresses peer0.bank.lk:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/peerOrganizations/bank.lk/peers/peer0.bank.lk/tls/ca.crt -c '{\"function\":\"MarkAsErased\",\"Args\":[\"DEED006\"]}'



docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"QueryDeed\",\"Args\":[\"DEED006\"]}'



DEMO 9: Show Verification Fails After Erasure

docker exec cli peer chaincode query -C landchannel -n landrecords -c '{\"function\":\"VerifyDeed\",\"Args\":[\"DEED006\",\"c8e3f5a7b2d4e6f8a1c3b5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3e5f7\"]}'
