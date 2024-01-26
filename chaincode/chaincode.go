package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	c "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zbohm/lirisi/client"
	"github.com/zbohm/lirisi/ring"
)

// KVContractGo defines the Contract structure
type KVContractGo struct {
	c.Contract
}

// PutFoldedPublicKeys stores a private message in a specified collection
func (t *KVContractGo) PutFoldedPublicKeys(ctx c.TransactionContextInterface, value string) error {
	context := ctx.GetStub()
	// If exists, remove the previous folded public keys.
	_ = context.DelPrivateData("foldedPublicKeys", "foldedPublicKeys")
	return context.PutPrivateData("foldedPublicKeys", "foldedPublicKeys", []byte(value))
}

// GetFoldedPublicKeys retrieves a private message from a specified collection
func (t *KVContractGo) GetFoldedPublicKeys(ctx c.TransactionContextInterface) (string, error) {
	foldedPublicKeys, err := ctx.GetStub().GetPrivateData("foldedPublicKeys", "foldedPublicKeys")
	if err != nil {
		return "", err
	}
	return string(foldedPublicKeys), nil
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

// PutVote stores a vote in the ledger, where key is a UUID and value is a JSON string.
func (t *KVContractGo) PutVote(ctx c.TransactionContextInterface, key, value string) (string, error) {

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
func (t *KVContractGo) GetVotes(ctx c.TransactionContextInterface) (string, error) {

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
	chaincode, err := c.NewChaincode(&KVContractGo{})
	if err != nil {
		fmt.Printf("Error creating KVContractGo chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting KVContractGo chaincode: %s", err)
	}
}
