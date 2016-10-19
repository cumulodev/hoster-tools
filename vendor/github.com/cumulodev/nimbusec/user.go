package nimbusec

import "fmt"

const (
	// RoleUser is the restricted role for an user
	RoleUser = "user"

	// RoleAdministrator is the unrestricted role for an user
	RoleAdministrator = "administrator"
)

// User represents an human user able to login and receive notifications.
type User struct {
	Id           int    `json:"id,omitempty"`           // unique identification of user
	Login        string `json:"login"`                  // login name of user
	Mail         string `json:"mail"`                   // e-mail contact where mail notifications are sent to
	Role         string `json:"role"`                   // role of an user (`administrator` or `user`
	Company      string `json:"company"`                // company name of user
	Surname      string `json:"surname"`                // surname of user
	Forename     string `json:"forename"`               // forename of user
	Title        string `json:"title"`                  // academic title of user
	Mobile       string `json:"mobile"`                 // phone contact where sms notifications are sent to
	Password     string `json:"password,omitempty"`     // password of user (only used when creating or updating a user)
	SignatureKey string `json:"signatureKey,omitempty"` // secret for SSO (only used when creating or updating a user)
}

// Notification represents an notification entry for a user and domain.
type Notification struct {
	Id         int    `json:"id,omitempty"` // unique identification of notification
	Domain     int    `json:"domain"`       // domain for which notifications should be sent
	Transport  string `json:"transport"`    // transport over which notifications are sent (mail, sms)
	ServerSide int    `json:"serverside"`   // minimum severity of serverside results before a notification is sent
	Content    int    `json:"content"`      // minimum severity of content results before a notification is sent
	Blacklist  int    `json:"blacklist"`    // minimum severity of backlist results before a notification is sent
}

// CreateUser issues the nimbusec API to create the given user.
func (a *API) CreateUser(user *User) (*User, error) {
	dst := new(User)
	url := a.BuildURL("/v2/user")
	err := a.Post(url, Params{}, user, dst)
	return dst, err
}

// CreateOrUpdateUser issues the nimbusec API to create the given user. Instead of
// failing when attempting to create a duplicate user, this method will update the
// remote user instead.
func (a *API) CreateOrUpdateUser(user *User) (*User, error) {
	dst := new(User)
	url := a.BuildURL("/v2/user")
	err := a.Post(url, Params{"upsert": "true"}, user, dst)
	return dst, err
}

// CreateOrGetUser issues the nimbusec API to create the given user. Instead of
// failing when attempting to create a duplicate user, this method will fetch the
// remote user instead.
func (a *API) CreateOrGetUser(user *User) (*User, error) {
	dst := new(User)
	url := a.BuildURL("/v2/user")
	err := a.Post(url, Params{"upsert": "false"}, user, dst)
	return dst, err
}

// GetUser fetches an user by its ID.
func (a *API) GetUser(user int) (*User, error) {
	dst := new(User)
	url := a.BuildURL("/v2/user/%d", user)
	err := a.Get(url, Params{}, dst)
	return dst, err
}

// GetUserByLogin fetches an user by its login name.
func (a *API) GetUserByLogin(login string) (*User, error) {
	users, err := a.FindUsers(fmt.Sprintf("login eq \"%s\"", login))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("login %q matched too many users. please contact nimbusec.", login)
	}

	return &users[0], nil
}

// FindUsers searches for users that match the given filter criteria.
func (a *API) FindUsers(filter string) ([]User, error) {
	params := Params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]User, 0)
	url := a.BuildURL("/v2/user")
	err := a.Get(url, params, &dst)
	return dst, err
}

// UpdateUser issues the nimbusec API to update an user.
func (a *API) UpdateUser(user *User) (*User, error) {
	dst := new(User)
	url := a.BuildURL("/v2/user/%d", user.Id)
	err := a.Put(url, Params{}, user, dst)
	return dst, err
}

// DeleteUser issues the nimbusec API to delete an user. The root user or tennant
// can not be deleted via this method.
func (a *API) DeleteUser(user *User) error {
	url := a.BuildURL("/v2/user/%d", user.Id)
	return a.Delete(url, Params{})
}

// GetDomainSet fetches the set of allowed domains for an restricted user.
func (a *API) GetDomainSet(user *User) ([]int, error) {
	dst := make([]int, 0)
	url := a.BuildURL("/v2/user/%d/domains", user.Id)
	err := a.Get(url, Params{}, &dst)
	return dst, err
}

