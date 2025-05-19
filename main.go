package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

var BootyConfigDir string

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Couldnt find user's config dir!") // Handle error

	}
	BootyConfigDir = filepath.Join(configDir, "booty")
}

func initialize() {
	// can also add metadata.json
	fmt.Println("Initializing booty")

	fmt.Printf("Default config path: %s\n", BootyConfigDir)
	err := os.MkdirAll(BootyConfigDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	fmt.Printf("Initializing project\nPlease Enter project's name: ")
	var projectName string
	fmt.Scanln(&projectName)

	projectConfigPath := fmt.Sprintf("%s.yaml", filepath.Join(BootyConfigDir, projectName))
	fmt.Printf("Creating config file: %s\n", projectConfigPath)
	if FileExists(projectConfigPath) {
		fmt.Printf("Project exists! Do you wish to overwrite(y/n): ")
		var choice string
		fmt.Scanln(&choice)
		if choice != "y" {
			os.Exit(1)
		}
	}
	projectDir, err := os.Getwd()
	if err != nil {
		projectDir = "" // Handle error
		log.Println(err.Error())
	}
	// change template
	templateContent := fmt.Sprintf(`

# Booty Project Configuration

# General Settings
project_name: %s
workdir: %s
services: 
   # service-name: # Customize name
    # command: "docker run -d -p 8000:8000 container" # Configure commands
    # terminal: false # False for Detatched mode

  # Add more services as per requirement
  # backend:
    # command: "cd backend/cmd && go run main.go"
    # terminal: true

`, projectName, projectDir)

	err = os.WriteFile(projectConfigPath, []byte(templateContent), 0644)
	if err != nil {
		log.Fatalf("Error writing template to file: %v", err)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}
	if runtime.GOOS == "linux" {
		exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("%s %s", editor, projectConfigPath)).Start()
	} else if runtime.GOOS == "windows" {
		exec.Command("notepad", projectConfigPath).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "%s %s"`, editor, projectConfigPath)).Start()
	}
}

func boot() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: booty boot <project-name>")
		os.Exit(1)
	}
	project_name := os.Args[2]
	project_config := fmt.Sprintf("%s.yaml", filepath.Join(BootyConfigDir, project_name))
	fmt.Printf("Running project: %s\n", project_name)
	config, err := os.ReadFile(project_config)
	if err != nil {
		fmt.Println("Project does not exist")
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
		} else if runtime.GOOS == "darwin" {
			if services.Terminal {
				exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Terminal\" to do script \"%s\"", services.Command)).Start()
			} else {
				exec.Command("bash", "-c", services.Command).Run()
			}
		} else {
			fmt.Printf("%s is not supported yet!\n", runtime.GOOS)
		}
	}

}

func list() {
	fmt.Println("Found following projects: ")
	matches, err := filepath.Glob(filepath.Join(BootyConfigDir, "*.yaml"))
	if err != nil {
		log.Fatalf("Error reading the config dir!") // Handle error
	}
	for i := range matches {
		projName := strings.Split(matches[i], "/")
		fmt.Printf("%s\t", projName[len(projName)-1][0:len(projName[len(projName)-1])-5])
	}
	fmt.Println()
}

func config() {
	var editor string
	if os.Getenv("EDITOR") != "" {
		editor = os.Getenv("EDITOR")
	} else {
		editor = "nano"
	}
	projectName := os.Args[2]
	path := filepath.Join(BootyConfigDir, projectName)
	configPath := fmt.Sprintf("%s.yaml", path)
	if runtime.GOOS == "linux" {
		exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("%s %s", editor, configPath)).Start()
	} else if runtime.GOOS == "windows" {
		exec.Command("notepad", configPath).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "%s %s"`, editor, configPath)).Start()
	}
}

func remove() {
	for {
		list()
		fmt.Print("Enter the project name to remove: ")
		var projectName string
		fmt.Scanln(&projectName)
		path := filepath.Join(BootyConfigDir, projectName)
		matches, err := filepath.Glob(fmt.Sprintf("%s.yaml", path))
		if err != nil {
			log.Fatalf("Error searching for file!") // Handle error
		}

		if matches != nil {
			err = os.Remove(matches[0])
			if err != nil {
				fmt.Println("The project could not be deleted") // Handle error
				os.Exit(1)
			}
			fmt.Println("Project was successfully deleted.")
			return
		}

		fmt.Printf("Project does not exist, do you wish to retry(y/n): ")
		var choice string
		fmt.Scanln(&choice)
		if choice != "y" {
			return
		} else {
			continue
		}
	}
}

func main() {
	// fetaures to be implemented
	// help, edit, remove, add(service), status

	// need to implement auto yaml generation, yaml edits and then config edits without yaml
	// kill proccesses on exit
	if len(os.Args) < 2 {
		fmt.Printf("Usage: booty <options> COMMAND\n")
	}
	feature := os.Args[1]
	if feature == "boot" {
		// implement workdir based cd and source
		boot()
	} else if feature == "help" {
		fmt.Printf("Usage: booty <options> COMMAND\n")
	} else if feature == "init" {
		initialize()
	} else if feature == "list" {
		list()
	} else if feature == "config" {
		config()
	} else if feature == "remove" {
		remove()
	}
}
