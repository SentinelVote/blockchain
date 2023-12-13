#!/bin/sh

token='a2cca5d0-99c2-11ee-b33d-a5e2baffb329-admin'
uuid=$(uuidgen)
message='"{ \"name\": \"John Doe\", \"vote\": \"yes\" }"'
curl --request POST \
  --url http://localhost:8801/invoke/vote-channel/chaincode_vote_go \
  --header "Authorization: Bearer ${token}" \
  --header 'Content-Type: application/json' \
  --data "{\"method\": \"KVContractGo:put\", \"args\": [\"${uuid}\", ${message}]}"
