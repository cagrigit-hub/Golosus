package main

import (
	"fmt"
	"log"
	"os"
)

type file struct {
	name    string
	content string
}

func main() {
	// open CLI menu
	openCLIMenu()
	// get option
	var option int
	fmt.Scanln(&option)

	var name string
	var githubProfile string
	switch option {
	case 1:
		fmt.Println("Enter project name: ")
		fmt.Scanln(&name)
		fmt.Println("Enter your github profile: ")
		fmt.Scanln(&githubProfile)
		fmt.Println("Creating project...")
	case 2:
		fmt.Println("Exiting Golosus-Web CLI")
		return
	default:
		fmt.Println("Invalid option")
	}

	ct := &Content{}

	folders := []string{
		"assets",
		"assets/jscode",
		"assets/bundled",
		"cmd",
		"view",
		"model",
		"handler",
		"view/components",
		"view/layout",
		"view/example",
		"typescript",
		".",
	}

	for _, folder := range folders {
		err := createFolders(name + "/" + folder)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	files := map[string][]file{
		"cmd": {
			{"main.go", ct.Main(name, githubProfile)},
		},
		"view/layout": {
			{"base.templ", ct.Layout(name)},
		},
		"view/example": {
			{"example.templ", ct.ExampleView(githubProfile, name)},
		},
		"view/components": {
			{"input.templ", ct.ExampleComponent()},
		},
		"handler": {
			{"util.go", ct.Util()},
			{"example.go", ct.ExampleHandler(githubProfile, name)},
		},
		"model": {
			{"example.go", ct.ExampleModel()},
		},
		".": {
			{"go.mod", ct.GoMod(githubProfile, name)},
			{"Makefile", ct.Make()},
			{".air.toml", ct.Air()},
		},
		"typescript": {
			{"tsconfig.json", ct.TsConfig()},
			{"index.ts", ct.TypescriptIndex()},
			{"scripts.ts", ct.ScriptsTs()},
			{"package.json", ct.PackageJson(name)},
		},
	}

	for _, folder := range folders {
		for _, file := range files[folder] {
			err := createFiles(name, folder, file.name)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = writeFiles(name, folder, file.name, file.content)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

	}
}

func createFolders(name string) error {
	if err := os.MkdirAll(name, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func createFiles(rn, target, name string) error {
	if _, err := os.Create(rn + "/" + target + "/" + name); err != nil {
		return err
	}
	return nil
}

func writeFiles(rn, target, name, content string) error {
	if err := os.WriteFile(rn+"/"+target+"/"+name, []byte(content), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func openCLIMenu() {
	fmt.Println("Welcome to Golosus-Web CLI")
	fmt.Println("Please select an option")
	fmt.Println("1. Create a new project")
	fmt.Println("2. Exit")
}
