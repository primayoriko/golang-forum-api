#!/bin/bash

go run main.go -- seed
# go test -v ./... # To run all test (order not defined)
go test -v ./api/controllers/post_test.go
go test -v ./api/controllers/thread_test.go
go test -v ./api/controllers/user_test.go
# go run main.go -- migrate     # For production use 
go run main.go -- seed
