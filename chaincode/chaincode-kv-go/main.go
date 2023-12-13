package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// KVContract defines the Smart Contract structure
type KVContract struct {
	contractapi.Contract
}

// Instantiate is called during chaincode instantiation to initialize any data
func (t *KVContract) Instantiate(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("KVContract Instantiated")
	return nil
}

// Put stores a key-value pair in the ledger
func (t *KVContract) Put(ctx contractapi.TransactionContextInterface, key, value string) error {
	return ctx.GetStub().PutState(key, []byte(value))
}

// Get retrieves a value from the ledger by its key
func (t *KVContract) Get(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	value, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", fmt.Errorf("NOT_FOUND")
	}
	return string(value), nil
}

// PutPrivateMessage stores a private message in a specified collection
func (t *KVContract) PutPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) error {
	transientData, _ := ctx.GetStub().GetTransient()
	message, ok := transientData["message"]
	if !ok {
		return fmt.Errorf("message not found in the transient data")
	}
	return ctx.GetStub().PutPrivateData(collection, "message", message)
}

// GetPrivateMessage retrieves a private message from a specified collection
func (t *KVContract) GetPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (string, error) {
	message, err := ctx.GetStub().GetPrivateData(collection, "message")
	if err != nil {
		return "", err
	}
	return string(message), nil
}

// VerifyPrivateMessage verifies the hash of a private message against the stored hash in a specified collection
func (t *KVContract) VerifyPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (bool, error) {
	transientData, _ := ctx.GetStub().GetTransient()
	message, ok := transientData["message"]
	if !ok {
		return false, fmt.Errorf("message not found in the transient data")
	}

	hasher := sha256.New()
	hasher.Write(message)
	currentHash := hex.EncodeToString(hasher.Sum(nil))

	privateDataHash, err := ctx.GetStub().GetPrivateDataHash(collection, "message")
	if err != nil {
		return false, err
	}

	if hex.EncodeToString(privateDataHash) != currentHash {
		return false, fmt.Errorf("VERIFICATION_FAILED")
	}
	return true, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&KVContract{})
	if err != nil {
		fmt.Printf("Error create KVContract chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting KVContract chaincode: %s", err)
	}
}
