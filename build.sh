#!/bin/sh
rm ./bin/*
GOOS=windows GOARCH=amd64 go build -o bin/Patrolman.exe main.go
GOOS=linux GOARCH=amd64 go build -o bin/Patrolman_linux main.go 
GOOS=darwin GOARCH=amd64 go build -o bin/Patrolman_mac main.go
