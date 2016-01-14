# hoster-tools
Collection of simple tools for hoster


## create-agent-config
Create an nimbusec agent configuration file based on an import CSV file. The `create-agent-config` tool outputs the configuration to stdout, so redirect it to the desired path (e.g. /opt/nimbusec/agent.conf).

### Installation
If you have Go installed, the `create-agent-config` can simply be installed by go get:

    go get github.com/cumulodev/hoster-tools/create-agent-config

### Example Usage
As `key` and `secret` please use an *Server Agent Token* (can be found at https://portal.nimbusec.com/einstellungen/serveragent).

    create-agent-config -key abc -secret abc -file import.csv > /opt/nimbusec/agent.conf
    
Or define another path for writing the temporary results file
    
    create-agent-config -key abc -secret abc -tmpfile C:\\tmp\\nimbusec.tmp -file import.csv > /opt/nimbusec/agent.conf
  
An example for the import.csv file is in the create-agent-config directory.

## sync-domains
Syncs the provided import CSV file with the nimbusec system. `sync-domains` provides a 2-way sync, which means that all domains in the CSV will be created in our system, while any domains missing in the CSV (but present in our system) will be removed. The 2-way delete can be disabled via `delete` option.

### Installation
If you have Go installed, the `sync-domains` can simply be installed by go get:

    go get github.com/cumulodev/hoster-tools/sync-domains

### Usage
* *delete*: default FALSE; delete domains from nimbusec if not provided in the CSV
* *update*: default FALSE; updates domain info (e.g. bundle ID); ELSE just inserts new domains without modifying existing
* *workers*: default 1; number of parallel workers which will increase sync time (please do not use too many workers)

As `key` and `secret` please use your assigned API key and secret (can be found at https://portal.nimbusec.com/einstellungen/serveragent).

    sync-domains -key abc -secret abc -file import.csv
    
Or to disable the removal of domains from nimbusec:

    sync-domains -delete false -key abc -secret abc -file import.csv
  
To keep existing domains in their defined bundles, but add new ones with another bundle entry

    sync-domains -update false -key abc -secret abc -file import.csv

An example for the import.csv file is in the sync-domains directory.

## infected-domain-trigger
This tool polls the nimbusec API in an specified interval for infected domains and performs certain actions on it. An example use case would be the automatic disabling of infected domains.

### Installation
If you have Go installed, the `infected-domain-trigger` can simply be installed by go get:

    go get github.com/cumulodev/hoster-tools/infected-domain-trigger
    
### Usage
As `key` and `secret` please use your assigned API key and secret (can be found at https://portal.nimbusec.com/einstellungen/serveragent).

    infected-domain-trigger -key abc -secret abc -action 'echo "infected $DOMAIN"' -reload 'echo "reloading httpd"'
    
* *action*: The action command will be executed for each infected domain. The command will be executed in an shell, where the environment variable `DOMAIN` is set to the name of the infected domain.
* *reload*: The reload command will be executed after each interval if nimbusec reported infected domains. This can be used to issue e.g. Apache to reload the configuration.
    
To disable for example all infected domains hosted by Apache, specify the following actions:

    infected-domain-trigger -key abc -secret abc -action 'a2dissite $DOMAIN' -reload 'apachectl graceful'
    
If one of the actions is not required, just specify for example the shell builtin `true` command:

    infected-domain-trigger -key abc -secret abc -action 'disable.sh' -reload 'true'

## infected-resources
This tool polls the nimbusec API for infected domains and returns the infected resources (files and paths). An example use case would be the automatic delete/move/quarantine of infected files.

### Installation
If you have Go installed, the `infected-resources` can simply be installed by go get:

    go get github.com/cumulodev/hoster-tools/infected-resources
    
### Usage
As `key` and `secret` please use your assigned API key and secret (can be found at https://portal.nimbusec.com/einstellungen/serveragent).

    infected-resources -key abc -secret abc -domain www.example.com
    
* *domain*: The domain command is OPTIONAL and can limit the output of the resources to one specific domain.

To get the infected resources of ALL domains just call:

    infected-resources -key abc -secret abc
    
If you want the resources of just one domain you may limit it like this:

    infected-resources -key abc -secret abc -domain www.example.com

The output has csv format and can be written to a file like this:

    infected-resources -key abc -secret abc > infected-resources.csv

Perform actions on the output with pipes - infinite possibilities - example:
    
    infected-resources -key abc -secret abc -domain www.example.com | sed -E 's/.*,(.*)$/echo \1/' > process.sh && sh process.sh

## show-cms
With this handy script it is possible to get an overview over all used cms on the infrastructure. 

### Installation
If you have Go installed, the `show-cms` can simply be installed by go get:

    go get github.com/cumulodev/hoster-tools/show-cms
    
### Usage
As `key` and `secret` please use your assigned API key and secret (can be found at https://portal.nimbusec.com/einstellungen/serveragent).

    show-cms -key abc -secret abc -domain www.example.com
    
* *domain*: The domain command is OPTIONAL and can limit the output of the resources to one specific domain. Else it will output ALL provisioned domains
* *outdated*: TRUE (default) shows only domains with outdated CMS, FALSE shows every domain where a CMS version could be determined. However, the latest CMS version field will stay empty when there is no newer version
* *humanreadable*: TRUE (default) prints out the commercial name of a CMS, FALSE uses the CPE-ID
 
To get ALL domains with identified CMS versions:

    show-cms -key abc -secret abc -outdated=false
    
You can also export this to csv:

    show-cms -key abc -secret abc > domains-cms-outdated.csv
