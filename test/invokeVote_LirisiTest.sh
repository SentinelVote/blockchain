#!/bin/sh

token () {
  curl -fsSL --request POST --url http://localhost:8801/user/enroll \
 --header 'Authorization: Bearer' \
 --data "{\"id\": \"admin\", \"secret\": \"adminpw\"}" | jq -r '.token'
}

token="00000000-0000-0000-0000-000000000000-admin"
token=$(token)
uuid=$(uuidgen)

curl --request POST \
  --url http://localhost:8801/invoke/vote-channel/chaincode_vote_go \
  --header "Authorization: Bearer ${token}" \
  --header 'Content-Type: application/json' \
  --data "{\"method\": \"KVContractGo:GenerateAndStorePublicKey\", \"args\": [\"${uuid}\"]}"
