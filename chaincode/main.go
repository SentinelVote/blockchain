package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zbohm/lirisi/client"
	"github.com/zbohm/lirisi/ring"
)

// KVContractGo defines the Contract structure
type KVContractGo struct {
	contractapi.Contract
}

// +----------------------------------------------------------------------------------------------+
// |                                General Contract API Functions                                |
// +----------------------------------------------------------------------------------------------+

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

// +----------------------------------------------------------------------------------------------+
// |                          Linkable Ring Signature E-Voting Functions                          |
// +----------------------------------------------------------------------------------------------+

// VoteContent is the JSON structure of a vote.
type VoteContent struct {
	Candidate    string      `json:"vote"`               // Capitalized name of the candidate.
	Signature    string      `json:"voteSignature"`      // Linkable Ring Signature of Vote.
	Constituency string      `json:"constituency"`       // UPPERCASE name of the constituency.
	Hour         json.Number `json:"hour"`               // A number >= 0 and <= 23.
	Valid        bool        `json:"verified,omitempty"` // True if the signature is valid.
}

func (t *KVContractGo) PutFoldedPublicKeys(ctx contractapi.TransactionContextInterface, value string) (string, error) {
	_ = ctx.GetStub().DelState("0") // TODO: Remove. Only for testing to reinsert different folded public keys.
	err := ctx.GetStub().PutState("0", []byte(value))
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// PutVote stores a vote in the ledger, where key is a UUID and value is a JSON string.
func (t *KVContractGo) PutVote(ctx contractapi.TransactionContextInterface, key, value string) (string, error) {

	// Get the stored folded public keys.
	foldedPublicKeys, err := ctx.GetStub().GetState("0")
	if err != nil {
		return "", err
	}

	// Unmarshal the vote.
	var req VoteContent
	if err := json.Unmarshal([]byte(value), &req); err != nil {
		return "", err
	}

	// Validate the vote.
	req.Valid = true
	if req.Constituency == "" {
		req.Valid = false
	} else if hour, err := req.Hour.Int64(); err != nil || hour < 0 || hour > 23 {
		req.Valid = false
	} else if req.Candidate == "" {
		req.Valid = false
	} else if matched, _ := regexp.Match(`-+BEGIN RING SIGNATURE`, []byte(req.Signature)); !matched {
		req.Valid = false
	} else if client.VerifySignature(foldedPublicKeys, []byte(req.Signature), []byte(req.Candidate), []byte("")) != ring.Success {
		req.Valid = false
	}

	// Marshal the vote.
	valueWithValidity, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// Store the vote in the ledger.
	err = ctx.GetStub().PutState(key, valueWithValidity)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// GetVotes retrieves all votes from the ledger, with statistics.
func (t *KVContractGo) GetVotes(ctx contractapi.TransactionContextInterface) (string, error) {

	var votes []VoteContent
	keys, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", err
	}
	defer keys.Close()

	var countTotal = 0
	var countHour = [24]int{}
	var countCandidate = make(map[string]int)
	var countConstituency = make(map[string]int)

	// Loop over all keys, append the value to the `votes` array.
	for keys.HasNext() {

		// Fetch the next key.
		key, err := keys.Next()
		if err != nil {
			return "", err
		}

		// Skip the folded public keys.
		if matched, err := regexp.Match(`-+BEGIN FOLDED PUBLIC KEYS`, key.Value); err != nil || matched {
			continue
		}

		// Unmarshal the vote.
		var vote VoteContent
		err = json.Unmarshal(key.Value, &vote)
		if err != nil {
			return "", err
		}
		// TODO: handle invalid votes, tally them nonetheless.

		hour, err := vote.Hour.Int64()
		if err != nil {
			return "", err
		}

		// Increment the counters.
		countCandidate[vote.Candidate]++
		countConstituency[vote.Constituency]++
		countHour[hour]++
		countTotal++

		votes = append(votes, vote)
	}

	response, err := json.Marshal(struct {
		CountCandidate    map[string]int `json:"countCandidate"`    // Number of votes per candidate.
		CountConstituency map[string]int `json:"countConstituency"` // Number of votes per constituency.
		CountHour         [24]int        `json:"countHour"`         // Number of votes per hour.
		CountTotal        int            `json:"countTotal"`        // Total number of votes.
		Raw               []VoteContent  `json:"raw"`               // Raw votes.
	}{
		CountCandidate:    countCandidate,
		CountConstituency: countConstituency,
		CountHour:         countHour,
		CountTotal:        countTotal,
		Raw:               votes,
	})
	if err != nil {
		return "", err
	}

	return string(response), nil
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
