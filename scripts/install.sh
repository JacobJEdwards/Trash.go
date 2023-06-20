#!/bin/bash
cd ../cmd/trash || exit
go install trash.go
cd ../../scripts || exit 
echo "trash installed"
echo "run trash -h for help" 

