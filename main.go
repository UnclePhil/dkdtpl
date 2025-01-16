package main

import (
        "fmt"
        "io/ioutil"
        "log"
        "os"
        // "regexp"
        "strings"
        "crypto/md5"
        "gopkg.in/yaml.v3"
)

type Volume struct {
        Driver     string            `yaml:"driver,omitempty"`
        DriverOpts []string          `yaml:"driver_opts,omitempty"`
}

type Volumes map[string]Volume

type Config struct {
	File     string          `yaml:"file,omitempty"`
	Name     string          `yaml:"name,omitempty"`
}

type Configs map[string]Config

type Secret struct {
	File     string          `yaml:"file,omitempty"`
	Name     string          `yaml:"name,omitempty"`
}

type Secrets map[string]Secret


type DockerCompose struct {
        Version string   `yaml:"version"`
        Services map[string]interface{} `yaml:"services"`
        Networks map[string]interface{} `yaml:"networks"`
        Volumes  Volumes  `yaml:"volumes"`
		Configs  Configs  `yaml:"configs"`
		Secrets  Secrets  `yaml:"secrets"`
}




func main() {
	// get param
	if len(os.Args) < 2 {
		fmt.Println("Usage: dkdtpl <INfilename> <OUTfilename>")
		os.Exit(1)
	}
	/////// Args
	// 1: Infilename  docker-stack.yml original  mandatory
	// 2: OUTfilename : the output file  mandatory
	// 3: Volume driver : default: local   any other value  
	// 4: Volume driver options : all options separate by comma    
	// ex: 
	//			"remote='s3fs:${CLUSTER_NAME}/data'",
	//			"allow_other='true'",
	//			"vfs_cache_mode=full",
	//			"poll_interval=0",
	// ================================================================
    // x

	var DKDdriver string
	var DKDdriveropts []string
	var INfilename string
	var OUTfilename string

	INfilename = os.Args[1]
	OUTfilename = os.Args[2]

	if os.Args[3] == "" { 
		DKDdriver = "local"
	} else {
		DKDdriver = os.Args[3]
	} 

	if os.Args[4] == "" { 
		DKDdriveropts = []string
	} else {
		DKDdriveropts = strings.Split(os.Args[4], ",")
	} 


	// Read the Docker Compose file
	yamlFile, err := ioutil.ReadFile(INfilename)
	if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the YAML
	var data DockerCompose
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
			log.Fatalf("Error parsing YAML: %v", err)
	}

	//////////////////////////////////////////////////////////////////////
	// Iterate through volumes
	for volumeName, volume := range data.Volumes {
		if volume.Driver == "" {
			// Set default driver

			vol := data.Volumes[volumeName]
			vol.Driver = DKDdriver
			vol.DriverOpts = DKDdriveropts

// 			vol.DriverOpts = []string{
//				"remote='s3fs:${CLUSTER_NAME}/data'",
//				"allow_other='true'",
//				"vfs_cache_mode=full",
//				"poll_interval=0",
//			}

			data.Volumes[volumeName] = vol

		}
	}
	
	//////////////////////////////////////////////////////////////////////	
	// Iterate through configs
	for configName := range data.Configs {
		filename := data.Configs[configName].File
		cfg := data.Configs[configName]
		fdata, err := ioutil.ReadFile(filename)
		if err != nil {
				log.Fatal(err)
		}
		hash := md5.Sum(fdata)
		md5sum := fmt.Sprintf("%x", hash)
		cfg.Name = configName+"."+md5sum
		data.Configs[configName] = cfg
	}

	//////////////////////////////////////////////////////////////////////	
	// Iterate through configs
	for secretName := range data.Secrets {
		filename := data.Secrets[secretName].File
		scr := data.Secrets[secretName]
		fdata, err := ioutil.ReadFile(filename)
		if err != nil {
				log.Fatal(err)
		}
		hash := md5.Sum(fdata)
		md5sum := fmt.Sprintf("%x", hash)
		scr.Name = secretName+"."+md5sum
		data.Secrets[secretName] = scr
	}



	// Marshal the modified data back to YAML
	modifiedYaml, err := yaml.Marshal(data)
	if err != nil {
			log.Fatalf("Error marshalling YAML: %v", err)
	}

	// Write the modified YAML to a new file
	err = ioutil.WriteFile(OUTfilename, modifiedYaml, 0644)
	if err != nil {
			log.Fatalf("Error writing YAML file: %v", err)
	}

}