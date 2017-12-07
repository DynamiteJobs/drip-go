package drip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const baseURL = "https://api.getdrip.com/v2/"

var (
	// ErrBadAPIKey is returned by New if bad api key.
	ErrBadAPIKey = fmt.Errorf("bad api key")
	// ErrBadAccountID is returned by New if bad account ID.
	ErrBadAccountID = fmt.Errorf("bad drip account id")
	// ErrIDorEmailEmpty is returned if id and email are both empty.
	ErrIDorEmailEmpty = fmt.Errorf("ID and Email both cannot be empty")
	// ErrInvalidInput is returned if input is invalid.
	ErrInvalidInput = fmt.Errorf("invalid input")
)

// Client is a client to interact with the Drip API.
// Use https://www.getdrip.com/docs/rest-api for extra documentation.
type Client struct {
	HTTPClient *http.Client
	UserAgent  string
	apiKey     string
	accountID  string
}

// New returns a new Client.
func New(apiKey, accountID string) (*Client, error) {
	if apiKey == "" {
		return nil, ErrBadAPIKey
	}
	if accountID == "" {
		return nil, ErrBadAccountID
	}
	return &Client{
		HTTPClient: http.DefaultClient,
		UserAgent:  "drip-go client",
		apiKey:     apiKey,
		accountID:  accountID,
	}, nil
}

func (c *Client) getReq(method, url string, body interface{}) (*http.Request, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiKey, "")
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", "application/vnd.api+json")
	return req, nil
}

func (c *Client) decodeResp(resp *http.Response, response interface{}) error {
	var err error
	if resp.StatusCode == 204 || strings.Contains(resp.Header.Get("Content-Type"), "No Content") {
		return nil
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "json") {
		var b []byte
		var body string
		b, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			body = string(b)
		}
		return fmt.Errorf("StatusCode(%d) %s", resp.StatusCode, body)
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	return err
}

// Links are links send in responses.
type Links struct {
	Account    string   `json:"account,omitempty"`
	Forms      []string `json:"forms,omitempty"`
	Subscriber string   `json:"subscriber,omitempty"`
}

// Meta is data related to pagination.
// https://www.getdrip.com/docs/rest-api#pagination
type Meta struct {
	Page       int `json:"page,omitempty"`
	Count      int `json:"count,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

// Subscriber is a subscriber.
type Subscriber struct {
	ID               string            `json:"id,omitempty"`
	Status           string            `json:"status,omitempty"`
	Email            string            `json:"email,omitempty"`
	TimeZone         string            `json:"time_zone,omitempty"`
	UTCOffset        int               `json:"utc_offset,omitempty"`
	VisitorUUID      string            `json:"visitor_uuid,omitempty"`
	CustomFields     map[string]string `json:"custom_fields,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	IPAddress        string            `json:"ip_address,omitempty"`
	UserAgent        string            `json:"user_agent,omitempty"`
	OriginalReferrer string            `json:"original_referrer,omitempty"`
	LandingURL       string            `json:"landing_url,omitempty"`
	Prospect         bool              `json:"prospect,omitempty"`
	LeadScore        int               `json:"lead_score,omitempty"`
	LifetimeValue    int               `json:"lifetime_value,omitempty"`
	CreatedAt        time.Time         `json:"created_at,omitempty"`
	HREF             string            `json:"href,omitempty"`
	UserID           string            `json:"user_id,omitempty"`
	BaseLeadScore    int               `json:"base_lead_score,omitempty"`
	Links            Links             `json:"links,omitempty"`
}

// SubscribersResp is a response recieved with subscribers in it.
// List functions have Meta for pagination. StatusCode is included in resp.
type SubscribersResp struct {
	StatusCode  int           `json:"status_code,omniempty"`
	Links       Links         `json:"links,omitempty"`
	Meta        Meta          `json:"meta,omitempty"`
	Subscribers []*Subscriber `json:"subscribers,omitempty"`
	Errors      []CodeError   `json:"errors,omitempty"`
}

// Response is a basic response recieved.
type Response struct {
	StatusCode int         `json:"status_code,omitempty"`
	Errors     []CodeError `json:"errors,omitempty"`
}

