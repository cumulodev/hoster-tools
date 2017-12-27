package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/cumulodev/nimbusec"
)

func main() {
	filter := flag.String("filter", "severity ge 3 and (event eq \"malware\" or event eq \"webshell\")", "filter for when a domain is considered infected")
	domain := flag.String("domain", "ALL", "define specific domain or ALL to lookup over all domains and resources")
	url := flag.String("url", nimbusec.DefaultAPI, "url to nimbusec API")
	key := flag.String("key", "", "nimbusec API key")
	secret := flag.String("secret", "", "nimbusec API secret")
	flag.Parse()

	api, err := nimbusec.NewAPI(*url, *key, *secret)
	if err != nil {
		log.Fatal(err)
	}

	// find infected domains
	var domains []nimbusec.Domain
	if *domain != "ALL" {
		obj, err := api.GetDomainByName(*domain)
		if err != nil {
			log.Fatal(err)
		}

		domains = []nimbusec.Domain{*obj}
	} else {
		domains, err = api.FindInfected(*filter)
		if err != nil {
			log.Fatal(err)
		}
	}

	// fetch resources per domain
	writer := csv.NewWriter(os.Stdout)
	for _, domain := range domains {
		results, err := api.FindResults(domain.Id, *filter)
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			writer.Write([]string{domain.Name, strconv.Itoa(result.LastDate), result.Resource, result.Threatname, result.Reason})
			writer.Flush()
		}
	}

}
