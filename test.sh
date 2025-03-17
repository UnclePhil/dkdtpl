## run test
echo "Test dkdtpl v2.0.0"
echo "=================="
cd test 
rm stack.tmp
## ../bin/dkdtpl stack.yml stack.tmp 
../bin/dkdtpl -in stack.yml -out stack.tmp -driver rclone -opts "remote=s3fs:${CLUSTER_NAME}/${STNAME}/data,allow_other=true,vfs_cache_mode=full,poll_interval=0" -labels "environment=prod,team=devops,owner=koenigphil"
cat stack.tmp
cd ..
