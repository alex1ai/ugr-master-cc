# Priovision project to different VMs via Vagrant

Here I will document how I made provisioning to different providers (Azure and local) via Vagrant. 

To do so we need the following requirements installed:

- Vagrant
- Virtualization Hipervisor (VirtualBox to be preferred)
- [Azure-Vagrant Plugin](https://github.com/Azure/vagrant-azure)

Follow instructions on their official packages for installation.
Azure-Vagrant Plugin is easily installed via `vagrant plugin install vagrant-azure`. 

At first I will describe how to get a _local_ VM up and running (including provisioning). After that we add Azure provisioning.
If you get the local version to run, you see that Vagrant is functioning correctly.

## Provision on the local machine

## Provision on Azure-Cloud
