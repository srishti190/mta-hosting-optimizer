# mta-hosting-optimizer

## Introduction
Service that uncovers the inefficient servers hosting only few active MTAs.

## Dependencies:

- go version: go 1.20
    - Install using following command:
        - `brew install go@1.20`

## Setup
* Install dependency using
```shell
go mod tidy
```
* For checking test coverage
```shell
go test -cover
```
* For running the main program
```shell
go run service.go
```
* For checking the output
    * http://localhost:8082/mta-hosting-optimizer

* Output Data for the default threshold X=1

  {
    "ResultSet": [
    "mta-prod-1",
    "mta-prod-3"
    ],
    "Success": "True",
    "ErrorReason": ""
  }
