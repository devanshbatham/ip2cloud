#!/bin/bash

# Rename the ip2cloud.py file to watson
mv ip2cloud.py ip2cloud

python3 parse_data.py

# Move the ip2cloud file to /usr/local/bin
sudo mv ip2cloud /usr/local/bin/

# Make the ip2cloud file executable
sudo chmod +x /usr/local/bin/ip2cloud

# Remove the ip2cloud.pyc file if it exists
if [ -f ip2cloud.pyc ]; then
    rm ip2cloud.pyc
fi

echo "ip2cloud has been installed successfully!"