// ListSubscribersReq is a request for ListSubscribers.
type ListSubscribersReq struct {
	Status           string     `json:"status,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	SubscribedBefore *time.Time `json:"subscribed_before,omitempty"`
	SubscribedAfter  *time.Time `json:"subscribed_after,omitempty"`
	Page             *int       `json:"page,omitempty"`
	PerPage          *int       `json:"per_page,omitempty"`
}

// ListSubscribers returns a list of subscribers. Either an ID or Email can
func (c *Client) ListSubscribers(req *ListSubscribersReq) (*SubscribersResp, error) {
	url := fmt.Sprintf("%s/%s/subscribers", baseURL, c.accountID)
	httpReq, err := c.getReq(http.MethodGet, url, req)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(SubscribersResp)
	resp.StatusCode = httpResp.StatusCode
	err = c.decodeResp(httpResp, resp)
	return resp, err
}

// UpdateSubscriber is the info available to update or create a subscriber.
type UpdateSubscriber struct {
	Email         string            `json:"email,omitempty"`
	ID            string            `json:"id,omitempty"`
	NewEmail      string            `json:"new_email,omitempty"`
	UserID        string            `json:"user_id,omitempty"`
	TimeZone      string            `json:"time_zone,omitempty"`
	LifetimeValue *float32          `json:"lifetime_value,omitempty"`
	IPAddress     string            `json:"ip_address,omitempty"`
	CustomFields  map[string]string `json:"custom_fields,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	RemoveTags    []string          `json:"remove_tags,omitempty"`
	Prospect      *bool             `json:"prospect,omitempty"`
	BaseLeadScore *int              `json:"base_lead_score,omitempty"`
}

// UpdateSubscribersReq is a request for UpdateSubscriber.
type UpdateSubscribersReq struct {
	Subscribers []UpdateSubscriber `json:"subscribers,omitempty"`
}

// UpdateSubscriber creates or updates a subscriber.
// If you need to create or update a collection of subscribers at once, use our batch API instead.
func (c *Client) UpdateSubscriber(req *UpdateSubscribersReq) (*SubscribersResp, error) {
	url := fmt.Sprintf("%s/%s/subscribers", baseURL, c.accountID)
	httpReq, err := c.getReq(http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(SubscribersResp)
	resp.StatusCode = httpResp.StatusCode
	err = c.decodeResp(httpResp, resp)
	return resp, err
}

// DeleteSubscriber deletes a subscriber.
func (c *Client) DeleteSubscriber(idOrEmail string) (*Response, error) {
	if idOrEmail == "" {
		return nil, ErrIDorEmailEmpty
	}
	url := fmt.Sprintf("%s/%s/subscribers/%s", baseURL, c.accountID, idOrEmail)
	httpReq, err := c.getReq(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	resp.StatusCode = httpResp.StatusCode
	err = c.decodeResp(httpResp, resp)
	return resp, err
}

// FetchSubscriber fetches a subscriber.
func (c *Client) FetchSubscriber(idOrEmail string) (*SubscribersResp, error) {
	if idOrEmail == "" {
		return nil, ErrIDorEmailEmpty
	}
	url := fmt.Sprintf("%s/%s/subscribers/%s", baseURL, c.accountID, idOrEmail)
	httpReq, err := c.getReq(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(SubscribersResp)
	resp.StatusCode = httpResp.StatusCode
	err = c.decodeResp(httpResp, resp)
	return resp, err
}

// TagsReq is a request for TagSubscription.
type TagsReq struct {
	Tags []TagReq `json:"tags,omitempty"`
}

// TagReq is a part of a TagsReq.
type TagReq struct {
	Email string `json:"email,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

// TagSubscriber adds a tag to a subscriber.
func (c *Client) TagSubscriber(req *TagsReq) (*Response, error) {
	url := fmt.Sprintf("%s/%s/tags", baseURL, c.accountID)
	httpReq, err := c.getReq(http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	resp.StatusCode = httpResp.StatusCode
	err = c.decodeResp(httpResp, resp)
	return resp, err
}

// RemoveSubscriberTag adds a tag to a subscriber.
func (c *Client) RemoveSubscriberTag(req *TagReq) (*Response, error) {
	url := fmt.Sprintf("%s/%s/subscribers/%s/tags/%s", baseURL, c.accountID, req.Email, req.Tag)
	httpReq, err := c.getReq(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	resp.StatusCode = httpResp.StatusCode
	err = c.decodeResp(httpResp, resp)
	return resp, err
}
