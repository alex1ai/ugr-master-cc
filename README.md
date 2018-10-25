# infogration-rest
## Description

This project is a RESTful WebService, which will be used for _infogration_, an Android app, which I developed in a previous class in Germany.
It answers the most important questions asked by refugees coming to Germany. Questions range from job-hunting, housing to everyday-life. The content is created by mentors of Diakonie Würzburg, an institution which works with young refugees on a daily basis.
At the moment all the questions and answers are hard-coded in the application, which is undesirable as changes should be made by other persons, too. 

## Architecture

This service stores the data which will be requested (get) by the Android-App, and updated (update/post/delete) by another webservice from responsible persons via GUI. Here we only care about the data-storage REST-service as this needs to run permanently in the cloud, in order to allow continous content upgrades to the users of the app.

### Demands:

	- 24/7 uptime
	- Internationalization
	- Connection to Android-App & other Frontend-Webservice for Updates
	- Authentication (normal users & admin-users)

### Solution:
	- Microservice-based 
	- REST-API for communication (Written in Go)
		-> Framework: go-restful (https://github.com/emicklei/go-restful)
	- NoSQL data storage, needed feature: internationalization (TODO: search candidates)
	- Authentication service (might be included in Go naturally)
	- Optional: Use Google-Translate-API to get translations for content we could not translate manually

## Further Tools
- **Vagrant** is used to automate creation of virtual machines. It allows automation of creation of VMs and it management via a config file.
I will probably use two machines, one for the REST-API and one for the data base.
Default system will be Ubuntu 18.04
- **Travis** Continous Testing at GitHub
- **Docker** will add another layer of abstraction to this system and adds scalabillity.

