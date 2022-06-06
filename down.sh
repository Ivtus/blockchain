#!/bin/bash
docker-compose -f order1.yaml down --volume --remove-orphans
docker-compose -f p0o1.yaml down --volume --remove-orphans
docker rm -f $(docker ps -aq)
docker volume prune
docker netowrk prune

