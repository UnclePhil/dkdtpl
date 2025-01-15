package main

import (
        "fmt"
        "io/ioutil"
        "log"
        "os"
        // "regexp"
        // "strings"
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
		fmt.Println("Usage: dkdfill <INfilename> <OUTfilename>")
		os.Exit(1)
	}
	INfilename := os.Args[1]
	OUTfilename := os.Args[2]


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
					data.Volumes[volumeName].Driver = "rclone"

					// Set default driver options
					data.Volumes[volumeName].DriverOpts = []string{
							"remote='s3fs:${CLUSTER_NAME}/data'",
							"allow_other='true'",
							"vfs_cache_mode=full",
							"poll_interval=0",
					}
			}
	}
	
	//////////////////////////////////////////////////////////////////////	
	// Iterate through configs
	for configName := range data.Configs {
		filename := data.Configs[configName].File
		fdata, err := ioutil.ReadFile(filename)
		if err != nil {
				log.Fatal(err)
		}
		hash := md5.Sum(fdata)
		md5sum := fmt.Sprintf("%x", hash)
		data.Configs[configName].Name = filename+"."+md5sum
	}

	//////////////////////////////////////////////////////////////////////	
	// Iterate through configs
	for secretName := range data.Secrets {
		filename := data.Secrets[secretName].File
		fdata, err := ioutil.ReadFile(filename)
		if err != nil {
				log.Fatal(err)
		}
		hash := md5.Sum(fdata)
		md5sum := fmt.Sprintf("%x", hash)
		data.Secrets[secretName].Name = filename+"."+md5sum
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