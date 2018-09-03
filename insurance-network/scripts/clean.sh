#!/usr/bin/env bash

docker stop $(docker ps -q)

docker rm $(docker ps -aq)

#docker rmi $(docker images -aq)

rm -rf ../channel-artifacts/*

rm -rf ../crypto-config/*