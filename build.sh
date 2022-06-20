#!/usr/bin/env bash

go build -trimpath -o client ./cmd/client
go build -trimpath -o server ./cmd/server
