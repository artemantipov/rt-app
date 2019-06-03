# rt-app
Some basic REST app with deploy to GCP(GAE)

# System Diagram

* Blue lines - current solution 
* Grat lines - possible alternative solutions supported by app

![alt text](https://github.com/artemantipov/rt-app/blob/master/diagram.png)

# Deploy to GCP

## Requirements
* gcloud sdk 
* golang 1.11+

## Preparations
* Create file and fill with ENV variables *./deploy/.env* from *./deploy/tpl.env*
* Cause of specific behavior of Google App Engine builder you need to copy all files from dir *$GOPATH/src/github.com/labstack/echo* to *$GOPATH/src/github.com/labstack/echo/v4* and also install dependencies with go modules disabled after by command `GO111MODULE=off go get -d ./...` 


## Deploy command format
Run command from project root dir `./gcp.sh (init|deploy|rollback) version`

### Parameters
* init - Creation of CloudSQL instance, test and deploy app with specified version to Google App Engine (GAE)
* deploy - Test and deploy specified version of app to current GCP infrastructure (CloudSQL and GAE)
* rollback - rollback app to sprecific or previous version
* version - Specified app version for deploy or rollback

### Metrics
Available at *URL/metrics*

### Supported DB
* PostgreSQL
* MySQL
* CloudSQL
* Sqlite3 (with small changes)