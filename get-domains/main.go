package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"

	"github.com/cumulodev/nimbusec"
)

func main() {
	url := flag.String("url", nimbusec.DefaultAPI, "url to nimbusec API")
	key := flag.String("key", "", "nimbusec API key")
	secret := flag.String("secret", "", "nimbusec API secret")
	flag.Parse()

	api, err := nimbusec.NewAPI(*url, *key, *secret)
	if err != nil {
		log.Fatal(err)
	}

	// find domains
	var domains []nimbusec.Domain
	domains, err = api.FindDomains(nimbusec.EmptyFilter)
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	for _, domain := range domains {
		writer.Write([]string{domain.Name, "", domain.Scheme, domain.Bundle, ""})
		writer.Flush()
	}

}
