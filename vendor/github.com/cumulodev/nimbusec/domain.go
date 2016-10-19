package nimbusec

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Domain represents a nimbusec monitored domain.
type Domain struct {
	Id        int      `json:"id,omitempty"` // Unique identification of domain
	Bundle    string   `json:"bundle"`       // ID of assigned bundle
	Name      string   `json:"name"`         // Name of domain (usually DNS name)
	Scheme    string   `json:"scheme"`       // Flag whether the domain uses http or https
	DeepScan  string   `json:"deepScan"`     // Starting point for the domain deep scan
	FastScans []string `json:"fastScans"`    // Landing pages of the domain scanned
}

type DomainEvent struct {
	Time    Timestamp `json:"time"`
	Event   string    `json:"event"`
	Human   string    `json:"human"`
	Machine string    `json:"machine"`
}

type Timestamp struct {
	time.Time
}

type DomainMetadata struct {
	LastDeepScan Timestamp `json:"lastDeepScan"` // timestamp (in ms) of last external scan of the whole site
	NextDeepScan Timestamp `json:"nextDeepScan"` // timestamp (in ms) for next external scan of the whole site
	LastFastScan Timestamp `json:"lastFastScan"` // timestamp (in ms) of last external scan of the landing pages
	NextFastScan Timestamp `json:"nextFastScan"` // timestamp (in ms) for next external scan of the landing pages
	Agent        Timestamp `json:"agent"`        // status of server agent for the given domain
	Files        int       `json:"files"`        // number of downloaded files/URLs for last deep scan
	Size         int       `json:"size"`         // size of downloaded files for last deep scan (in byte)}
}

