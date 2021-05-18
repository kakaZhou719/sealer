#!/bin/bash
SEALER_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

#install ginkgo
echo "installing ginkgo ..."
go get github.com/onsi/ginkgo/ginkgo

#set ginkgo run options

# run test
echo "starting to test sealer ..."
cd $SEALER_ROOT && ginkgo -v test