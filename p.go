package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type P struct {
	projectsDir string
    openCmd string
}

func getProjects(rootDir string) []string {
	files, err := ioutil.ReadDir(rootDir)
    var projects []string

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
            projects = append(projects, file.Name())
		}
	}

    return projects
}

func (p *P) ls() {
    projects := getProjects(p.projectsDir)

    for _, projectName := range projects {
        fmt.Println(projectName)
    }
}

func (p *P) new(name string) {
    path := filepath.Join(p.projectsDir, name)
    if err := os.Mkdir(path, os.ModePerm); err != nil {
        fmt.Printf("Could not create project directory! %v", err)
    }
}

func (p *P) open(name string) {
    project := filepath.Join(p.projectsDir, name)
    err := os.Chdir(project)

    if err != nil {
        fmt.Printf("Could not change to project directory. Error: %v", err)
    }

    cmdToRun := "vim"
    if p.openCmd != "" {
        cmdToRun = p.openCmd
    }

    cmd := exec.Command(cmdToRun)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout

    err = cmd.Run()
    if err != nil {
        fmt.Println(err)
    }
}

func readConfig(p *P) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(dirname + "/.config/p.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 || scanner.Text()[0:1] != "#" { // Ignore comments
			conf := strings.Split(scanner.Text(), "=")
			if conf[0] == "projects" {
				p.projectsDir = conf[1]
			} else if conf[0] == "open" {

                if strings.Contains(conf[1], "${projects}") {
                    // Make sure the projects setting is set, when referenced in 'open' setting.
                    if p.projectsDir == "" {
                        fmt.Println("Failed to read config!")
                        fmt.Println("Make sure the 'projects' setting is set in your p.conf, and that it comes before the 'open' setting")
                        return
                    }
                }

                p.openCmd = conf[1]
            }
		}
	}
}

func main() {
	p := P{}
	readConfig(&p)
	numArgs := len(os.Args[1:])
    args := os.Args[1:]

	if numArgs <= 0 {
		fmt.Println("Common usage: ls")
        return
	}

    switch args[0] {
    case "ls":
        p.ls()
    case "new":
        if numArgs < 2 {
            fmt.Println("'new' command requires a project name")
            return
        }
        p.new(args[1])
    case "open":
        if numArgs < 2 {
            fmt.Println("'open' command requires a project name")
            return
        }
        p.open(args[1])
    default:
        fmt.Printf("'%v' command not supported\n", args[0])
    }
}
