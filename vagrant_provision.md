# Priovision project to different VMs via Vagrant

Here I will document how I made provisioning to different providers (Azure and local) via Vagrant. 

To do so we need the following requirements installed:

- Vagrant
- Virtualization Hipervisor (VirtualBox to be preferred)
- [Azure-Vagrant Plugin](https://github.com/Azure/vagrant-azure)

Follow instructions on their official packages for installation.
Azure-Vagrant Plugin is easily installed via `vagrant plugin install vagrant-azure`. 

At first I will describe how to get a _local_ Setup up and running (including provisioning). After that we add Azure provisioning.
If you get the local version to run, you see that Vagrant is functioning correctly.

## Provision on the local machine

I am going to create a seperate VM for Webserver and Database respectively. For this my Vagrantfile looks like this:

```ruby
BOX = "ubuntu/xenial64"

Vagrant.configure("2") do |config|

  config.vm.define 'data' do |data|
    data.vm.box = BOX
    data.vm.provision "ansible" do |ansible|
		ansible.playbook = "../provision/data_playbook.yml"
	end
	data.vm.network "private_network", ip: "192.168.50.4"
  end
  
  config.vm.define 'server' do |local|	
	local.vm.box = BOX
    local.vm.provision "ansible" do |ansible|
        ansible.playbook = "../provision/playbook.yml"
    end
	local.vm.network "private_network", ip: "192.168.50.5"
    local.vm.network "forwarded_port", guest:3000, host:4321  
  end

end
```
I setup data first because in provisioning of the webserver, we expect the mongo-daemon to be running already.


## Provision on Azure-Cloud
