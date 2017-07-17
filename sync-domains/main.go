package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cumulodev/goutils/pool"
	"github.com/cumulodev/nimbusec"
)

func main() {
	url := flag.String("url", nimbusec.DefaultAPI, "API Url")
	key := flag.String("key", "abc", "API key for authentication")
	secret := flag.String("secret", "abc", "API secret for authentication")
	file := flag.String("file", "import.csv", "path to import file")
	delete := flag.Bool("delete", false, "delete domains from nimbusec if not provided in the CSV")
	update := flag.Bool("update", false, "updates domain info; false to just insert new domains")
	workers := flag.Int("workers", 1, "number of parallel workers (please do not use too many workers)")
	throttle := flag.Int("throttle", 0, "delay in seconds per domain creation")
	flag.Parse()

	waitdur := time.Duration(*throttle) * time.Second

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

	pool := pool.New(*workers)
	pool.Start()

	// keep track of domain names (required for delete step later)
	ref := make(map[string]bool)

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

		// upsert domain
		ref[name] = true
		pool.Add(upsertJob{
			api:      api,
			domain:   *domain,
			update:   *update,
			throttle: waitdur,
		})
	}

	pool.Wait()

	// sync
	// delete domains not listed in new set
	if *delete {
		// read all domains from api
		domains, err := api.FindDomains(nimbusec.EmptyFilter)
		if err != nil {
			log.Fatal(err)
		}

		// cross reference domains in nimbusec with csv file and delete all
		// domains not present in csv file
		for _, domain := range domains {
			if !ref[domain.Name] {
				pool.Add(deleteJob{
					api:    api,
					domain: domain,
				})
			}
		}

		pool.Wait()
	}
}

type upsertJob struct {
	api      *nimbusec.API
	domain   nimbusec.Domain
	update   bool
	throttle time.Duration
}

func (job upsertJob) Work() {
	fmt.Printf("upsert domain: %+v\n", job.domain)
	if job.update {
		if _, err := job.api.CreateOrUpdateDomain(&job.domain); err != nil {
			log.Fatal(err)
		}
	} else {
		if _, err := job.api.CreateOrGetDomain(&job.domain); err != nil {
			log.Fatal(err)
		}
	}
	if job.throttle > 0 {
		time.Sleep(job.throttle)
	}
}

func (job upsertJob) Save() {}

type deleteJob struct {
	api    *nimbusec.API
	domain nimbusec.Domain
}

func (job deleteJob) Work() {
	fmt.Printf("delete domain: %s\n", job.domain.Name)
	job.api.DeleteDomain(&job.domain, true)
}

func (job deleteJob) Save() {}