type DomainApplication struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Path       string `json:"path"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	Latest     bool   `json:"latest"`
	Vulnerable bool   `json:"vulnerable"`
}

type DomainIssues struct {
	DomainID int    `json:"domainId,omitempty"`
	Category string `json:"category"`
	Issues   int    `json:"issues"`
	Severity int    `json:"severity"`
	Src      string `json:"src"`
}

type Screenshot struct {
	Target   string `json:"target"`
	Previous struct {
		Date     Timestamp `json:"date"`
		MimeType string    `json:"mime"`
		URL      string    `json:"url"`
	} `json:"previous"`
	Current struct {
		Date     Timestamp `json:"date"`
		MimeType string    `json:"mime"`
		URL      string    `json:"url"`
	} `json:"current"`
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Unix()
	stamp := strconv.FormatInt(ts*1000, 10)
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	ts, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	*t = Timestamp{time.Unix(ts/1000, 0)}
	return nil
}

// CreateDomain issues the API to create the given domain.
func (a *API) CreateDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.BuildURL("/v2/domain")
	err := a.Post(url, Params{}, domain, dst)
	return dst, err
}

// CreateOrUpdateDomain issues the nimbusec API to create the given domain. Instead
// of failing when attempting to create a duplicate domain, this method will update
// the remote domain instead.
func (a *API) CreateOrUpdateDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.BuildURL("/v2/domain")
	err := a.Post(url, Params{"upsert": "true"}, domain, dst)
	return dst, err
}

// CreateOrGetDomain issues the nimbusec API to create the given domain. Instead
// of failing when attempting to create a duplicate domain, this method will fetch
// the remote domain instead.
func (a *API) CreateOrGetDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.BuildURL("/v2/domain")
	err := a.Post(url, Params{"upsert": "false"}, domain, dst)
	return dst, err
}

// GetDomain retrieves a domain from the API by its ID.
func (a *API) GetDomain(domain int) (*Domain, error) {
	dst := new(Domain)
	url := a.BuildURL("/v2/domain/%d", domain)
	err := a.Get(url, Params{}, dst)
	return dst, err
}

// GetDomainByName fetches an domain by its name.
func (a *API) GetDomainByName(name string) (*Domain, error) {
	domains, err := a.FindDomains(fmt.Sprintf("name eq \"%s\"", name))
	if err != nil {
		return nil, err
	}

	if len(domains) == 0 {
		return nil, ErrNotFound
	}

	if len(domains) > 1 {
		return nil, fmt.Errorf("name %q matched too many domains. please contact nimbusec.", name)
	}

	return &domains[0], nil
}

// FindDomains searches for domains that match the given filter criteria.
func (a *API) FindDomains(filter string) ([]Domain, error) {
	params := Params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Domain, 0)
	url := a.BuildURL("/v2/domain")
	err := a.Get(url, params, &dst)
	return dst, err
}

// UpdateDOmain issues the nimbusec API to update a domain.
func (a *API) UpdateDomain(domain *Domain) (*Domain, error) {
	dst := new(Domain)
	url := a.BuildURL("/v2/domain/%d", domain.Id)
	err := a.Put(url, Params{}, domain, dst)
	return dst, err
}

// DeleteDomain issues the API to delete a domain. When clean=false, the domain and
// all assiciated data will only be marked as deleted, whereas with clean=true the data
// will also be removed from the nimbusec system.
func (a *API) DeleteDomain(d *Domain, clean bool) error {
	url := a.BuildURL("/v2/domain/%d", d.Id)
	return a.Delete(url, Params{
		"pleaseremovealldata": fmt.Sprintf("%t", clean),
	})
}

// FindInfected searches for domains that have pending Results that match the
// given filter criteria.
func (a *API) FindInfected(filter string) ([]Domain, error) {
	params := make(map[string]string)
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Domain, 0)
	url := a.BuildURL("/v2/infected")
	err := a.Get(url, params, &dst)
	return dst, err
}

// ListDomainConfigs fetches the list of all available configuration keys for the
// given domain.
func (a *API) ListDomainConfigs(domain int) ([]string, error) {
	dst := make([]string, 0)
	url := a.BuildURL("/v2/domain/%d/config", domain)
	err := a.Get(url, Params{}, &dst)
	return dst, err
}

// GetDomainConfig fetches the requested domain configuration.
func (a *API) GetDomainConfig(domain int, key string) (string, error) {
	url := a.BuildURL("/v2/domain/%d/config/%s/", domain, key)
	return a.getTextPlain(url, Params{})
}

// SetDomainConfig sets the domain configuration `key` to the requested value.
// This method will create the domain configuration if it does not exist yet.
func (a *API) SetDomainConfig(domain int, key string, value string) (string, error) {
	url := a.BuildURL("/v2/domain/%d/config/%s/", domain, key)
	return a.putTextPlain(url, Params{}, value)
}

// DeleteDomainConfig issues the API to delete the domain configuration with
// the provided key.
func (a *API) DeleteDomainConfig(domain int, key string) error {
	url := a.BuildURL("/v2/domain/%d/config/%s/", domain, key)
	return a.Delete(url, Params{})
}

func (a *API) GetDomainEvent(domain int, filter string, limit int) ([]DomainEvent, error) {
	params := Params{
		"limit": strconv.Itoa(limit),
	}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]DomainEvent, 0)
	url := a.BuildURL("/v2/domain/%d/events", domain)
	err := a.Get(url, params, &dst)
	return dst, err
}

func (a *API) CreateDomainEvent(domain int, log *DomainEvent) error {
	url := a.BuildURL("/v2/domain/%d/events", domain)
	return a.Post(url, Params{}, log, nil)
}

func (a *API) GetDomainMetadata(domain int) (*DomainMetadata, error) {
	dst := new(DomainMetadata)
	url := a.BuildURL("/v2/domain/%d/metadata", domain)
	err := a.Get(url, Params{}, &dst)
	return dst, err
}

func (a *API) GetDomainApplications(domain int) ([]DomainApplication, error) {
	dst := make([]DomainApplication, 0)
	url := a.BuildURL("/v2/domain/%d/applications", domain)
	err := a.Get(url, Params{}, &dst)
	return dst, err
}

func (a *API) GetDomainScreenshot(domain int) (*Screenshot, error) {
	return a.GetSpecificDomainScreenshot(domain, "EU", "desktop")
}

func (a *API) GetSpecificDomainScreenshot(domain int, region, viewport string) (*Screenshot, error) {
	dst := new(Screenshot)
	url := a.BuildURL("/v2/domain/%d/screenshot/%s/%s", domain, region, viewport)
	err := a.Get(url, Params{}, &dst)
	return dst, err
}

func (a *API) GetImage(url string) ([]byte, error) {
	resolved := a.BuildURL(url)
	return a.getBytes(resolved, Params{})
}

func (a *API) GetIssues() ([]DomainIssues, error) {
	dst := make([]DomainIssues, 0)
	url := a.BuildURL("/v2/domainissues")
	err := a.Get(url, Params{}, &dst)
	return dst, err
}
