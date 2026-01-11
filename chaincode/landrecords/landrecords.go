package main

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type LandDeed struct {
    DeedID       string `json:"deedID"`
    DocumentHash string `json:"documentHash"`
    RegisteredBy string `json:"registeredBy"`
    Timestamp    string `json:"timestamp"`
    Location     string `json:"location"`
    Status       string `json:"status"`
}

func (s *SmartContract) RegisterDeed(ctx contractapi.TransactionContextInterface, deedID string, documentHash string, registeredBy string, location string) error {
    existing, err := ctx.GetStub().GetState(deedID)
    if err != nil {
        return fmt.Errorf("failed to read: %v", err)
    }
    if existing != nil {
        return fmt.Errorf("deed %s already exists", deedID)
    }

    if len(documentHash) != 64 {
        return fmt.Errorf("invalid hash format")
    }

    txTimestamp, err := ctx.GetStub().GetTxTimestamp()
    if err != nil {
        return fmt.Errorf("failed to get timestamp: %v", err)
    }

    deed := LandDeed{
        DeedID:       deedID,
        DocumentHash: documentHash,
        RegisteredBy: registeredBy,
        Timestamp:    txTimestamp.String(),
        Location:     location,
        Status:       "active",
    }

    deedJSON, err := json.Marshal(deed)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState(deedID, deedJSON)
}

func (s *SmartContract) QueryDeed(ctx contractapi.TransactionContextInterface, deedID string) (*LandDeed, error) {
    deedJSON, err := ctx.GetStub().GetState(deedID)
    if err != nil {
        return nil, fmt.Errorf("failed to read: %v", err)
    }
    if deedJSON == nil {
        return nil, fmt.Errorf("deed %s does not exist", deedID)
    }

    var deed LandDeed
    err = json.Unmarshal(deedJSON, &deed)
    if err != nil {
        return nil, err
    }
    return &deed, nil
}

func (s *SmartContract) VerifyDeed(ctx contractapi.TransactionContextInterface, deedID string, providedHash string) (bool, error) {
    deed, err := s.QueryDeed(ctx, deedID)
    if err != nil {
        return false, err
    }
    if deed.Status == "erased" {
        return false, fmt.Errorf("deed data has been erased")
    }
    return deed.DocumentHash == providedHash, nil
}

func (s *SmartContract) MarkAsErased(ctx contractapi.TransactionContextInterface, deedID string) error {
    deed, err := s.QueryDeed(ctx, deedID)
    if err != nil {
        return err
    }
    deed.Status = "erased"
    deedJSON, err := json.Marshal(deed)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState(deedID, deedJSON)
}

func (s *SmartContract) QueryAllDeeds(ctx contractapi.TransactionContextInterface) ([]*LandDeed, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var deeds []*LandDeed
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }
        var deed LandDeed
        err = json.Unmarshal(queryResponse.Value, &deed)
        if err != nil {
            return nil, err
        }
        deeds = append(deeds, &deed)
    }
    return deeds, nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        fmt.Printf("Error creating chaincode: %v\n", err)
        return
    }
    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %v\n", err)
    }
}
