#!/bin/sh
# shellcheck shell=dash

if [ "${PWD##*/}" != "fabricBlockchain" ]; then
echo "Please run this script from the root directory of the project"
exit 1
fi

# Check if nix-shell is installed
if ! command -v nix-shell >/dev/null 2>&1;
then
echo "nix-shell could not be found"
echo "Please install nix-shell and try again"
exit 1
fi

# Check if nix daemon is running
if ! systemctl is-active --quiet nix-daemon >/dev/null 2>&1;
then
echo "nix-daemon is not running"
echo "Please start the nix-daemon and try again"
exit 1
fi

# Remove the existing go.mod and go.sum files
rm -f chaincode/chaincode-kv-go/go.mod
rm -f chaincode/chaincode-kv-go/go.sum

nix-shell -p go_1_20 --command "
go get github.com/hyperledger/fabric-chaincode-go;
go get github.com/hyperledger/fabric-contract-api-go;
go get github.com/hyperledger/fabric-protos-go;
go get github.com/stretchr/testify;
go get google.golang.org/protobuf;
go get github.com/hyperledger/fabric-chaincode-go/pkg/cid@v0.0.0-20231108144948-3542320d76a7;
go get github.com/hyperledger/fabric-chaincode-go/shim/internal@v0.0.0-20231108144948-3542320d76a7;
go get github.com/hyperledger/fabric-chaincode-go/shim@v0.0.0-20231108144948-3542320d76a7;
go get github.com/hyperledger/fabric-contract-api-go/internal/utils@v1.2.2;
go get github.com/hyperledger/fabric-contract-api-go/metadata@v1.2.2;
go get github.com/hyperledger/fabric-protos-go/peer@v0.3.1;
go get github.com/zbohm/lirisi
"

# From https://github.com/hyperledger/fabric-samples/blob/main/asset-transfer-basic/chaincode-go/go.mod
# go get github.com/hyperledger/fabric-chaincode-go
# go get github.com/hyperledger/fabric-contract-api-go
# go get github.com/hyperledger/fabric-protos-go
# go get github.com/stretchr/testify
# go get google.golang.org/protobuf

# From error messages:
# go get github.com/hyperledger/fabric-chaincode-go/pkg/cid@v0.0.0-20231108144948-3542320d76a7
# go get github.com/hyperledger/fabric-chaincode-go/shim/internal@v0.0.0-20231108144948-3542320d76a7
# go get github.com/hyperledger/fabric-chaincode-go/shim@v0.0.0-20231108144948-3542320d76a7
# go get github.com/hyperledger/fabric-contract-api-go/internal/utils@v1.2.2
# go get github.com/hyperledger/fabric-contract-api-go/metadata@v1.2.2
# go get github.com/hyperledger/fabric-protos-go/peer@v0.3.1
