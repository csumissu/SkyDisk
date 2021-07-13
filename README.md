# SkyDisk

This is an exercise program aimed to practice the GO language and related features.

## Features

- login
- logout
- get user profile
- upload a single file
- list dirs and files
- download a single file
- delete a file or dir
- create a dir
- rename a file or dir
- move a file or dir

## How to Run

> Please make sure Go (>= 1.16) is installed.

```bash
# the server will run on port 3000
./go.sh
```

- Some configurations can be customized, please refer to `./config/resources/dev.ini`.
- Postman could be used to have a try on the APIs, please refer to `docs/sky-disk.postman_collection.json`.

## How to Build

```bash
# a executable binary will be generated in the current directory.
go build -o SkyDisk main.go
# it could be executed directly.
./SkyDisk
```

## Tech Stack

- [Go](https://golang.org/) Go is an open source programming language that makes it easy to build simple, reliable, and
  efficient software.
- [Gin](https://github.com/gin-gonic/gin) Gin is a web framework written in Go (Golang).
- [GORM](https://github.com/go-gorm/gorm) The fantastic ORM library for Golang, aims to be developer friendly.
- [gqlgen](https://github.com/99designs/gqlgen) gqlgen is a Go library for building GraphQL servers without any fuss.
- [jwt-go](https://github.com/golang-jwt/jwt) A go (or 'golang' for search engine friendliness) implementation of JSON
  Web Tokens.