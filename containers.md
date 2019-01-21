
# Installation of docker & docker-compose

As I am working on a fedora machine, I am using [this](https://docs.docker.com/install/linux/docker-ce/fedora/#install-using-the-repository) guide of installing docker via repository.

Next I am starting the docker daemon via `$ sudo systemctl start docker` and I verfy that everything is running with

`$ sudo docker run hello-world`

After this returns correctly, we are ready to use it.

To use docker-compose on the local machine, we have to install that extra. I chose the global installation for Linux, but it is also available via [pip or Docker itself](https://docs.docker.com/compose/install/).

To installi it on linux we do
```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```
Which installs docker-compose 1.23.2 on the machine, which is a fairly late version. One only needs to change the version number for other versions.

Running `$ sudo docker-compose --version` indicates a successful installation. 

# Running in Docker environment

I will create 2 containers, one for the REST-API and the second one for MongoDB service. To be able to start different containers and let them communicate among them we will use **docker-compose**. 

## API-Dockerfile

The [Dockerfile](https://github.com/alex1ai/ugr-master-cc/Dockerfile) of the REST-API looks like this:

```Dockerfile
# Use latest stable golang base container, alpine version because it is much smaller than the normal one
FROM golang:1.11-alpine AS build

# Install tools required for project
RUN apk update
RUN apk add --no-cache git

# Get project from github master
RUN go get github.com/alex1ai/ugr-master-cc

# Switch to project dir and download dependencies
WORKDIR /go/src/github.com/alex1ai/ugr-master-cc
RUN go get -d

# Set environment variables to compile go application to use in scratch below (more below)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /bin/infogration

# Get a totally empty image ("start from scratch")
FROM scratch

# Copy the compiled binary from the intermediate container above to the new scratch container
COPY --from=build /bin/infogration /infogration

# Create no-root user to exectute binary for security reason
#RUN groupadd -r infogration && useradd --no-log-init -r -g infogration infogration
#USER infogration

# Start server
ENTRYPOINT ["/infogration"]
```

I commented (nearly) every line above with its function. As Go is a compiled language (in contrast to interpretated ones such as JS or Python), the container, in which the binary is executed later, does not have to have any other functionality. Therefore I could use [Scratch](https://docs.docker.com/samples/library/scratch/) as a base image, which is the smallest possible image (really empty). 

Yet I first ran into a few 

`web_1   | standard_init_linux.go:207: exec user process caused "no such file or directory"` 

errors when executing `$ sudo docker-compose up`. Through some online research I found [this](https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/) blog and also this [docker blog](https://blog.docker.com/2016/09/docker-golang/) which handles this issue.

The problem's solution is this section: 

```Dockerfile
# Set environment variables to compile go application to use in scratch below
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /bin/infogration
```

As go has the ability of cross-compilation (compile binaries for another system architecture as the one the command `go build ..` is exectuted) we need to set some environment variables to make an actual executable binary for the SCRATCH image. Setting those env-variables includes some extra C libraries into the binary application which are still needed for execution, but Scratch doesn't even have those.

## Starting the Webservice via docker-compose

Only using the above Dockerfile with `$ docker run .` will crash soon, as the database is not found (because it is not started and the rest-api does not find it on the default route localhost:27017). 
I am going to use docker-compose to do handle both of those problems (start DB & tell service the adress where to look).
The docker-compose.yml file looks like this:

```yaml
version: '3'
services:
  data:
    image: mongo
    restart: always
    command: --smallfiles
  web:
    build: .
    environment:
      - MONGO_IP=data
    ports:
     - "80:3000"
  ```
  
We are defining two containers to be created, "data" and "web". 
  
"data" uses the official mongo image from Docker Hub. Furthermore the mongo-daemon will be restarted everytime we call `docker-compose up`. On the Docker Hub page of the mongo image I found the tip with `command: --smallfiles`. To quote the documentation :
  
 > Sets MongoDB to use a smaller default file size. The --smallfiles option reduces the initial size for data files and limits the maximum size to 512 megabytes. --smallfiles also reduces the size of each journal file from 1 gigabyte to 128 megabytes. Use --smallfiles if you have a large number of databases that each holds a small quantity of data.
 
This means mongo will block less file space, which after all might save some cents when running in a cloud environment.
 
"web" is build from Dockerfile (see above) which is in the same directory as the docker-compose.yml file. Here we also set the environment variable of the IP (or local DNS in this case) where the webservice will find the database. Using `ports: 80:3000` the container port (**Last number**) will be accessible from the host machine (first number) on the specified port 80 (where the webserver is running). In cloud environments the first number will nearly always be 80/443 for HTTP/HTTPS respectively if there is no further internal (VM) port forwarding.
 
Running `$ sudo docker-compose up -d` will start the webservice in the docker environment. `-d` is for detached mode, which enables us to start it and leave the shell without the command being killed. To get some logs from the running containers, one can use `$ sudo docker-compose logs` to get logs of all running containers or specify a machine, e.g. `$ sudo docker-compose logs web`. 
 
![logs](./containers/logs.png)
 
 To kill the docker environment, use `$ sudo docker-compose down`, again one can specify single machines by adding its name in the end.
 
 The size of the machines is:
 
 ![image size](./containers/sizes.png)
 
 We can see that while MongoDB is HUGE, the webservice only has ~14MB in total.
 
 # Docker Hub
 
 I created an account on Docker Hub to make the webservice publicly available.
 
 Next I created a new repository following the GUI in the browser. The service can be found by the name _alexgrimm/infogration-backend_. 
 
 While I was at the creation, I also setup a Github-Hook for automatic Docker Hub building. I basically followed [this](https://docs.docker.com/docker-hub/builds/) official documentation for connecting the services. 
 
 The Docker Hub version will be updated automatically from now on at every push to master.
 
 [infogration-backend on docker hub](https://cloud.docker.com/u/alexgrimm/repository/docker/alexgrimm/infogration-backend)
 
 From now on the service can be used with its image name _alexgrimm/infogration-backend_ in any Dockerfile/docker-compose.yaml or directly on the CLI (e.g. `$sudo docker run alexgrimm/infogration-backend`). 
 
 **Warning**: Using the Dockerhub version only spins up the service and has no database connected. Yet when starting the service will look for a MongoDB service running somewhere. Therefore you need to provide a MongoDB-IP (and Port if not default) via environment variable when starting the docker hub container (add parameter -e MONGO_IP={ipaddress} to docker run).
 
