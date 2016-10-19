package nimbusec

// Result represents a finding of the nimbusec service that requires user action.
type Result struct {
	Id           int     `json:"id,omitempty"` // unique identification of a result
	Status       int     `json:"status"`       // status of the result (1 = pending, 2 = acknowledged, 3 = falsepositive, 4 = removed)
	Event        string  `json:"event"`        // event type of result (e.g added file)
	Category     string  `json:"category"`     // category of result
	Severity     int     `json:"severity"`     // severity level of result (1 = medium to 3 = severe)
	Probability  float64 `json:"probability"`  // probability the result is critical
	SafeToDelete bool    `json:"safeToDelete"` // flag indicating if the file can be safely deleted without loosing user data
	CreateDate   int     `json:"createDate"`   // timestamp (in ms) of the first occurrence
	LastDate     int     `json:"lastDate"`     // timestamp (in ms) of the last occurrence

	// the following fields contain more details about the result. Not all fields
	// must be filled or present.

	Threatname string `json:"threatname"` // name identifying the threat of a result
	Resource   string `json:"resource"`   // affected resource (e.g. file path or URL)
	MD5        string `json:"md5"`        // MD5 hash sum of the affected file
	Filesize   int    `json:"filesize"`   // filesize of the affected file
	Owner      string `json:"owner"`      // file owner of the affected file
	Group      string `json:"group"`      // file group of the affected file
	Permission int    `json:"permission"` // permission of the affected file as decimal integer
	Diff       string `json:"diff"`       // diff of a content change between two scans
	Reason     string `json:"reason"`     // reason why a domain/URL is blacklisted
}

// GetResult fetches a result by its ID.
func (a *API) GetResult(domain, result int) (*Result, error) {
	dst := new(Result)
	url := a.BuildURL("/v2/domain/%d/result/%d", domain, result)
	err := a.Get(url, Params{}, dst)
	return dst, err
}

// FindResults searches for results that match the given filter criteria.
func (a *API) FindResults(domain int, filter string) ([]Result, error) {
	params := make(map[string]string)
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Result, 0)
	url := a.BuildURL("/v2/domain/%d/result", domain)
	err := a.Get(url, params, &dst)
	return dst, err
}

// UpdateResult issues the nimbusec API to update a result, all fields except
// status will be ignored.
func (a *API) UpdateResult(domain int, result *Result) (*Result, error) {
	dst := new(Result)
	url := a.BuildURL("/v2/domain/%d/result/%d", domain, result.Id)
	err := a.Put(url, Params{}, result, dst)
	return dst, err
}
