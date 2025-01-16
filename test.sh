## run test
cd test 
../bin/dkdtpl stack.yml stack.tmp "rclone" "remote='s3fs:${CLUSTER_NAME}/data',allow_other='true',vfs_cache_mode=full,poll_interval=0"
cat stack.tmp
cd ..
