#!/bin/bash

# Install golang >= 1.23.0
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz -O /tmp/go1.23.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo "Add `export PATH=$PATH:/usr/local/go/bin` to your $HOME/.profile"

# Install node >= 22
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
nvm install 22
nvm use 22
nvm alias default 22
