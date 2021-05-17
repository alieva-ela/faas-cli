#!/bin/bash
make local-install
sudo rm /usr/local/bin/faas-cli
cp ~/go/bin/faas-cli  /usr/local/bin
echo "faas-cli done"
