#!/bin/bash
#https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes
sudo sysctl -w net.core.rmem_max=7500000
go build -o server . && sudo ./server config.json
