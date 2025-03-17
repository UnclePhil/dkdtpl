package main

import (
        "flag"
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
        DriverOpts map[string]string          `yaml:"driver_opts,omitempty"`
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

// Modifions d'abord la structure pour supporter les services avec deploy
type Deploy struct {
        Labels map[string]string `yaml:"labels,omitempty"`
}

type Service struct {
        Deploy     *Deploy                `yaml:"deploy,omitempty"`
        Image      string                 `yaml:"image,omitempty"`
        Command    interface{}            `yaml:"command,omitempty"`
        Ports      interface{}            `yaml:"ports,omitempty"`
        Environment interface{}           `yaml:"environment,omitempty"`
        Volumes    interface{}            `yaml:"volumes,omitempty"`
        Networks   interface{}            `yaml:"networks,omitempty"`
        Configs    interface{}            `yaml:"configs,omitempty"`
        Secretss   interface{}            `yaml:"secretss,omitempty"`

}

type DockerCompose struct {
        Version  string                 `yaml:"version"`
        Services map[string]interface{} `yaml:"services"`
        Networks map[string]interface{} `yaml:"networks"`
        Volumes  Volumes               `yaml:"volumes"`
        Configs  Configs               `yaml:"configs"`
        Secrets  Secrets               `yaml:"secrets"`
}

func main() {
        // Version information
        const version = "2.0.0"
        
        // Define command line flags
        inFile := flag.String("in", "", "Input docker-compose file (required)")
        outFile := flag.String("out", "", "Output file (required)")
        volDriver := flag.String("driver", "local", "Volume driver (default: local)")
        volOpts := flag.String("opts", "", "Volume driver options (comma-separated key=value pairs)")
        volLabels := flag.String("labels", "", "Service labels (comma-separated key=value pairs)")
        showVersion := flag.Bool("version", false, "Show version information")
        
        // Custom usage message
        flag.Usage = func() {
                fmt.Fprintf(os.Stderr, "dkdtpl: Docker compose template processor - v%s\n\n", version)
                fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
                fmt.Fprintf(os.Stderr, "Options:\n")
                flag.PrintDefaults()
                fmt.Fprintf(os.Stderr, "\nExamples:\n")
                fmt.Fprintf(os.Stderr, "  %s -in docker-stack.yml -out output.yml -driver s3fs -opts \"remote=s3fs:cluster/data,allow_other=true\"\n", os.Args[0])
                fmt.Fprintf(os.Stderr, "  %s -in docker-stack.yml -out output.yml -labels \"environment=prod,owner=team1;monitoring=true,type=backend\"\n", os.Args[0])
        }

        flag.Parse()

        // Show version if requested
        if *showVersion {
                fmt.Printf("dkdtpl: with <3 by UnclePhil - v%s\n", version)
                os.Exit(0)
        }

        // Validate required flags
        if *inFile == "" || *outFile == "" {
                fmt.Fprintf(os.Stderr, "Error: Input and output files are required\n\n")
                flag.Usage()
                os.Exit(1)
        }

        // Parse driver options
        driverOpts := make(map[string]string)
        if *volOpts != "" {
                items := strings.Split(*volOpts, ",")
                for _, item := range items {
                        parts := strings.Split(item, "=")
                        if len(parts) == 2 {
                                driverOpts[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
                        }
                }
        }

        // Parse labels - modification pour supporter plusieurs groupes de labels
        labelsGroups := make([]map[string]string, 0)
        if *volLabels != "" {
                // Séparer les différents groupes de labels (séparés par des points-virgules)
                groups := strings.Split(*volLabels, ";")
                for _, group := range groups {
                        labels := make(map[string]string)
                        items := strings.Split(group, ",")
                        for _, item := range items {
                                parts := strings.Split(item, "=")
                                if len(parts) == 2 {
                                        labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
                                }
                        }
                        if len(labels) > 0 {
                                labelsGroups = append(labelsGroups, labels)
                        }
                }
        }

        // Read the Docker Compose file
        yamlFile, err := ioutil.ReadFile(*inFile)
        if err != nil {
                log.Fatalf("Error reading YAML file: %v", err)
        }

        // Parse the YAML
        var data DockerCompose
        err = yaml.Unmarshal(yamlFile, &data)
        if err != nil {
                log.Fatalf("Error parsing YAML: %v", err)
        }

        // Process volumes
        for volumeName, volume := range data.Volumes {
                if volume.Driver == "" {
                        vol := data.Volumes[volumeName]
                        vol.Driver = *volDriver
                        vol.DriverOpts = driverOpts
                        data.Volumes[volumeName] = vol
                }
        }
        
        // Iterate through configs
        //   add name with md5 value
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

        // Iterate through secrets
        //   add name with md5 value
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

        // Process services to add labels
        for serviceName, service := range data.Services {
                // Convert interface{} to map
                serviceMap, ok := service.(map[string]interface{})
                if !ok {
                        continue
                }

                // Check if deploy section exists
                deploy, exists := serviceMap["deploy"]
                if !exists {
                        // Only create deploy section if we have labels to add
                        if len(labelsGroups) > 0 {
                                deploy = make(map[string]interface{})
                                serviceMap["deploy"] = deploy
                        } else {
                                continue
                        }
                }

                deployMap, ok := deploy.(map[string]interface{})
                if !ok {
                        deployMap = make(map[string]interface{})
                }

                // Get existing labels
                var existingLabels map[string]string
                if labels, exists := deployMap["labels"]; exists {
                        // Try to convert existing labels to map[string]string
                        switch v := labels.(type) {
                        case map[string]string:
                                existingLabels = v
                        case map[string]interface{}:
                                // Convert interface{} map to string map
                                existingLabels = make(map[string]string)
                                for k, val := range v {
                                        if str, ok := val.(string); ok {
                                                existingLabels[k] = str
                                        }
                                }
                        default:
                                existingLabels = make(map[string]string)
                        }
                } else {
                        existingLabels = make(map[string]string)
                }

                // Add new labels while preserving existing ones
                if len(labelsGroups) > 0 {
                        // Add all label groups to the service
                        for _, labels := range labelsGroups {
                                for k, v := range labels {
                                        // Only add if key doesn't exist
                                        if _, exists := existingLabels[k]; !exists {
                                                existingLabels[k] = v
                                        }
                                }
                        }

                        // Update the deploy section with merged labels
                        deployMap["labels"] = existingLabels
                        serviceMap["deploy"] = deployMap
                        data.Services[serviceName] = serviceMap
                }
        }

        // Marshal the modified data back to YAML
        modifiedYaml, err := yaml.Marshal(data)
        if err != nil {
                log.Fatalf("Error marshalling YAML: %v", err)
        }

        // Write the modified YAML to a new file
        err = ioutil.WriteFile(*outFile, modifiedYaml, 0644)
        if err != nil {
                log.Fatalf("Error writing YAML file: %v", err)
        }
}