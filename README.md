# rt-app
Some basic REST app with deploy to GCP(GAE)

# System Diagram

![alt text](https://github.com/artemantipov/rt-app/blob/master/diagram.png)

# Deploy to GCP

## Requirements
* gcloud sdk 
* golang 1.11+

## Preparations
Create file and fill with ENV variables *./deploy/.env* from *./deploy/tpl.env*

## Deploy command format
Run command from project root dir `./gcp.sh (init|deploy) version`

### Parameters
* init - Creation of CloudSQL instance, test and deploy app with specified version to Google App Engine (GAE)
* deploy - Test and deploy specified version of app to current GCP infrastructure (CloudSQL and GAE)
* version - Specified app version for deploy