// UpdateDomainSet updates the set of allowed domains of an restricted user.
func (a *API) UpdateDomainSet(user *User, domains []int) ([]int, error) {
	dst := make([]int, 0)
	url := a.BuildURL("/v2/user/%d/domains", user.Id)
	err := a.Put(url, Params{}, domains, &dst)
	return dst, err
}

// LinkDomain links the given domain id to the given user and adds the priviledges for
// the user to view the domain.
func (a *API) LinkDomain(user *User, domain int) error {
	url := a.BuildURL("/v2/user/%d/domains", user.Id)
	return a.Post(url, Params{}, domain, nil)
}

// UnlinkDomain unlinks the given domain id to the given user and removes the priviledges
// from the user to view the domain.
func (a *API) UnlinkDomain(user *User, domain int) error {
	url := a.BuildURL("/v2/user/%d/domains/%d", user.Id, domain)
	return a.Delete(url, Params{})
}

// ListuserConfigs fetches the list of all available configuration keys for the
// given domain.
func (a *API) ListUserConfigs(user int) ([]string, error) {
	dst := make([]string, 0)
	url := a.BuildURL("/v2/user/%d/config", user)
	err := a.Get(url, Params{}, &dst)
	return dst, err
}

// GetUserConfig fetches the requested user configuration.
func (a *API) GetUserConfig(user int, key string) (string, error) {
	url := a.BuildURL("/v2/user/%d/config/%s/", user, key)
	return a.getTextPlain(url, Params{})
}

// SetUserConfig sets the user configuration `key` to the requested value.
// This method will create the user configuration if it does not exist yet.
func (a *API) SetUserConfig(user int, key string, value string) (string, error) {
	url := a.BuildURL("/v2/user/%d/config/%s/", user, key)
	return a.putTextPlain(url, Params{}, value)
}

// DeleteUserConfig issues the API to delete the user configuration with
// the provided key.
func (a *API) DeleteUserConfig(user int, key string) error {
	url := a.BuildURL("/v2/user/%d/config/%s/", user, key)
	return a.Delete(url, Params{})
}

// GetNotification fetches a notification by its ID.
func (a *API) GetNotification(user int, id int) (*Notification, error) {
	dst := new(Notification)
	url := a.BuildURL("/v2/user/%d/notification/%d", user, id)
	err := a.Get(url, Params{}, dst)
	return dst, err
}

// FindNotifications fetches all notifications for the given user that match the
// filter criteria.
func (a *API) FindNotifications(user int, filter string) ([]Notification, error) {
	params := Params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Notification, 0)
	url := a.BuildURL("/v2/user/%d/notification", user)
	err := a.Get(url, params, &dst)
	return dst, err
}

// CreateNotification creates the notification for the given user.
func (a *API) CreateNotification(user int, notification *Notification) (*Notification, error) {
	dst := new(Notification)
	url := a.BuildURL("/v2/user/%d/notification", user)
	err := a.Post(url, Params{}, notification, dst)
	return dst, err
}

// CreateOrUpdateNotification issues the nimbusec API to create the given notification. Instead of
// failing when attempting to create a duplicate notification, this method will update the
// remote notification instead.
func (a *API) CreateOrUpdateNotification(user int, notification *Notification) (*Notification, error) {
	dst := new(Notification)
	url := a.BuildURL("/v2/user/%d/notification", user)
	err := a.Post(url, Params{"upsert": "true"}, notification, dst)
	return dst, err
}

// CreateOrGetNotifcation issues the nimbusec API to create the given notification. Instead of
// failing when attempting to create a duplicate notification, this method will fetch the
// remote notification instead.
func (a *API) CreateOrGetNotification(user int, notification *Notification) (*Notification, error) {
	dst := new(Notification)
	url := a.BuildURL("/v2/user/%d/notification", user)
	err := a.Post(url, Params{"upsert": "false"}, notification, dst)
	return dst, err
}

// UpdateNotification updates the notification for the given user.
func (a *API) UpdateNotification(user int, notification *Notification) (*Notification, error) {
	dst := new(Notification)
	url := a.BuildURL("/v2/user/%d/notification/%d", user, notification.Id)
	err := a.Put(url, Params{}, notification, dst)
	return dst, err
}

// DeleteNotification deletes the given notification.
func (a *API) DeleteNotification(user int, notification *Notification) error {
	url := a.BuildURL("/v2/user/%d/notification/%d", user, notification.Id)
	return a.Delete(url, Params{})
}
