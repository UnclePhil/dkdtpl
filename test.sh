## run test
echo "Test dkdtpl v1.0.0"
echo "=================="
cd test 
rm stack.tmp
## ../bin/dkdtpl stack.yml stack.tmp 
../bin/dkdtpl stack.yml stack.tmp rclone remote=s3fs:${CLUSTER_NAME}/${STNAME}/data,allow_other=true,vfs_cache_mode=full,poll_interval=0
cat stack.tmp
cd ..
