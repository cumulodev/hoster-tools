package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/cumulodev/nimbusec"
)

type AgentConfig struct {
	Key           string            `json:"key"`
	Secret        string            `json:"secret"`
	Domains       map[string]string `json:"domains"`
	TmpFile       string            `json:"tmpfile"`
	ExcludeDir    []string          `json:"excludeDir"`
	ExcludeRegexp []string          `json:"excludeRegexp"`
	APIServer     string            `json:"apiserver"`
}

func main() {
	api := flag.String("url", nimbusec.DefaultAPI, "API URL")
	key := flag.String("key", "abc", "Agent Key")
	secret := flag.String("secret", "abc", "Agent Secret")
	filename := flag.String("file", "import.csv", "path to import file")
	tmpfile := flag.String("tmpfile", "/tmp/nimbusec.tmp", "path of the tmpfile that writes interim results")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	docroots := make(map[string]string)
	for _, row := range rows {
		url := row[0]
		docroot := row[1]

		if url != "" {
			docroots[url] = docroot
		}
	}

	conf := AgentConfig{
		Key:           *key,
		Secret:        *secret,
		APIServer:     *api,
		TmpFile:       *tmpfile,
		ExcludeDir:    []string{},
		ExcludeRegexp: []string{},
		Domains:       docroots,
	}

	data, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(data)
	os.Stdout.Write([]byte{'\n'})
}
