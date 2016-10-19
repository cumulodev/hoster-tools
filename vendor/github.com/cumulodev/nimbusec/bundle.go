package nimbusec

type Bundle struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name"`
	Start      Timestamp `json:"startDate"`
	End        Timestamp `json:"endDate"`
	Quota      string    `json:"quota"`
	Depth      int       `json:"depth"`
	Fast       int       `json:"fast"`
	Deep       int       `json:"deep"`
	Contingent int       `json:"contingent"`
	Active     int       `json:"active"`
	Engines    []string  `json:"engines"`
	Amount     int       `json:"amount"`
	Currency   string    `json:"currency"`
}

func (a *API) GetBundle(bundle string) (*Bundle, error) {
	dst := new(Bundle)
	url := a.BuildURL("/v2/bundle/%s", bundle)
	err := a.Get(url, Params{}, dst)
	return dst, err
}

func (a *API) FindBundles(filter string) ([]Bundle, error) {
	params := Params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Bundle, 0)
	url := a.BuildURL("/v2/bundle")
	err := a.Get(url, params, &dst)
	return dst, err
}
