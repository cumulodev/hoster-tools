package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cumulodev/goutils/pool"
	"github.com/cumulodev/nimbusec"
)

func main() {
	url := flag.String("url", nimbusec.DefaultAPI, "API Url")
	key := flag.String("key", "abc", "API key for authentication")
	secret := flag.String("secret", "abc", "API secret for authentication")
	file := flag.String("file", "import.csv", "path to import file")
	dryrun := flag.Bool("dry-run", false, "simulate what would be done without writing")

	flag.Parse()

	// creates a new nimbusec API instance
	api, err := nimbusec.NewAPI(*url, *key, *secret)
	if err != nil {
		log.Fatal(err)
	}

	// open csv input file and parse it
	fh, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()
	reader := csv.NewReader(fh)
	reader.FieldsPerRecord = -1 // see the Reader struct information below
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	pool := pool.New(1)
	pool.Start()

	// keep track of domain names (required for delete step later)
	ref := make(map[string]int)

	domains, err := api.FindDomains(nimbusec.EmptyFilter)
	if err != nil {
		log.Fatal(err)
	}

	for _, domain := range domains {
		ref[domain.Name] = domain.Id
	}

	for _, row := range rows {
		name := row[0]
		scheme := row[2]
		bundle := row[3]

		url := scheme + "://" + name
		//		if len(row) > 4 {
		//			deeplink := row[4]
		//			if deeplink == "" {
		//				// do nothing
		//			} else if strings.HasPrefix(deeplink, "/") {
		//				url = url + deeplink
		//			} else {
		//				url = deeplink
		//			}
		//		}

		// construct domain
		domain := &nimbusec.Domain{
			Name:      name,
			Bundle:    bundle,
			Scheme:    scheme,
			DeepScan:  url,
			FastScans: []string{url},
		}

		if _, ok := ref[domain.Name]; ok {
			if *dryrun {
				fmt.Printf("i would now delete '%s'\n", domain.Name)
			} else {
				pool.Add(deleteJob{
					api: api,
					domain: nimbusec.Domain{
						Id:   ref[domain.Name],
						Name: domain.Name,
					},
				})
			}

		}

	}

	pool.Wait()

}

type deleteJob struct {
	api    *nimbusec.API
	domain nimbusec.Domain
}

func (job deleteJob) Work() {
	fmt.Printf("delete domain: %s\n", job.domain.Name)
	job.api.DeleteDomain(&job.domain, true)
}

func (job deleteJob) Save() {}
