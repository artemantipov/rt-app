# rt-app
Some basic test app with deploy to GCP

# System Diagram

![alt text](https://github.com/artemantipov/rt-app/blob/master/diagram.png)

# Deploy to GCP

## Requirements
* gcloud sdk
* golang 1.11+

## Deploy command format
Run from project root dir:
* ./gcp.sh (init|deploy) version

init - Creation of CloudSQL instance, test and deploy app with specified version to Google App Engine (GAE)
deploy - Test and deploy specified version of app to current GCP infrastructure (CloudSQL and GAE)

