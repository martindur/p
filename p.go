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
    openFromProject bool
    git bool
    readme bool
}

// func createPattern(pattern string) {
//     switch pattern {
//     case ".git":
//         cmd := exec.Command("git", "init")
//         cmd.Run()
//     case ".gitignore":
//         os.NewFile
//     }
// }

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

func resolveOpen(cmd string, project string) string {
    return strings.ReplaceAll(cmd, "${PROJECT}", project)
}

func (p *P) ls() {
    projects := getProjects(p.projectsDir)

    for _, projectName := range projects {
        fmt.Println(projectName)
    }
}

func (p *P) new(name string) error {
    path := filepath.Join(p.projectsDir, name)
    if err := os.Mkdir(path, os.ModePerm); err != nil {
        return err
    }

    curDir, err := os.Getwd()

    if p.git {
        if err != nil {
            return err
        }

        err = os.Chdir(path)
        if err != nil {
            return err
        }

        cmd := exec.Command("git", "init")
        cmd.Run()

        err = os.Chdir(curDir)
        if err != nil {
            return err
        }
    }

    if p.readme {
        readmeFilePath := filepath.Join(p.projectsDir, name, "readme.md")
        _, err := os.Create(readmeFilePath)
        if err != nil {
            return err
        }
    }

    return nil
}

func (p *P) open(name string) {
    project := filepath.Join(p.projectsDir, name)

    var argsToRun []string
    var err error
    if p.openCmd != "" {
        argsToRun = strings.Split(resolveOpen(p.openCmd, project), " ")
    } else {
        argsToRun = append(argsToRun, "vim")
    }

    if p.openFromProject {
        err = os.Chdir(project)
        if err != nil {
            fmt.Printf("Could not change to project directory. Error: %e", err)
            return
        }
    }

    // First arg is 'some' named executable. The rest are unpacked as arguments
    var cmd *exec.Cmd
    if len(argsToRun) > 1 {
        cmd = exec.Command(argsToRun[0], argsToRun[1:]...)
    } else {
        cmd = exec.Command(argsToRun[0])
    }
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

            switch conf[0] {
            case "projects":
                p.projectsDir = conf[1]
            case "open":
                if strings.Contains(conf[1], "${projects}") {
                    // Make sure the projects setting is set, when referenced in 'open' setting.
                    if p.projectsDir == "" {
                        fmt.Println("Failed to read config!")
                        fmt.Println("Make sure the 'projects' setting is set in your p.conf, and that it comes before the 'open' setting")
                        return
                    }
                }
                p.openCmd = strings.ReplaceAll(conf[1], "\"", "")
            case "open_from_project":
                if conf[1] == "true" {
                    p.openFromProject = true
                } else if conf[1] == "false" {
                    p.openFromProject = false
                } else {
                    fmt.Println("Failed to read config! the 'open_from_project' setting only supports the values 'true' and 'false'")
                }
            case "git":
                if conf[1] == "true" {
                    p.git = true
                } else if conf[1] == "false" {
                    p.git = false
                } else {
                    fmt.Println("Failed to read config! the 'git' setting only supports the values 'true' and 'false'")
                }
            case "readme":
                if conf[1] == "true" {
                    p.readme = true
                } else if conf[1] == "false" {
                    p.readme = false
                } else {
                    fmt.Println("Failed to read config! the 'readme' setting only supports the values 'true' and 'false'")
                }
            }
		}
	}
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}


func main() {
    p := P{projectsDir: "", openCmd: "", openFromProject: false, git: false, readme: false}
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
        err := p.new(args[1])
        if err != nil {
            fmt.Printf("Could not create new project! %e", err)
        }
    case "open":
        if numArgs < 2 {
            fmt.Println("'open' command requires a project name")
            return
        }
        p.open(args[1])
    default:
        if contains(getProjects(p.projectsDir), args[0]) {
            p.open(args[0])
        }
        fmt.Printf("'%v' command not supported\n", args[0])
    }
}
