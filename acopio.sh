#!/bin/bash

group="CCFrance"
location="francecentral"
name="webserver"
playbook_path="./provision/playbook.yml"

ex=$(az group exists -n $group )
echo "Group does exist: $ex"

if [ "$ex" == "false" ]; then
	echo "Creating group $group in $location"
	az group create -l $location -n $group
fi

# Create webserver and store ip address in tmp file
az vm create -g $group -n $name --size Standard_B1s --image ubuntults | jq '.publicIpAddress' > hosts.txt

# Open Port 80 (and others if wanted) to connect to server
az vm open-port -g $group -n $name --port 80 > /dev/null

# Provision server
ansible-playbook -i hosts.txt $playbook_path

# Delete tmp host file
rm hosts.txt


