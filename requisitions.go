package nordigen

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const requisitionsPath = "requisitions"

type Requisitions struct {
	Count    int64         `json:"count,omitempty"`
	Next     string        `json:"next,omitempty"`
	Previous string        `json:"previous,omitempty"`
	Results  []Requisition `json:"results,omitempty"`
}

type Requisition struct {
	Id       string    `json:"id,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Redirect string    `json:"redirect,omitempty"`
	Status   string    `json:"status,omitempty"`
	// There is an issue in the api, the status is still a string
	// like in v1
	//Status        Status    `json:"status,omitempty"`
	InstitutionId string   `json:"institution_id,omitempty"`
	Agreement     string   `json:"agreement,omitempty"`
	Reference     string   `json:"reference,omitempty"`
	Accounts      []string `json:"accounts,omitempty"`
	UserLanguage  string   `json:"user_language,omitempty"`
	Link          string   `json:"link,omitempty"`
}

type Status struct {
	Short       string `json:"short,omitempty"`
	Long        string `json:"long,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *Client) CreateRequisition(ctx context.Context, r Requisition) (Requisition, error) {
	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Path: strings.Join([]string{requisitionsPath, ""}, "/"),
		},
	}

	req = req.WithContext(ctx)
	data, err := json.Marshal(r)

	if err != nil {
		return Requisition{}, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	resp, err := c.c.Do(req)

	if err != nil {
		return Requisition{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Requisition{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Requisition{}, &APIError{resp.StatusCode, string(body), err}
	}
	err = json.Unmarshal(body, &r)

	if err != nil {
		return Requisition{}, err
	}

	return r, nil
}

func (c *Client) GetRequisition(ctx context.Context, id string) (r Requisition, err error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{requisitionsPath, id, ""}, "/"),
		},
	}

	req = req.WithContext(ctx)
	resp, err := c.c.Do(req)

	if err != nil {
		return Requisition{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Requisition{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Requisition{}, &APIError{resp.StatusCode, string(body), err}
	}
	err = json.Unmarshal(body, &r)

	if err != nil {
		return Requisition{}, err
	}

	return r, nil
}

func (c *Client) GetRequisitions(ctx context.Context, limit int64, offset int64) (r Requisitions, err error) {
	url := &url.URL{
		Path: strings.Join([]string{requisitionsPath, ""}, "/"),
	}
	queryParams := url.Query()
	queryParams.Add("limit", strconv.FormatInt(limit, 10))
	queryParams.Add("offset", strconv.FormatInt(offset, 10))
	url.RawQuery = queryParams.Encode()

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}
	req = req.WithContext(ctx)
	resp, err := c.c.Do(req)

	if err != nil {
		return Requisitions{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Requisitions{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Requisitions{}, &APIError{resp.StatusCode, string(body), err}
	}
	err = json.Unmarshal(body, &r)

	if err != nil {
		return Requisitions{}, err
	}

	return r, nil
}
