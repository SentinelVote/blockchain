#!/bin/sh
# shellcheck shell=dash

red='\033[0;31m'
RED='\033[1;31m'
blue='\033[0;34m'
BLUE='\033[1;34m'
green='\033[0;32m'
GREEN='\033[1;32m'
default='\033[0m'

# POSIX `seq`.
pseq () { j=$1 ; while [ "$j" -le "$2" ] ; do printf %s\\n "$j" ; j=$(( j + 1 )) ; done ; }
# Fail with a message, RED for "Error:", unset red, then message in default color.
fail () { printf %s%s\\n "${RED}Error:${default}" "$1" ; }

gen_number() {
# Generate a random number in a specified range.
# $1 - Lower limit (inclusive)
# $2 - Upper limit (inclusive)
lower_limit=$1
upper_limit=$2
# Generate a random number using /dev/urandom
random_number=$(od -An -N4 -tu4 /dev/urandom | tr -d ' ')
# Scale the number to the specified range
printf %s\\n "$((lower_limit + random_number % (upper_limit - lower_limit + 1)))"
unset lower_limit upper_limit random_number
}

gen_candidate() {
_random_number=$(gen_number 0 2)
case $_random_number in
0) printf %s "Alice"   ;;
1) printf %s "Bob"     ;;
2) printf %s "Charlie" ;;
esac
}

printf 'User,Password,Token\n' > users.csv

# 1. Admin enrolls himself/herself.
_user='admin'
_password='adminpw'
_token=$(curl -sSL --request POST --url http://localhost:8801/user/enroll \
                   --header 'Authorization: Bearer' \
                   --header 'Content-Type: application/json' \
                   --data "{\"id\": \"$_user\", \"secret\": \"$_password\"}" | jq -r '.token')
printf '%s,%s,%s\n' "$_user" "$_password" "$_token" >> users.csv
_admin_token=$_token

mkdir -p pk sksk/

for i in $(pseq 0 50); do

# 2. Admin registers a user.
_user="user$i"
_password="user$i"
curl -sSL --request POST --url http://localhost:8801/user/register \
          --header "Authorization: Bearer ${_admin_token}" \
          --header 'Content-Type: application/json' \
          --data   "{\"id\": \"${_user}\", \"secret\": \"${_password}\"}"

done

### Voting Phase

for i in $(pseq 0 50); do
:
# 3. User enrolls himself/herself.
_token=$(curl -sSL --request POST --url http://localhost:8801/user/enroll \
                   --header 'Authorization: Bearer' \
                   --header 'Content-Type: application/json' \
                   --data "{\"id\": \"$_user\", \"secret\": \"$_password\"}" | jq -r '.token')
printf '%s,%s,%s\n' "$_user" "$_password" "$_token" >> users.csv

lirisi genkey -out "sk/$_user.sk.pem"
lirisi pubout -in "sk/$_user.sk.pem" -out "pk/$_user.pk.pem"

# 4. User votes for a candidate.
_candidate=$(rand_candidate)
_uuid=$(uuidgen)
curl -sSL --request POST \
          --url http://localhost:8801/invoke/my-channel1/chaincode1 \
          --header "Authorization: Bearer ${_token}" \
          --header 'Content-Type: application/json' \
          --data "{\"method\": \"KVContract:put\", \"args\": [\"${_uuid}\", \"${_candidate}\"]}"

done
