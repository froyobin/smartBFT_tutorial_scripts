#!/bin/bash
docker rm -f $(docker ps -aq)
docker network prune
docker rmi dev-jdoe-mycc-1.2-e2b81d2d57d8667b9c094c15f4a726966eb6c659dd00e52f8c505432af98d4bf
