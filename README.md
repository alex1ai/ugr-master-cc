# Project Overview

This repository will be filled with the project for assignment 'Cloud Computing' at the UGR.


## infogration-backend

This project will be a WebService backend for an (Android-)App which is already available in the PlayStore. Until now the content of the App is hardcoded (Prototype-style) in the application. 
It provides a Web-Interface for authorized persons to add, edit, style and translate their own contents. 
Optional Functionality might be Peer-review of content as well...
As the first concern should always be security, different microservices and framework will be used. 

More information about the project coming soon...

## Tools 

- This project will use Docker as a container system to enable easy shipping on different possible servers
- Furthermore it will be a WebService providing a REST-API. This will be written in Go 
- Content must be available in different languages, this means for one instance of a Text, we need to store different translations. Therefore data-storage needs to be efficient about language storage and changes.
- As the first concern should always be security, different microservices and framework will be used for authentication, data storage and requests. 

(I wanted to use Docker and Go for a long time, let's try these trends)
