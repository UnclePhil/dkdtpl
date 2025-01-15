!#/bin/bash
mv .git .gitold
docker context use default
# docker run -it --rm --user $(id -u):$(id -g) -v ${PWD}:/data -w /data golang:latest go build -v -o dkdtpl   
docker run -it --rm  -v ${PWD}:/data -w /data golang:latest go build -v -o dkdtpl   
mv .gitold .git