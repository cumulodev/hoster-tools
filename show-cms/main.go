package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"

	"github.com/cumulodev/nimbusec"
)

func main() {
	domain := flag.String("domain", "ALL", "define specific domain or ALL to lookup over all domains and resources")
	url := flag.String("url", nimbusec.DefaultAPI, "url to nimbusec API")
	key := flag.String("key", "", "nimbusec API key")
	secret := flag.String("secret", "", "nimbusec API secret")
	humanReadable := flag.Bool("humanreadable", true, "default true: return human readable cms name; set false for internal cpeID")
	outdated := flag.Bool("outdated", true, "default true: show only domains with outdated cms versions; false: show all")
	flag.Parse()

	api, err := nimbusec.NewAPI(*url, *key, *secret)
	if err != nil {
		log.Fatal(err)
	}

	// find domains
	var domains []nimbusec.Domain
	if *domain != "ALL" {
		obj, err := api.GetDomainByName(*domain)
		if err != nil {
			log.Fatal(err)
		}

		domains = []nimbusec.Domain{*obj}
	} else {
		domains, err = api.FindDomains(nimbusec.EmptyFilter)
		if err != nil {
			log.Fatal(err)
		}
	}

	writer := csv.NewWriter(os.Stdout)
	for _, domain := range domains {
		cms, err := api.GetDomainCMS(domain.Id)
		if err != nil {
			log.Fatal(err)
		}

		for _, elem := range cms {
			if !*outdated || elem.LatestStable != "" {
				currcms := elem.CMS
				latest := elem.LatestStable
				if *humanReadable {
					currcms, err = api.GetCMSName(elem.CMS)
					if err != nil {
						log.Fatal(err)
					}
					if elem.LatestStable != "" {
						latest, err = api.GetCMSName(elem.LatestStable)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
				writer.Write([]string{domain.Name, currcms, latest})
				writer.Flush()
			}
		}
	}

}
