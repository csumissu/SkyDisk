#!/usr/bin/env bash

go get github.com/99designs/gqlgen/internal/imports@v0.13.0
go get github.com/99designs/gqlgen/internal/code@v0.13.0
go get github.com/99designs/gqlgen/cmd@v0.13.0
go run github.com/99designs/gqlgen

docker compose up -d
go run main.go
