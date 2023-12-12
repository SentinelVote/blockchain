#!/bin/sh

curl -sS --request POST --url http://localhost:8801/user/enroll \
 --header 'Authorization: Bearer' \
 --data "{\"id\": \"admin\", \"secret\": \"adminpw\"}"
