package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/RTradeLtd/thc"
)

var (
	cidInputFile = flag.String("file.name", "", "specify input file for scripting")
	user         = flag.String("user.name", "", "specify user login name")
	pass         = flag.String("user.pass", "", "specify user login password")
	dev          = flag.Bool("dev", true, "indicate dev or prod api")
	mode         = flag.String("mode", "index", "specify run mode")
	url          string
)

func init() {
	flag.Parse()
	if *user == "" {
		log.Fatal("user.name must not be empty")
	}
	if *pass == "" {
		log.Fatal("user.pass must not be empty")
	}
	if *dev {
		url = thc.DevURL
	} else {
		url = thc.ProdURL
	}
}

func main() {
	v2 := thc.NewV2(*user, *pass, url)
	if err := v2.Login(); err != nil {
		log.Fatal(err)
	}
	if *mode == "index" {
		if *cidInputFile == "" {
			log.Fatal("file.name cant be empty")
		}
		cids, err := readFile(*cidInputFile)
		if err != nil {
			log.Fatal(err)
		}
		for _, cid := range cids {
			if resp, err := v2.IndexHash(cid, true); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(resp)
			}
		}
	}
}

// readFile is used to read the content of the file into an array
func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
