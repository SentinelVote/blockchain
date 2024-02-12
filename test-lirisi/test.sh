#!/bin/sh

#
# This shell script tests the functionalities of the
# lirisi, a Linkable Ring Signature library that has
# a CLI interface.
#

# POSIX `seq`.
pseq () { j=$1 ; while [ "$j" -le "$2" ] ; do printf %s\\n "$j" ; j=$(( j + 1 )) ; done ; }

printf %s\\n "Removing old keys..."
rm -rf pk/
rm -f sk-*.pem
rm -f pk-*.pem
rm -f pk-folded.pem
rm -f vote-*.txt
rm -f sign-*.pem
rm -f verify-*.txt
printf %s\\n "Successfully removed old keys."
sleep 5

amt=5
printf %s\\n "Generating ${amt} key pairs..."
mkdir -p pk
for i in $(pseq 1 ${amt}); do
lirisi genkey -out "sk-${i}.pem"
lirisi pubout -in "sk-${i}.pem" -out "pk/pk-${i}.pem"
done
printf %s\\n "Successfully generated ${amt} public keys at pk/"
printf %s\\n "Successfully generated ${amt} secret keys at sk/"

# Create a double-voter key.
printf %s\\n "Generating a double-voter key called double.pem..."
lirisi genkey -out "sk-double.pem"
lirisi pubout -in "sk-double.pem" -out "pk/pk-double.pem"

# Fold the public keys.
printf %s\\n "Folding the public keys to the current directory..."
lirisi fold-pub -inpath pk -out pk-folded.pem

# Bring the public keys to the current directory.
mv pk/pk-* .
rmdir pk

# Create a malicious key, who is not part of the group.
printf %s\\n "Generating a malicious key called mallory.pem..."
lirisi genkey -out "sk-mallory.pem"
lirisi pubout -in "sk-mallory.pem" -out "pk-mallory.pem"

###############################################################################

# Prepare utility functions to randomize the candidate selection.
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
0) printf %s "Tharman Shanmugaratnam"   ;;
1) printf %s "Ng Kok Song"              ;;
2) printf %s "Tan Kin Lian"             ;;
esac
}

# Test the randomization functions.
printf %s\\n "Testing the randomization functions 5 times."
printf %s\\n "$(gen_candidate)"
printf %s\\n "$(gen_candidate)"
printf %s\\n "$(gen_candidate)"
printf %s\\n "$(gen_candidate)"
printf %s\\n "$(gen_candidate)"

###############################################################################

# Generate the votes.
printf %s\\n "Generating the votes..."
for i in $(pseq 1 ${amt}); do
candidate=$(gen_candidate)
printf %s "$candidate" > "vote-${i}.txt"
lirisi sign -message "$candidate" -inpub "pk-folded.pem" -inkey "sk-${i}.pem" -out "sign-${i}.pem"
done

# Generate the double-vote.
printf %s\\n "Generating the double-vote..."
candidate=$(gen_candidate)
printf %s "$candidate" > "vote-double.txt"
lirisi sign -message "$candidate" -inpub "pk-folded.pem" -inkey "sk-double.pem" -out "sign-double-1.pem"
lirisi sign -message "$candidate" -inpub "pk-folded.pem" -inkey "sk-double.pem" -out "sign-double-2.pem"

# Generate a malicious vote.
printf %s\\n "Generating a malicious vote..."
candidate=$(gen_candidate)
printf %s "$candidate" > "vote-mallory.txt"
lirisi sign -message "$candidate" -inpub "pk-folded.pem" -inkey "sk-mallory.pem" -out "sign-mallory.pem" || true

# Verify the votes.
for i in $(pseq 1 ${amt}); do
lirisi verify -message "vote-${i}.txt" -inpub "pk-folded.pem" -in "sign-${i}.pem" > "verify-${i}.txt"
done

# Verify the both of the double-voter's votes.
printf %s\\n "Verifying the double-voter's votes..."
set -x
lirisi verify -message "vote-double.txt" -inpub "pk-folded.pem" -in "sign-double-1.pem" > "verify-double-1.txt"
lirisi verify -message "vote-double.txt" -inpub "pk-folded.pem" -in "sign-double-2.pem" > "verify-double-2.txt"
diff "verify-double-1.txt" "verify-double-2.txt" --report-identical-files
set +x

printf done\\n
