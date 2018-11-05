# infogration-rest
The usage of the most recent version is documented [here](https://github.com/alex1ai/ugr-master-cc/blob/gh-pages/usage.md).
## Description

This project is a RESTful WebService, which will be used for _infogration_, an Android app, which I developed in a previous class in Germany.
It answers the most important questions asked by refugees coming to Germany. Questions range from job-hunting, housing to everyday-life. The content is created by mentors of Diakonie Würzburg, an institution which works with young refugees on a daily basis.
At the moment all the questions and answers are hard-coded in the application, which is undesirable as changes should be made by other persons, too. 

## Architecture

This service stores the data which will be requested (get) by the Android-App, and updated (update/post/delete) by another webservice from responsible persons via GUI. Here we only care about the data-handling REST-service as this needs to run permanently in the cloud, in order to allow continous content upgrades to the users of the app.

### Highlights:
- Microservice-based in order to have good maintainability and seperation of concerns.
- Go REST-API for communication with App and GUI, I use [Gorilla Mux](https://github.com/gorilla/mux) router as this is the mostly used one in Go. (Service #1)
- Messenger-service between the microservices will be [RabbitMQ](https://github.com/streadway/amqp) as this 'is the most widely deployed open source message broker'.
- SQL data storage service with PostgreSQL, will run as a stand-alone service with REST (Serivce #2)
- Authentication of normal users (get info only) and admins (edit content) (Service #3)
- Optional: Use Google-Translate-API to get translations for content we could not translate manually (Service #4)

## Deployment [https://infogration.now.sh](https://infogration.now.sh)
While there are different PaaS deployment services out there, I chose [zeit.co](https://zeit.co/) as it offers nice tutorials and a lightweight deployment mechanism. Furthermore it offers intergration of Docker (i.e. I can use any language I want besides Node/JS, which others offer exclusively), CI-support of Travis for automatic testing, and Github-Integration which deploys by pushing to master (if configured this way). 
**No configuration is needed** if you want to use Now. Everything is build automatically and aliased as _infogration_ if someone pushes to master or via Pull Requests. 

## Further Tools (will be included in the near future)
- **Vagrant** is used to automate creation of virtual machines. It allows automation of creation of VMs and it management via a config file.
I will probably use two machines, one for the REST-API and one for the data base.
Default system will be Ubuntu 18.04
- **Travis** Continous Testing at GitHub
- **Docker** will add another layer of abstraction to this system and adds scalabillity.

