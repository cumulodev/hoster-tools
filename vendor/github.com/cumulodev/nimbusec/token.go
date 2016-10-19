package nimbusec

// Token represents the credentials of an API or agent for the nimbusec API.
type Token struct {
	Id       int    `json:"id"`       // unique identification of a token
	Name     string `json:"name"`     // given name for a token
	Key      string `json:"key"`      // oauth key
	Secret   string `json:"secret"`   // oauth secret
	LastCall int    `json:"lastCall"` // last timestamp (in ms) an agent used the token
	Version  int    `json:"version"`  // last agent version that was seen for this key
}

// CreateToken issues the nimbusec API to create a new agent token.
func (a *API) CreateToken(token *Token) (*Token, error) {
	dst := new(Token)
	url := a.BuildURL("/v2/agent/token")
	err := a.Post(url, Params{}, token, dst)
	return dst, err
}

// GetToken fetches a token by its ID.
func (a *API) GetToken(token int) (*Token, error) {
	dst := new(Token)
	url := a.BuildURL("/v2/agent/token/%d", token)
	err := a.Get(url, Params{}, dst)
	return dst, err
}

// FindTOkens searches for tokens that match the given filter criteria.
func (a *API) FindTokens(filter string) ([]Token, error) {
	params := Params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Token, 0)
	url := a.BuildURL("/v2/agent/token")
	err := a.Get(url, params, &dst)
	return dst, err
}
