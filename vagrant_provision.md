# Priovision project to different VMs via Vagrant

Here I will document how I made provisioning to different providers (Azure and local) via Vagrant. 

To do so we need the following requirements installed:

- Vagrant
- Virtualization Hipervisor (VirtualBox to be preferred)
- [Azure-Vagrant Plugin](https://github.com/Azure/vagrant-azure)
- Ansible for Provisioning

Follow instructions on their official packages for installation.
Azure-Vagrant Plugin is easily installed via 

`vagrant plugin install vagrant-azure`

At first I will describe how to get a _local_ Setup up and running (including provisioning). After that we add Azure provisioning.
If you get the local version to run, you see that Vagrant is functioning correctly.

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

Last but not least, type `$ vagrant up --provider azure` to start and provision the servers in your Azure subscription.

### Choices made in Vagrantfile

Originally I wanted to do the same as locally - separating service and data in respective VMs. This turned out to be harder than it should be because all of a sudden I got "no resource named XY" errors or ports have not been correctly opened. This is why I refused to spend even more time trying this and went for the simpler solution of just deploying one VM at Azure cloud running Data and Service simultaneously. The next milestone will abstract this separation into Docker containers anyway...

```ruby
# -*- mode: ruby -*-
# vi: set ft=ruby :

require 'vagrant-azure'
# Azure information
TENANT_ID = ENV['AZURE_TENANT_ID']
CLIENT_ID = ENV['AZURE_CLIENT_ID']
CLIENT_SECRET = ENV['AZURE_CLIENT_SECRET']
SUBSCRIPTION_ID = ENV['AZURE_SUBSCRIPTION_ID']

# VM specification
VM_SIZE="Standard_B1s"
LOCATION="francecentral"
RESOURCE_GROUP="vagrant"
NAME="vagrant-server"

Vagrant.configure("2") do |config|
    config.vm.define 'server' do |server|
        server.vm.box = 'azure'
        server.vm.provider :azure do |az, override|
            az.tenant_id = TENANT_ID 
            az.client_id = CLIENT_ID
            az.client_secret = CLIENT_SECRET
            az.subscription_id = SUBSCRIPTION_ID

	    az.vm_name = NAME
            az.vm_size = VM_SIZE

            # az.vm_image_urn = "canonical:ubuntuserver:16.04-LTS:latest"
            az.tcp_endpoints = 80
            az.location = LOCATION
            az.resource_group_name = RESOURCE_GROUP
        end

        server.vm.provision "ansible" do |ansible|
            ansible.compatibility_mode = "2.0"
            ansible.playbook = "./provision/playbook.yml"
        end
    end
    config.ssh.private_key_path = '~/.ssh/id_rsa'
end
```

The script is straightforward to understand. The choices of image/location/size are the same as made and justified in the previous milestone. `az.tcp_endpoints = 80` automatically opens port 80 for server access.

In the end use `vagrant halt` to stop the server again and prevent money loss.
