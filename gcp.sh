#!/bin/bash
. ./deploy/.env

ACTION=$1
VERSION=$2

[[ -z $ACTION ]] && echo "Not enough arguments! e.g. './gcp.sh init v1' or './gcp.sh deploy v2'" && exit 1
[[ $# -ne 2 ]] && echo "Please provide app version as second arg." && exit 1

# Check for gcloud sdk install and env file exist
echo "Checking gcloud SDK..."
gcloud version || (echo "gcloud sdk required" && exit 1)
test -f ./deploy/.env || (echo "Please create ENV file ./deploy/.env from ./deploy/tpl.env" && exit 1)


case $ACTION in
    "init")
        echo "Imitiatiom started: Create SQL instance, deploy and test app"
        ./deploy/db.sh
        ./deploy/app.sh $ACTION $VERSION
        ;;
    "deploy")
        echo "Deployment started: Deploy and test new version ($VERSION) of app "
        ./deploy/app.sh $ACTION $VERSION
        ;;
    *)
        echo "1st arg must be 'init' or 'deploy'"
        exit 1
        ;;
esac
        



