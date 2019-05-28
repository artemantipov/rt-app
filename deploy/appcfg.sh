#!/bin/bash
. ./deploy/.env

TPL=./app.tpl.yaml
CFG=./app.yaml

GAE_PROJECT_ID=$1
TEST=$2
[[ -z $TEST ]] || DB_NAME=$DB_NAME-test # change db name for test db name

sed -e "s/host/$DB_HOST/" \
    -e "s/port/$DB_PORT/" \
    -e "s/db/$DB_NAME/" \
    -e "s/user/$DB_USER/" \
    -e "s/password/$DB_PASS/" \
    -e "s/type/$DB_TYPE/" \
    -e "s/PROJECT/$GAE_PROJECT_ID/" \
    -e "s/REGION/$GCP_REGION/" \
    -e "s/INSTANCE/$CSQL_INSTANCE_NAME/" \
    $TPL > $CFG