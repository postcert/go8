#!/bin/bash

# Go stuff
echo "Go info:"
go version && echo $(which go)
echo ""

# Delve stuff
echo "Delve info:"
dlv version && echo $(which dlv)
echo ""

# Run expect test
./dlv_test.exp
