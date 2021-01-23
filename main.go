package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Vikash082/registrylogin/server"
	"gopkg.in/yaml.v3"
)

func main() {
	inputFile := flag.String("f", "", "inventory file for servers where docker/podman login need to be done")
	dockerUser := flag.String("du", "", "docker username")
	dockerPassword := flag.String("dp", "", "docker password")
	podmanUser := flag.String("pu", "", "podman user")
	podmanPassword := flag.String("pp", "", "podman password")
	flag.Usage = usage
	flag.Parse()
	validateArgs(*inputFile, *dockerUser, *dockerPassword, *podmanUser, *podmanPassword)

	input, err := ioutil.ReadFile(*inputFile)
	logAndExit(err)
	var servers []server.Server
	logAndExit(yaml.Unmarshal(input, &servers))
	for _, serv := range servers {
		cmdString := serv.GetLoginCommand(*dockerUser, *dockerPassword, *podmanUser, *podmanPassword)
		log.Println(cmdString)
		if cmdString == "" {
			log.Printf("Unable to determine container cli for server %s", serv.IP)
			continue
		}
		err := serv.Execute(cmdString)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("eg: registrylogin -f /path/to/<inventory.yaml> -du user -dp nopass")
	flag.PrintDefaults()
}

func validateArgs(fname, duser, dpasswd, puser, ppaswd string) {
	if fname == "" {
		printMsgAndExit(fmt.Sprint("Please provide input file"))
	}
	if duser != "" {
		if dpasswd == "" {
			printMsgAndExit(fmt.Sprintf("Please specify docker password for docker user"))
		}
	}
	if puser != "" {
		if ppaswd == "" {
			printMsgAndExit(fmt.Sprintf("Please specify podman password for podman user"))
		}
	}
	if duser == "" && puser == "" {
		printMsgAndExit(fmt.Sprint("Please provide at least one registry creds "))
	}
}

func printMsgAndExit(msg string) {
	log.Println(msg)
	usage()
	os.Exit(1)
}

func logAndExit(err error) {
	if err != nil {
		log.Fatalf("Fatal Error: %v", err)
	}
}
