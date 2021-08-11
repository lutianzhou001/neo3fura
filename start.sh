#!/bin/sh
echo shut down existed docker service
docker-compose down
echo remove image of neo3fura_web
docker rmi neo3fura_web
echo restart docker service
docker-compose up -d

