#!/bin/sh
echo shut down existed docker service
echo you env is $1
if [ $1 == "TEST" ]
then
    export RUNTIME="test"
    docker stop service_ws_test
    docker stop service_http_test
    docker stop cache_test
    docker stop mongo_test

    docker container rm service_ws_test
    docker container rm service_http_test
    docker container rm cache_test
    docker container rm mongo_test

    docker rmi test_neofura_http -f
    docker rmi test_neofura_ws -f
    docker-compose -p "test" up -d
fi

if [ $1 == "STAGING" ]
then
    export RUNTIME="staging"
    docker stop service_ws_staging
    docker stop service_http_staging
    docker stop cache_staging
    docker stop mongo_staging

    docker container rm service_ws_staging
    docker container rm service_http_staging
    docker container rm cache_staging
    docker container rm mongo_staging

    docker rmi service_ws_staging -f
    docker rmi service_http_staging -f
    docker-compose -p "staging" up -d
fi


