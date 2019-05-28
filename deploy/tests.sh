#!/bin/bash
. ./deploy/.env

APP_URL=$1
CHECK_DATE="2006-01-02"

# GET-TEST
curl -i -X GET $APP_URL | grep "HTTP/2 200"
[[ $? != 0 ]] && echo "GET-TEST FAILED!" && exit 1
echo "GET-TEST PASSED!"

# PUT-DATA-TEST
curl -i -X PUT \
   -H "Content-Type:application/json" \
   -d \
'{ 
  "dateOfBirth": "2006-01-02"
}' \
 $APP_URL/hello/test | grep "HTTP/2 204"
 [[ $? != 0 ]] && echo "PUT-DATA-TEST FAILED!" && exit 1
 echo "PUT_DATA_TEST PASSED!"

# HELLO-USER-TEST
curl -i -X GET $APP_URL/hello/test | grep "Hello test!"
[[ $? != 0 ]] && echo "HELLO-USER-TEST FAILED!" && exit 1
echo "HELLO-USER-TEST PASSED!"
