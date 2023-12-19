#!/bin/sh

token () {
  curl -fsSL --request POST --url http://localhost:8801/user/enroll \
 --header 'Authorization: Bearer' \
 --data "{\"id\": \"admin\", \"secret\": \"adminpw\"}" | jq -r '.token'
}

token="00000000-0000-0000-0000-000000000000-admin"
token=$(token)
uuid=$(uuidgen)
message='"{ \"name\": \"John Doe\", \"vote\": \"yes\" }"'

curl --request POST \
  --url http://localhost:8801/invoke/vote-channel/chaincode_vote \
  --header "Authorization: Bearer ${token}" \
  --header 'Content-Type: application/json' \
  --data "{\"method\": \"KVContractGo:put\", \"args\": [\"${uuid}\", ${message}]}"
