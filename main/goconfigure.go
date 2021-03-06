// GoConfigure is an application for templating configurations meant to be pushed via SSH.
// It can be consumed as a command line tool. Passing arguments of the form `goconfigure_0.1_arm64
// -t TEMPLATE_FILE_NAME -i DATA_FILE_NAME` will render and push the render to the devices
// defined in the inventory file.
package main

import (
	"flag"
	"github.com/dyntek-services-inc/goconfigure/client"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"github.com/dyntek-services-inc/goconfigure/render"
	"log"
	"os"
	"strings"
)

var pwd, pwdErr = os.Getwd()

func main() {
	if pwdErr != nil {
		log.Fatal(pwdErr)
	}
	invFilename := flag.String("i", "", "inventory filename")
	tplFilename := flag.String("t", "", "template filename")
	keyFilename := flag.String("k", "", "PEM key filename")
	flag.Parse()
	if len(*invFilename) == 0 && len(*tplFilename) == 0 {
		// No inventory or template flags were passed, start manual mode
		// TODO: implement manual mode
	} else {
		if len(*invFilename) == 0 || len(*tplFilename) == 0 {
			// One of the flags was passed but not the other
			log.Fatal("one flag was passed, but not both")
		} else {
			// Both flags were passed, begin loading
			var inv *inventory.Inventory
			var err error
			if strings.HasSuffix(*invFilename, ".csv") {
				// The passed inventory file was a CSV
				// If a keyfile was provided passwords are not required in the CSV
				inv, err = inventory.LoadFromCSV(*invFilename, len(*keyFilename) == 0)
			} else if strings.HasSuffix(*invFilename, ".yml") || strings.HasSuffix(*invFilename, ".yaml") {
				inv, err = inventory.LoadFromYAML(*invFilename)
			} else {
				log.Fatal("the passed inventory file is not of type CSV or YAML")
			}
			if err != nil {
				log.Fatal(err)
			}
			// Begin loading the template
			tplString, err := render.FileToString(*tplFilename)
			if err != nil {
				log.Fatal(err)
			}
			// Determine authentication method
			var auth client.Authentication
			if len(*keyFilename) == 0 {
				// No auth key was passed, use Basic
				auth, err = client.BasicConnect()
			} else {
				// an auth key was provided, use Key
				auth, err = client.KeyConnect(*keyFilename)
			}
			if err != nil {
				log.Fatal(err)
			}
			// Begin the deployment
			deployment := client.NewDeployment(tplString, pwd)
			if err := deployment.Deploy(inv, auth); err != nil {
				log.Fatal(err)
			}
		}
	}
}
