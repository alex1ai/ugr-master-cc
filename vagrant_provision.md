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
# -*- mode: ruby -*-
# vi: set ft=ruby :

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
    local.vm.network "forwarded_port", guest:3000, host:3000  
  end
end

```
I setup data first because in provisioning of the webserver, we expect the mongo-daemon to be running already.
Furthermore I have modified playbooks in use here (to seperate DB and server concerns). The plays can be found in the subfolder [provision](https://github.com/alex1ai/ugr-master-cc/tree/master/orquestacion/provision). The most important part though is in data_playbook.yml:

```yml
...

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

Here we need to enable MongoDB to be accessed by other other machines instead only localhost. THIS IS ESSENTIAL. Providing "0.0.0.0" enables connections from any computer basically which suggests a big security issue. For this milestone and testing issues it is okay, but **never use this in production**.

The only thing you have to do to execute everything is (if you are on root level of the project)

```bash
cd orquestacion
vagrant up
``` 

and after some time to initialize the machines, we can access the webserver also via localhost (see portforwarding in the Vagrantfile above) by:

```
curl localhost:3000
{"status": "OK"}%                              
```

Everything working locally! Let's go further and deploy both in Azure...

## Provision on Azure-Cloud
