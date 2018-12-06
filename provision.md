# Provision to VM

This milestone of creating a configuration management for easy provisioning to any VM will be documented in the following.

As a cloud infrastructure I chose Azure, as we have got a Sponsorship there to finish this work.

## Configuration of Virtual Machine

As this project does not make any extensive computations and only deals with text files, 
I chose one of the smalles VM configurations possible (Azure B1s):

- Name of the machine: infogration
- Region: Europe West
- Authentification: SSH
- Open Ports: HTTP, SSH
- Disk: SSD Standard (30 GiB)
- vCPUs: 1
- RAM: 1 Gb
- Operating System: Ubuntu 18.04 LTS server
- DNS name: infogration.westeurope.cloudapp.azure.com

![Azure screenshot](./provision/azure.png)

## Ansible

For configuration management of virtual machines I chose Ansible as we did have a dedicated seminar about it in class and I liked its handling a little more then, for example, chef.
I installed Ansible for Linux as it is described on the Website via pip in a virtual environment.
All files that we need for provisioning are in the subfolder provision in the master branch of this repository. 
it consists of two files:

 - hosts
 - playbook.yml
 
The _hosts_ file contains all ips (or DNS-aliases) which you want to configure, it is also possible to create certain groups:

```
[azure]
infogration.westeurope.cloudapp.azure.com
```

The _playbook.yml_ describes all tasks that should be executed, while idempotence should always be taken care of.
My Playbook for provisioning of a Go-Project looks like this:

- Update APT
- Install Git
- Install daemon (need this in the end to start service)
- Download and Extract Go 1.11.2
- Create typical Go-workspace structure
- Add symbolic link for go/bin/go to execute it via `$ go build`
- Clone github repository in workspace
- Add Port forwarding to iptables (Really important to reach the webservice)
- Download depencies for project
- Install project (creates binary in $HOME/go/bin)
- Kill still running server-process if it is running (important for idempotence)
- Run webservice (via daemon)

Executing it from the command line from the provision folder looks like this:
`$ playbook-ansible -i hosts playbook.yml`

![Image configure VM](./provision/provision.png)

Calling the ip from azure via browser indicates a running webservice:

![Image running service](./provision/running.png)

## Project Extension

In order to go along with the project next to the CC-hitos, I implemented a logging-service for all requests and responses.
As the golang "log"-library is really basic (only has log-levels error/panic/fatal), and we want to document more about what is going on in the webservce (info/warning) I chose [Logrus](https://github.com/Sirupsen/logrus) as it has an active community and a lot of features for improved logging.
As it is api-compatible with the native logger of Go, all it needed basically was to import the package via 
`import log github.com/sirupsen/logrus` and install the package locally via `go get github.com/sirupsen/logrus`.

## Checking playbook of other student

To ensure that everybody can build the webservice on his VM, we check our playbooks by deploying the project of another student.
Sergio Samaniego Martinez ran my playbook successfully, as documented on his [page](https://github.com/samahetfield/PersonalCC-1819/blob/master/docs/hito3.md).

Next, I also downloaded his playbook from his repository and ran it on my Azure VM (after killing running process of my application to make the port available). The result can be seen in the following picture.

![Image sergio playbook](./provision/playbook-sergio.png)
