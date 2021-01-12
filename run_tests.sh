#!/bin/bash

# go run main.go -- seed
# go test -v ./...
go test -v ./api/controllers
# go run main.go -- migrate     # For production use 
# go run main.go -- seed
