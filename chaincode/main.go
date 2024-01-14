package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zbohm/lirisi/client"
	"github.com/zbohm/lirisi/ring"
)

// KVContractGo defines the Smart Contract structure
type KVContractGo struct {
	contractapi.Contract
}

// Instantiate is called during chaincode instantiation to initialize any data
func (t *KVContractGo) Instantiate(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("KVContractGo Instantiated")
	return nil
}

// Put stores a key-value pair in the ledger
func (t *KVContractGo) Put(ctx contractapi.TransactionContextInterface, key, value string) (string, error) {
	err := ctx.GetStub().PutState(key, []byte(value))
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// Get retrieves a value from the ledger by its key
func (t *KVContractGo) Get(ctx contractapi.TransactionContextInterface, key string) (string, error) {
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
func (t *KVContractGo) PutPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) error {
	transientData, _ := ctx.GetStub().GetTransient()
	message, ok := transientData["message"]
	if !ok {
		return fmt.Errorf("message not found in the transient data")
	}
	return ctx.GetStub().PutPrivateData(collection, "message", message)
}

// GetPrivateMessage retrieves a private message from a specified collection
func (t *KVContractGo) GetPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (string, error) {
	message, err := ctx.GetStub().GetPrivateData(collection, "message")
	if err != nil {
		return "", err
	}
	return string(message), nil
}

// VerifyPrivateMessage verifies the hash of a private message against the stored hash in a specified collection
func (t *KVContractGo) VerifyPrivateMessage(ctx contractapi.TransactionContextInterface, collection string) (bool, error) {
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

// Functions for Linkable Ring Signature ----------------------------------------------------------

// PutVote stores a key-value pair in the ledger
func (t *KVContractGo) PutVote(ctx contractapi.TransactionContextInterface, key, value string) (string, error) {

	// Parse the JSON input.
	type Request struct {
		FoldedPublicKeys string `json:"foldedPublicKeys"`
		Signature        string `json:"signature"`
		Message          string `json:"message"`
	}
	var request Request
	err := json.Unmarshal([]byte(value), &request)
	if err != nil {
		return "", err
	}

	// Validate required parameters.
	if request.FoldedPublicKeys == "" {
		return "", fmt.Errorf("foldedPublicKeys is required")
	}
	if request.Signature == "" {
		return "", fmt.Errorf("signature is required")
	}
	if request.Message == "" {
		return "", fmt.Errorf("message is required")
	}

	// Convert JSON fields to byte arrays.
	foldedPublicKeys := []byte(request.FoldedPublicKeys)
	signature := []byte(key)
	message := []byte(request.Message)

	// Validate the signature.
	verify := client.VerifySignature(foldedPublicKeys, signature, message, []byte(""))
	if verify != ring.Success {
		return "", fmt.Errorf("signature verification failed: %s", ring.ErrorMessages[verify])
	}

	// Store the vote in the ledger.
	err = ctx.GetStub().PutState(key, []byte(value))
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&KVContractGo{})
	if err != nil {
		fmt.Printf("Error creating KVContractGo chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting KVContractGo chaincode: %s", err)
	}
}
