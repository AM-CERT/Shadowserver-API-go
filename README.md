# Shadowserver-API-go
Simple Golang program to access Shadowserver API

[![Go Report Card](https://goreportcard.com/badge/github.com/AM-CERT/Shadowserver-API-go)](https://goreportcard.com/report/github.com/AM-CERT/Shadowserver-API-go)
![Discord](https://img.shields.io/discord/681699554189377567)
![GitHub last commit](https://img.shields.io/github/last-commit/AM-CERT/Shadowserver-API-go)
![Docker Pulls](https://img.shields.io/docker/pulls/AM-CERT/Shadowserver-API-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/AM-CERT/Shadowserver-API-go)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/AM-CERT/Shadowserver-API-go)

This software is based on [official API client](https://github.com/The-Shadowserver-Foundation/api_utils) and additional details can be found [here](https://github.com/The-Shadowserver-Foundation/api_utils/wiki)

## Description
TODO

## Usage

Download precompiled binaries from 
```
./shadowserver-api-go -method reports/query -param '{"query":{"geo":"AM", "type":["|sinkhole","|honeypot"],"date":"2023-02-14"},"limit":1}'
```
