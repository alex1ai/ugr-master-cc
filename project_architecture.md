# Project architecture

## Demands:

	- 24/7 uptime
	- Internationalization i18n
	- Connection to Android-App & other Frontend-Webservice for Updates
	- Authentication (normal users & admin-users)

## Solution:
	- Microservice
	- REST-API for communication (Written in Go)
		-> Framework: go-restful (https://github.com/emicklei/go-restful)
	- NoSQL data storage, needed feature: internationalization (TODO: search candidates)
	
## Orchestration Tools
**Vagrant** is used to automate creation of virtual machines. It allows automation of creation of VMs and it management via a config file.
I will probably use two machines, one for the REST-API and one for the data base.
Default system will be Ubuntu 18.04

## Deployment Tools
TODO

## Further Abstraction
**Docker** will add another layer of abstraction to this system add adds scalabillity.
