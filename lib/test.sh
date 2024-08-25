#!/bin/bash
go clean -testcache
go test ./...
npm run test
