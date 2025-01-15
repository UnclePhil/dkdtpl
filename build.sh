!#/bin/bash
mkdir ./bin
mv .git .gitold
docker context use default
# build 
docker run -it --rm  -v ${PWD}:/data -w /data golang:latest go build -v -o ./bin/dkdtpl   
sudo chown -R $(id -u):$(id -g) ./bin 
mv .gitold .git

## run test
cd test 
../bin/dkdtpl stack.yml stack.tmp
cat stack.tmp
cd ..
