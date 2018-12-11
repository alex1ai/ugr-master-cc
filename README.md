[![Build Status](https://travis-ci.org/alex1ai/ugr-master-cc.svg?branch=master)](https://travis-ci.org/alex1ai/ugr-master-cc)

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
- Authentication of normal users (get info only) and admins (edit content) (Middleware)

## Deployment
Deployment https://infogration.now.sh

While there are different PaaS deployment services out there, I chose [zeit.co](https://zeit.co/) as it offers nice tutorials and a lightweight deployment mechanism. Furthermore it offers integration of Docker (i.e. I can use any language I want besides Node/JS, which others offer exclusively), CI-support of Travis for automatic testing, and Github-Integration which deploys by pushing to master (if configured this way). 
**No configuration is needed** if you want to use Now. Everything is build automatically and aliased as _infogration_ if someone pushes to master or via Pull Requests. 
Deployment through _Now_ is part of **Travis** execution after it passed all tests. If the tests fail, there will be no new deployment.

## Provision

For easy provisioning on any virtual machine of this webservice, I chose Ansible for its configuration management. Ansible works, in contrast to e.g. Puppet, in a push-way and is a good choice for this kind of project because one doesn't need to install anything on the Client-vm.   
This Project is deployed in Microsoft's Azure cloud with a low-resourced VM. 
How to provision to your cloud and more information about the chosen configuration can be found [here](https://github.com/alex1ai/ugr-master-cc/blob/gh-pages/provision.md). 

MV: 20.188.34.125

As of milestone 4, we can create VMs in Azure automatically via `acopio.sh`. [This](https://github.com/alex1ai/ugr-master-cc/blob/gh-pages/cli-provisioning.md) provides all the documentation and justification for the chosen image, ressources and location. Currently the server is running on Ubuntu 16.04 LTS in France-Central on a cheap virtual machine (B1s).
