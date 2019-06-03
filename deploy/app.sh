#!/bin/bash
. ./deploy/.env

ACTION=$1
VERSION=$2
GAE_PROJECT_ID=$(gcloud app describe | grep id: | awk '{print $2}')

# Generate app.yaml for GAE from template
./deploy/appcfg.sh $GAE_PROJECT_ID

case $ACTION in
    "init")
        gcloud app create --region="$GCP_REGION"
        gcloud app deploy -q -v $VERSION --promote || exit 1
        echo "Use URL below for check"
        APP_URL=$(gcloud app browse --no-launch-browser -s default 2>&1)
        echo "$APP_URL"
        ;;
    "deploy")
        ./deploy/appcfg.sh $GAE_PROJECT_ID TEST # TEST arg generate config for testdb
        gcloud app deploy -q -v $VERSION --no-promote || exit 1
        echo "Use URL below for check"
        APP_URL=$(gcloud app browse --no-launch-browser -s default -v $VERSION 2>&1)
        echo "$APP_URL"
        echo "Starting tests..."
        ./deploy/tests.sh $APP_URL
        echo "Tests Passed!"
        echo "Deploying version $VERSION..."
        ./deploy/appcfg.sh $GAE_PROJECT_ID
        gcloud app deploy -q -v $VERSION --promote || exit 1
        echo "App deployed! Use URL below for check"
        APP_URL=$(gcloud app browse --no-launch-browser -s default 2>&1)
        echo "$APP_URL"
        ;;
    "rollback")
        echo "Starting rollback..."
        gcloud app deploy -q -v $VERSION --promote || exit 1
        echo "Rollback complete"
        ;;
esac


