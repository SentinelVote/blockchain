#!/bin/sh

token () {
  curl -fsSL --request POST --url http://localhost:8801/user/enroll \
 --header 'Authorization: Bearer' \
 --data "{\"id\": \"admin\", \"secret\": \"adminpw\"}" | jq -r '.token'
}

token="00000000-0000-0000-0000-000000000000-admin"
token=$(token)
message='test' # Your private message to store
transient_data=$(printf "{\"message\":\"%s\"}" "${message}") # Convert the message to a transient data format

curl --request POST \
  --url http://localhost:8801/invoke/vote-channel/chaincode_vote \
  --header "Authorization: Bearer ${token}" \
  --header 'Content-Type: application/json' \
  --data "{\"method\": \"KVContractGo:PutPrivateMessage\", \"args\": [\"collection_public_keys\"], \"transient\": ${transient_data}}"
