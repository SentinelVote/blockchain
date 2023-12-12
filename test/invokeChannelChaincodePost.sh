#!/bin/sh

token='20219750-98bc-11ee-ad2b-8d2d942971ad-admin'
uuid=$(uuidgen)
message='"{ \"name\": \"John Doe\", \"vote\": \"yes\" }"'
curl --request POST \
  --url http://localhost:8801/invoke/vote-channel/chaincode_vote \
  --header "Authorization: Bearer ${token}" \
  --header 'Content-Type: application/json' \
  --data "{\"method\": \"KVContract:put\", \"args\": [\"${uuid}\", ${message}]}"
