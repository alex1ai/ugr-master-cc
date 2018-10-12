# Project architecture

Demands:
	- 24/7 uptime
	- Internationalization i18n
	- Connection to Android-App & other Frontend-Webservice for Updates
	- Authentication (normal users & admin-users)

Brainstorming:

	- Microservice
	- REST-API for communication
	- Other Service: NoSQL data storage, needed feature: internationalization (search candidates)
	- Framework: go-restful (https://github.com/emicklei/go-restful)
