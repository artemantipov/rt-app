#!/bin/bash
. ./deploy/.env

gcloud sql instances create "$CSQL_INSTANCE_NAME" --region="$GCP_REGION" --tier db-f1-micro --database-version=MYSQL_5_7 --async
gcloud sql instances describe "$CSQL_INSTANCE_NAME" | grep RUNNABLE
while [[ $? != 0 ]]
do 
    echo "CloudSQL instance deployment in progress, please wait..."
    sleep 10
    gcloud sql instances describe "$CSQL_INSTANCE_NAME" | grep RUNNABLE
done
gcloud sql databases create "$DB_NAME" --instance="$CSQL_INSTANCE_NAME" # prod db
gcloud sql databases create "$DB_NAME"-test --instance="$CSQL_INSTANCE_NAME" # test db
gcloud sql users create "$DB_USER" -i "$CSQL_INSTANCE_NAME" --password="$DB_PASS"



