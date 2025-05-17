package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ProjectName string                    `yaml:"project_name"`
	Services    map[string]ServiceConfigs `yaml:"services"`
	// Order       []string                  `yaml:"order"`
}

type ServiceConfigs struct {
	Command  string `yaml:"command"`
	Terminal bool   `yaml:"terminal"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: booty <project-name> <project-config.yaml>")
		os.Exit(1)
	}
	project_name := os.Args[1]
	project_config := os.Args[2]
	fmt.Printf("Running project: %s\n", project_name)
	config, err := os.ReadFile(project_config)
	if err != nil {
		fmt.Println("Error reading config!")
	}
	var config_loader Config
	err = yaml.Unmarshal(config, &config_loader)
	if err != nil {
		fmt.Printf("Error unmarshalling yaml file!\n")
	}
	for name, services := range config_loader.Services {
		fmt.Printf("Running service: %s\n", name)
		if runtime.GOOS == "linux" {
			if services.Terminal {
				exec.Command("gnome-terminal", "--", "bash", "-c", services.Command).Start()
			} else {
				exec.Command("bash", "-c", services.Command).Run()
			}
		} else if runtime.GOOS == "windows" {
			if services.Terminal {
				cmd := fmt.Sprintf("\"%s\"", services.Command)
				exec.Command("cmd", "/C", "start", cmd).Start()
			} else {
				exec.Command("cmd", "/C", services.Command).Run()
			}
		} else {
			fmt.Printf("%s is not supported yet!\n", runtime.GOOS)
		}
	}
}
