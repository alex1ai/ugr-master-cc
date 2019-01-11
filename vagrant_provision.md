# Provision project to different VMs via Vagrant

Here I will document how I made provisioning to different providers (Azure and local) via Vagrant. 

To do so we need the following requirements installed:

- Vagrant
- Ansible for Provisioning
- Local provisioning: Virtualization Hipervisor (VirtualBox to be preferred)
- Azure provisioning: [Azure-Vagrant Plugin](https://github.com/Azure/vagrant-azure)

Follow instructions on their official packages for installation.
Azure-Vagrant Plugin is easily installed via 

`vagrant plugin install vagrant-azure`

At first I will describe how to get a _local_ setup up and running (including provisioning). After that we add Azure provisioning.

## Provision on the local machine

I am going to create a seperate VM for Webserver and Database respectively. For this my Vagrantfile looks like this:

```ruby
BOX = "ubuntu/xenial64"
IP_DATA = "192.168.50.2"
IP_SERVER = "192.168.50.3"

Vagrant.configure("2") do |config|
  config.vm.define 'data' do |data|
    data.vm.box = BOX
    data.vm.provision "ansible" do |ansible|
		ansible.playbook = "./provision/data_playbook.yml"
    end
    data.vm.network "private_network", ip: IP_DATA
  end

  config.vm.define 'server' do |local|	
    local.vm.box = BOX
    local.vm.provision "ansible" do |ansible|
      ansible.playbook = "./provision/playbook.yml"
    end
    local.vm.network "private_network", ip: IP_SERVER 
    local.vm.network "forwarded_port", guest:3000, host:8080
  end
end

```
I setup data first because in provisioning of the webserver, we expect the mongo-daemon to be running already.
Furthermore I have modified playbooks in use here (to seperate DB and server concerns). The plays can be found in the subfolder [provision](https://github.com/alex1ai/ugr-master-cc/tree/master/orquestacion/local/provision). The most important part though is in data_playbook.yml:

```yml

- hosts: all
  name: Start Mongo Server
  gather_facts: true
  vars:
    BIND_IP: "0.0.0.0" 
  tasks:
    - include: mongo.yml
    - name: accept Traffic only from {{BIND_IP}} Subnet 
      become: yes
      become_method: sudo
      replace:
        path: /etc/mongod.conf
        regexp: '127\.0\.0\.1'
        replace: "{{ BIND_IP }}"
        backup: yes

    - name: Restart daemon
      become: yes
      become_method: sudo
      service:
         name: mongod
         state: restarted
```

Here we need to enable MongoDB to be accessed by other other machines instead of only localhost. THIS IS ESSENTIAL. Providing "0.0.0.0" enables connections from any computer basically which suggests a big security issue. For this milestone and testing issues it is okay, but **never use this in production**.

As you can see on the Vagrantfile above, I manually added private IP-addresses to both VMs. With this given I can set the environment variable for the webserver, where to look for the mongodb, in the `server_playbook.yml':

```
- environment: 
    MONGO_IP: "192.168.50.2" # IP address of data server
```

The only thing you have to do to execute everything is (if you are on root level of the project)

```bash
cd orquestacion
vagrant up
``` 

and after some time to initialize the machines, we can access the webserver also via localhost (see portforwarding in the Vagrantfile above) by:

```
curl localhost:8080
{"status": "OK"}%                              
```

Everything working locally! Let's go further and deploy in Azure...

## Provision on Azure-Cloud

At first, make sure you are logged in locally in your azure subscription. You can do that through `$ az login`.

Next, you will need to set some environment variables. Either follow these steps at ["Create an Azure Active Directory (AAD) Application"](https://github.com/Azure/vagrant-azure) or if you have the package _jq_ installed you can also just type `$ source set_env.sh` from the orchestration folder - this will set all needed environment variables automatically.

If you haven't already, install the Azure-Vagrant plugin now

`$ vagrant plugin install vagrant-azure`

Finally you need to install the azure-dummy box via

`$ vagrant box add azure https://github.com/azure/vagrant-azure/raw/v2.0/dummy.box --provider azure`

Last but not least, type `$ vagrant up` to start and provision the servers in your Azure subscription.

### Choices made in Vagrantfile

Provisioning through Azure is nearly a full copy of the local version above, except of 2 minor changes. The first change is you have to use the Azure plugin to create the servers: 

```ruby
require 'vagrant-azure'

# Azure information
TENANT_ID = ENV['AZURE_TENANT_ID']
CLIENT_ID = ENV['AZURE_CLIENT_ID']
CLIENT_SECRET = ENV['AZURE_CLIENT_SECRET']
SUBSCRIPTION_ID = ENV['AZURE_SUBSCRIPTION_ID']

# VM specification
VM_SIZE="Standard_B1s"
LOCATION="francecentral"
RESOURCE_GROUP="vagrant-info"
SERVER_NAME="vagrant-server"
DB_NAME="vagrant-data"

IMAGE="Canonical:UbuntuServer:16.04-LTS:latest"

Vagrant.configure("2") do |config|
    config.vm.define 'data' do |server|
        server.vm.box = 'azure'
        server.vm.provider :azure do |az, override|
            az.tenant_id = TENANT_ID 
            az.client_id = CLIENT_ID
            az.client_secret = CLIENT_SECRET
            az.subscription_id = SUBSCRIPTION_ID

            az.vm_name = DB_NAME
            az.vm_size = VM_SIZE

            az.vm_image_urn = IMAGE
            az.tcp_endpoints = 27017 # Allow MongoDB Connections
            az.location = LOCATION
            az.resource_group_name = RESOURCE_GROUP
        end

        server.vm.provision "ansible" do |ansible|
            ansible.compatibility_mode = "2.0"
            ansible.playbook = "./provision/data_playbook.yml"
        end
    end
    config.vm.define 'server' do |server|
        server.vm.box = 'azure'
        server.vm.provider :azure do |az, override|
            az.tenant_id = TENANT_ID 
            az.client_id = CLIENT_ID
            az.client_secret = CLIENT_SECRET
            az.subscription_id = SUBSCRIPTION_ID

            az.vm_name = SERVER_NAME
            az.vm_size = VM_SIZE

            az.vm_image_urn = IMAGE
            az.tcp_endpoints = 80 # Webservice entry point
            az.location = LOCATION
            az.resource_group_name = RESOURCE_GROUP
        end

        server.vm.provision "ansible" do |ansible|
            ansible.compatibility_mode = "2.0"
            ansible.playbook = "./provision/server_playbook.yml"
        end
    end
    config.ssh.private_key_path = '~/.ssh/id_rsa'
end

```

The script is straightforward to understand. The choices of image/location/size are the same as made and justified in the previous milestone. `az.tcp_endpoints = 80` automatically opens port 80 for server access, same with port 27017 for MongoDB.

As they are provisioned together, they are also deployed in the same virtual network automatically by Vagrant. This means we can ping/reach the other virtual machine via its machine name (e.g. `$ ping vagrant-data` from vagrant-server).

For MongoDB location, again, the MONGO_IP environment variable is set in the server_provision.yml file, in this case `MONGO_IP: vagrant-data`.

The IP Adress of the server can be found when looking through the output of `az vm list-ip-addresses`

Screenshot after deploying:
![vagrant deployment](https://github.com/alex1ai/ugr-master-cc/blob/gh-pages/orquestacion/screen.png)

In the end use `vagrant halt` to stop the server again and prevent money loss (or destroy it for good).
