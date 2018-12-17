#!/bin/bash
set -e

images=( "CentOS" "Debian" "UbuntuLTS")
group="OSTest"
location="francecentral"

ex=$(az group exists -n $group )
echo "Group does exist: $ex"

if [ "$ex" == "false" ]; then
	echo "Creating group $group in $location"
	az group create -l $location -n $group
fi

# Disable manual yes input for provisioning (do not do this in usual scripts!)
export ANSIBLE_HOST_KEY_CHECKING=false


for i in "${images[@]}" 
do
		dnsserver="rhel-dns"
		name="server-$i"
		SECONDS=0
		if [ -z "$(az vm list -g $group | jq '.[] | .name' | grep "$name" )" ]
		then
		  echo "Creating $i server "
		  j=$(az vm create -n $name -g $group --image $i --size Standard_B2s --data-disk-sizes-gb 20 --public-ip-address-dns-name $dnsserver)
		  echo $j | jq '.publicIpAddress' > iptmp.txt
		else 
	      echo "VM exists"
		  if [ "$(az vm show -g $group -n $name -d | jq '.powerState')" != "VM running" ]
		  then
		    echo "Starting stopped VM"
            az vm start -g $group -n $name
	      fi
		  echo "Found IP from existing VM"
		  echo $(az vm show -g $group -n $name -d | jq '.publicIps') > iptmp.txt
        fi
		
		echo "Open port "
		az vm open-port -g $group -n $name --port 80 > /dev/null
		
		echo "Provision server"
		ansible-playbook -i iptmp.txt /home/alex/go/src/github.com/alex1ai/ugr-master-cc/provision/playbook.yml
		
		echo "Running speedtest"
		echo "$i: $(ab -n 500 -c 200 http://$dnsserver.francecentral.cloudapp.azure.com/ | grep "Requests per second") " >> tests.txt
		
		echo "Creating and setting up the server took $SECONDS sec"
		echo "$i: $SECONDS sec" >> tests.txt

		echo "Stop this vm again but already start with next server"
		az vm stop -g $group -n $name  --no-wait
done

# Delete all resources
echo "Deleting resource group"
az group delete -g $group --no-wait --yes

rm iptmp.txt
