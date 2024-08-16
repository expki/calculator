#!/bin/bash

# Build Client
cd client
rm -r dist/*
npm run build
cd ..

# Copy Client files to Server public
cp -r client/dist/* server/public/

# Build Server
mkdir -p build/
cd server
go build -o ../build/calculator -trimpath -ldflags="-s -w" .
cd ..
chmod +x build/calculator
