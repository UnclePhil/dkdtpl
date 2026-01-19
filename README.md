# DKDTPL

- Simplify the config and secret rotation in docker SwarmMode
- fill default values for volumes driver and options 


BTW: This is not perfect ... like me ...  but it works !

## Time line
20260119

- public repo
- extend documentation

20250116

- initial commit

## Secret/config rotation
If you are using docker swarm, you know that secret and config rotation is a s...t.
You need to give a new name to the config and change it in the stack file .

But there is a trick !!

Config ans secret can have a name !!

```
configs:
    bookmarks.yaml:
        file: ./config/bookmarks.yaml
        name: bookmarks.yaml.5b8497360c6aa8ae96e30b4a63257a51
```
- And this name is used in the effective docker configuration
- And a name can be modified with a md5 (or any) tail.

So if , before the load of the stack, you can evaluate the md5 of the config and add the correct name in a temporary file, you do not need to modify the stack at each load.

**That's the first purpose of dkdtpl**

```
dkdtpl stack.yml tmpstack.yml  
```
Give you something like the example above for each config or secrets

Now,  
```
docker stack deploy -c tmpstack.yml mystackname
``` 
will (re)load the ad-hoc config like a breeze !

## Volume driver
If, like me , you are using multiple docker swarm with multiple storages configuration. But don't like to change or multiply stack file on purpose.
You need a solution to inject driver and options on purpose

**That's the other purpose of DKDTPL**
```
dkdtpl: with <3 by UnclePhil - 1.0.0
    Usage: dkdtpl <INfilename> <OUTfilename> [voldriver  [voldriveropts]]
```

Set some variables for the specific cluster
```
VOLDRIVER=rclone 
VOLDRIVEROPTS=remote=s3fs:dcltst01/$STNAME/$VOLNAME,allow_other=true,vfs_cache_mode=full,poll_interval=0
```
Then run 
```
DKDTPL stack.yml tmpstack.yml $VOLDRIVER $VOLDRIVEROPTS
```

Any volume definition in stack.yml without driver definition , will be filled with the ad-hoc VOLDRIVER and VOLOPTIONS 

