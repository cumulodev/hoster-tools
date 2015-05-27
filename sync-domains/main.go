package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cumulodev/nimbusec"
)

func main() {

	apiUrlPtr := flag.String("url", nimbusec.DefaultAPI, "API Url")
	apiKeyPtr := flag.String("key", "abc", "API key for authentication")
	apiSecretPtr := flag.String("secret", "abc", "API secret for authentication")
	filePtr := flag.String("file", "import.csv", "path to import file")
	flag.Parse()

	api, err := nimbusec.NewAPI(*apiUrlPtr, *apiKeyPtr, *apiSecretPtr)
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * READ CSV FILE
	 */
	importfile, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}

	defer importfile.Close()
	reader := csv.NewReader(importfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below
	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// BUILD MAP WITH NEW DOMAINS
	ref := make(map[string]struct{})

	for _, each := range rawCSVdata {
		url := each[0]
		scheme := each[1]
		bundle := each[2]

		// BUILD REF
		ref[url] = struct{}{}

		// ADD DOMAIN TO SET
		domain := &nimbusec.Domain{
			Name:      url,
			Bundle:    bundle,
			Scheme:    scheme,
			DeepScan:  scheme + "://" + url,
			FastScans: []string{scheme + "://" + url},
		}

		// UPSERT DOMAIN
		fmt.Printf("UPSERT DOMAIN: %+v\n", domain)
		_, err := api.CreateOrUpdateDomain(domain)
		if err != nil {
			log.Fatal(err)
		}

	}

	// READ ALL DOMAINS FROM API
	currDomains, err := api.FindDomains(nimbusec.EmptyFilter)
	if err != nil {
		log.Fatal(err)
	}

	// SYNC
	// DELETE DOMAINS NOT LISTED IN NEW SET
	for _, d := range currDomains {
		if _, ok := ref[d.Name]; !ok {
			fmt.Println("I would now delete Domain " + d.Name)
			//api.DeleteDomain(d,true)

		}
	}
}
