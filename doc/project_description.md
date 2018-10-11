# Project description
# infogration-rest

This project is a RESTful WebService, which will be used for the infogration-project (More on this later). Infogration is an android app, which I developed in a previous class in Germany, at this moment looking for funding to establish as a startup. 
This service stores the data which will be requested (get) by the Android-App, and updated (update/post/delete) by another webservice from responsible persons via GUI. Here we only care about the data-storage REST-service as this needs to run permanently in the cloud, in order to allow continous content upgrades to the users of the app.

The content that is being queried are **Content** instances, where each instance has as **Question** and **Answer** instance. 
In this product it is all about internationalization, this means multi-languages (i18n) is one important feature in this project. 

# Suitability for CC-18-19

- It is not clear where this service will run later, thus it is important to establish an environment which works out of the box at any cloud server.
- The actual implementation of functionality should not be too long (haha), thus the focus can be on setting up a CI-environment with testing and everything.
- This project will actually be used and deployed after the course Cloud Computing, rather then ending up a dead github-repo
