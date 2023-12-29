package omglol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	timeout = 60 * time.Second
)

type DNSRecord struct {
	ID        int    `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	Data      string `json:"data,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type listResponse struct {
	Message string
	Dns     []DNSRecord
}

type listOutput struct {
	Response listResponse
}

type createReceived struct {
	Data DNSRecord `json:"data"`
}

type createResponse struct {
	ResponseReceived createReceived `json:"response_received"`
}

type createOutput struct {
	Response createResponse `json:"response"`
}

type Client struct {
	username string
	apiKey   string
}

func NewClient(username, apiKey string) (Client, error) {
	return Client{username: username, apiKey: apiKey}, nil
}

func (c Client) getRecords() (records []DNSRecord, err error) {
	client := http.Client{Timeout: timeout}
	endpoint := fmt.Sprintf("%s/address/%s/dns", baseUrl, c.username)
	authKey := fmt.Sprintf("Bearer %s", c.apiKey)
	reqUrl, err := url.Parse(endpoint)
	if err != nil {
		return records, err
	}

	req := http.Request{Method: "GET", URL: reqUrl, Header: map[string][]string{"Authorization": {authKey}}}

	resp, err := client.Do(&req)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()

	var output listOutput
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return records, err
	}

	return output.Response.Dns, nil
}

func (c Client) GetRecord(id int) (DNSRecord, error) {
	records, err := c.getRecords()
	if err != nil {
		return DNSRecord{}, err
	}

	for _, record := range records {
		if id == record.ID {
			return record, nil
		}
	}

	return DNSRecord{}, fmt.Errorf("unable to find DNS record with ID %d", id)
}

func (c Client) CreateRecord(record DNSRecord) (DNSRecord, error) {
	client := http.Client{Timeout: timeout}
	endpoint := fmt.Sprintf("%s/address/%s/dns", baseUrl, c.username)
	authKey := fmt.Sprintf("Bearer %s", c.apiKey)
	reqUrl, err := url.Parse(endpoint)
	if err != nil {
		return record, err
	}

	m, err := json.Marshal(record)
	if err != nil {
		return record, err
	}
	bodyReader := io.NopCloser(strings.NewReader(string(m)))

	req := http.Request{
		Method: "POST",
		URL:    reqUrl,
		Header: map[string][]string{"Authorization": {authKey}},
		Body:   bodyReader,
	}

	resp, err := client.Do(&req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, subErr := io.ReadAll(resp.Body)
		if subErr != nil {
			return record, subErr
		}
		return record, fmt.Errorf("error creating record: %s", respBody)
	}

	var output createOutput
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return record, err
	}

	record.ID = output.Response.ResponseReceived.Data.ID
	return record, nil
}

func (c Client) doDeleteRecord(id int) error {
	client := http.Client{Timeout: timeout}
	endpoint := fmt.Sprintf("%s/address/%s/dns/%d", baseUrl, c.username, id)
	authKey := fmt.Sprintf("Bearer %s", c.apiKey)
	reqUrl, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	req := http.Request{
		Method: "DELETE",
		URL:    reqUrl,
		Header: map[string][]string{"Authorization": {authKey}},
	}

	resp, err := client.Do(&req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, subErr := io.ReadAll(resp.Body)
		if subErr != nil {
			return subErr
		}
		return fmt.Errorf("error deleting record: %s", respBody)
	}

	return nil
}

func (c Client) deleteRecord(id int) error {
	records, err := c.getRecords()
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.ID == id {
			subErr := c.doDeleteRecord(id)
			if subErr != nil {
				return subErr
			}
		}
	}

	return nil
}
