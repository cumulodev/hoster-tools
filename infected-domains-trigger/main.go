package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cumulodev/nimbusec"
)

func main() {
	filter := flag.String("filter", "severity ge 3 and (event eq \"malware\" or event eq \"webshell\")", "filter for when a domain is considered infected")
	sleep := flag.Int("sleep", 5, "sleep interval in minutes between checks")

	action := flag.String("action", "echo \">> infected: $DOMAIN\"", "execute command for each infected domain")
	reload := flag.String("reload", "echo \"reload trigger\"", "execute command after processing of infected domains (only called if there were infected domains)")

	url := flag.String("url", nimbusec.DefaultAPI, "url to nimbusec API")
	key := flag.String("key", "", "nimbusec API key")
	secret := flag.String("secret", "", "nimbusec API secret")
	flag.Parse()

	//TODO: validate parameters

	api, err := nimbusec.NewAPI(*url, *key, *secret)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// find infected domains
		domains, err := api.FindInfected(*filter)
		if err != nil {
			log.Fatal(err)
		}

		// execute action hook for each infected domain
		for _, domain := range domains {
			run(*action, domain.Name)
		}

		// execute reload hook if the filter matched something
		if len(domains) > 0 {
			run(*reload, "")
		}

		time.Sleep(time.Duration(*sleep) * time.Minute)
	}
}

func run(name string, domain string) {
	cmd := exec.Command("sh", "-c", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"DOMAIN=" + domain}

	err := cmd.Run()
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}
