package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type P struct {
	projectsDir string
}

func (p *P) ls() {
	files, err := ioutil.ReadDir(p.projectsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name())
		}
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
		if scanner.Text()[0:1] != "#" { // Ignore comments
			conf := strings.Split(scanner.Text(), "=")
			if conf[0] == "projects" {
				p.projectsDir = conf[1]
			}
		}
	}
}

func main() {
	p := P{}
	readConfig(&p)
	numArgs := len(os.Args[1:])

	if numArgs <= 0 {
		fmt.Println("Common usage: ls")
	} else {
		p.ls()
	}
}
