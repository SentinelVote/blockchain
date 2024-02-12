package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

// VoteContent is the JSON structure of a vote.
type VoteContent struct {
	Candidate    string      `json:"vote"`          // Capitalized name of the candidate.
	Signature    string      `json:"voteSignature"` // Linkable Ring Signature of Vote.
	Constituency string      `json:"constituency"`  // UPPERCASE name of the constituency.
	Hour         json.Number `json:"hour"`          // A number >= 0 and <= 23.
	Valid        bool        `json:"verified"`      // True if the signature is valid.
}

// PutVote stores a vote in the ledger, where key is a UUID and value is a JSON string.
func (t *KVContractGo) PutVote(ctx c.TransactionContextInterface, value string) (string, error) {

	// Get the stored folded public keys.
	foldedPublicKeys, err := ctx.GetStub().GetPrivateData("foldedPublicKeys", "foldedPublicKeys")
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	// Unmarshal the vote.
	var req VoteContent
	if err := json.Unmarshal([]byte(value), &req); err != nil {
		fmt.Println("Error: ", err)
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
	jsonVote, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	// Store the vote in the ledger.
	uuidv7, err := uuid.NewV7()
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	err = ctx.GetStub().PutState(uuidv7.String(), jsonVote)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return "OK", nil
}

// GetVotes retrieves all votes from the ledger, with statistics.
func (t *KVContractGo) GetVotes(ctx c.TransactionContextInterface) (string, error) {

	keys, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		fmt.Println("Error getting keys: ", err)
		return "", err
	}
	defer func(keys shim.StateQueryIteratorInterface) {
		if err := keys.Close(); err != nil {
			fmt.Println("Error closing keys: ", err)
		}
	}(keys)

	var countTotal = 0
	var countHour = [24]int{}
	var countCandidate = make(map[string]int)
	var countConstituency = make(map[string]int)
	var countInvalid = 0

	// Loop over all keys, append the value to the `votes` array.
	var votes []VoteContent
	for keys.HasNext() {
		var vote VoteContent

		// Fetch the next key.
		key, err := keys.Next()
		if err != nil {
			fmt.Println("Error fetching next key: ", err)
			return "", err
		}

		// Unmarshal the vote.
		err = json.Unmarshal(key.Value, &vote)
		if err != nil {
			fmt.Println("Error unmarshalling vote: ", err)
			return "", err
		}

		// Tally invalid votes (if any), then skip to the next vote.
		if vote.Valid != true {
			countInvalid++
			countTotal++
			continue
		}

		hour, err := vote.Hour.Int64()
		if err != nil {
			fmt.Println("Error converting hour to int64: ", err)
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
		CountInvalid      int            `json:"countInvalid"`      // Total number of INVALID votes.
	}{
		CountCandidate:    countCandidate,
		CountConstituency: countConstituency,
		CountHour:         countHour,
		CountTotal:        countTotal,
		Raw:               votes,
		CountInvalid:      countInvalid,
	})
	if err != nil {
		fmt.Println("Error marshalling response: ", err)
		return "", err
	}

	return string(response), nil
}

// PutFoldedPublicKeys stores the folded public keys in its collection.
//
// If the folded public keys were already set, the previous folded public keys will be deleted.
// This accommodates for debugging and testing.
func (t *KVContractGo) PutFoldedPublicKeys(ctx c.TransactionContextInterface, value string) error {
	context := ctx.GetStub()
	_ = context.DelPrivateData("foldedPublicKeys", "foldedPublicKeys")
	return context.PutPrivateData("foldedPublicKeys", "foldedPublicKeys", []byte(value))
}

// GetFoldedPublicKeys retrieves the folded public keys in its collection.
//
// If the folded public keys are not set, it returns "Missing/Unset".
// This allows the client to perform a check without handling an error.
func (t *KVContractGo) GetFoldedPublicKeys(ctx c.TransactionContextInterface) (string, error) {
	foldedPublicKeys, err := ctx.GetStub().GetPrivateData("foldedPublicKeys", "foldedPublicKeys")
	if err != nil {
		fmt.Println("Error: ", err)
		return "Missing/Unset", err
	} else if matched, err := regexp.Match(`-+BEGIN`, foldedPublicKeys); err != nil {
		fmt.Println("Error: ", err)
		return "Missing/Unset", nil
	} else if !matched {
		return "Missing/Unset", nil
	}
	return string(foldedPublicKeys), nil
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
