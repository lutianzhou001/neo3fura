#!/bin/sh
echo shut down existed docker service
docker-compose down
echo remove images
docker kill $(docker ps -q)
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q) -f 
echo restart docker service
docker-compose -p "TEST" up -d

