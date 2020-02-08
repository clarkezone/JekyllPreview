package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	clone()
}

func clone() {
	os.Mkdir("/srv/jekyll/source", os.ModePerm)
	os.Mkdir("/app/_site", os.ModePerm)

	cmd := exec.Command("git", "clone", "https://github.com/clarkezone/clarkezone.github.io.git", ".")
	cmd.Dir = "/srv/jekyll/source"
	err := cmd.Run()

	if err != nil {
		fmt.Println("an error occurred.\n")
		log.Fatal(err)
	}

	//os.Chown("/srv/jekyll/source", 1000, 1000)

	cmd = exec.Command("git", "checkout", "acme-ngrok")
	cmd.Dir = "/srv/jekyll/source"
	err = cmd.Run()

	if err != nil {
		fmt.Println("an error occurred.\n")
		log.Fatal(err)
	}

	cmd = exec.Command("chown", "-R", "jekyll:jekyll", "/srv/jekyll/source")
	err = cmd.Run()

	if err != nil {
		fmt.Println("an error occurred.\n")
		log.Fatal(err)
	}

	cmd = exec.Command("chown", "-R", "jekyll:jekyll", "/app/_site")
	err = cmd.Run()

	if err != nil {
		fmt.Println("an error occurred.\n")
		log.Fatal(err)
	}
}
