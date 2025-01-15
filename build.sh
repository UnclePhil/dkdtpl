!#/bin/bash
docker context use default
docker run -it --rm -v ${PWD}:/data -w /data golang:latest go build -v -o dkdtpl   