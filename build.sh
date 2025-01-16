!#/bin/bash
## cleaning
mkdir ./bin
mv .git .gitold
rm ./test/stack.tmp

## prepare
docker context use default
## build 
docker run -it --rm  -v ${PWD}:/data -w /data golang:latest go build -v -o ./bin/dkdtpl   
sudo chown -R $(id -u):$(id -g) ./bin 

## rollback 
mv .gitold .git

## run test
cd test 
../bin/dkdtpl stack.yml stack.tmp "toto" "toto1","toto2","toto3"
cat stack.tmp
cd ..
