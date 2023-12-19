#!/bin/sh

token () {
  curl -fsSL --request POST --url http://localhost:8801/user/enroll \
 --header 'Authorization: Bearer' \
 --data "{\"id\": \"admin\", \"secret\": \"adminpw\"}" | jq -r '.token'
}

token="00000000-0000-0000-0000-000000000000-admin"
token=$(token)
message_key="00000000-0000-0000-0000-000000000000" # Look up on Hyperledger Explorer.

curl --request POST \
  --url http://localhost:8801/query/vote-channel/cc_vote \
  --header "Authorization: Bearer ${token}" \
  --header 'Content-Type: application/json' \
  --data "{\"method\": \"KVContractGo:get\", \"args\": [\"$message_key\"]}"
