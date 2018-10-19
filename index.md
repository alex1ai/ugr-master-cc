# Project description
# infogration-rest

This project is a RESTful WebService, which will be used for the infogration-project (More on this later). Infogration is an android app, which I developed in a previous class in Germany, at this moment looking for funding to establish as a startup. 
This service stores the data which will be requested (get) by the Android-App, and updated (update/post/delete) by another webservice from responsible persons via GUI. Here we only care about the data-storage REST-service as this needs to run permanently in the cloud, in order to allow continous content upgrades to the users of the app.

In this product it is all about internationalization, this means multi-languages (i18n) is one important feature in this project. 

# Suitability for CC-18-19

- It is not clear where this service will run later, thus it is important to establish an environment which works out of the box in a production pipeline at any cloud server.
- The actual implementation of functionality should not be too long (haha), thus the focus can be on setting up a CI-environment with testing and everything.
- This project will actually be used and deployed after the course Cloud Computing, rather then ending up a dead github-repo

# Architecture description

The architecture will be described on an extra page [here](project_architecture.md).

# infoGration

It is an app for Android phones which answers the most important questions asked by refugees coming to Germany. Questions range from job-hunting, housing to everyday-life. The content is created by mentors of Diakonie Würzburg, an institution which works with young refugees on a daily basis.
