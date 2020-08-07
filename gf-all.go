package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	gfPath, err := exec.LookPath("gf")
	if err != nil {
		log.Fatal("gf is not installed / not in your PATH")
	}
	fmt.Printf("gf is available at %s\n", gfPath)

	command := exec.Command(gfPath, "-list")

	var out bytes.Buffer

	command.Stdout = &out
	err = command.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.RemoveAll("gf-out")
	os.Mkdir("gf-out", 0755)

	patterns := strings.Split(out.String(), "\n")

	for _, pattern := range patterns {
		if len(pattern) == 0 {
			continue
		}

		executePattern(gfPath, pattern)
	}
}

func executePattern(gfPath string, pattern string) bool {
	fmt.Printf("Executing %s\n", pattern)
	command := exec.Command(gfPath, pattern)
	command.Wait()
	gfOutput, err := command.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	gfOutputString := string(gfOutput)

	// Some gf patterns are not recursive
	gfOutputString = strings.Replace(gfOutputString, "grep: .: Is a directory\n", "", -1)

	if len(gfOutputString) == 0 {
		return false
	}

	f, err := os.Create("gf-out/" + pattern + ".txt")
	if err != nil {
		fmt.Println(err)
		return false
	}

	l, err := f.WriteString(gfOutputString)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return false
	}

	fmt.Println(l, "bytes written successfully")

	return true
}
