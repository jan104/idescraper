#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push jan104/idescraper

if [[ -n $1 ]]; then
    docker push jan104/idescraper:$1
fi
