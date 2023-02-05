package omglol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	timeout = 60 * time.Second
)

// dnsRecordIdStr represent records returned from read-only endpoints.
type dnsRecordIdStr struct {
	Id        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	Data      string `json:"data,omitempty"`
	TTL       string `json:"ttl,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type listResponse struct {
	Message string
	Dns     []dnsRecordIdStr
}

type listOutput struct {
	Response listResponse
}

type dnsRecord struct {
	Id        int    `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	Data      string `json:"data,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type createReceived struct {
	Data dnsRecord `json:"data"`
}

type createResponse struct {
	ResponseReceived createReceived `json:"response_received"`
}

type createOutput struct {
	Response createResponse `json:"response"`
}

func getRecords(a auth) (records []dnsRecord, err error) {
	client := http.Client{Timeout: timeout}
	endpoint := fmt.Sprintf("%s/address/%s/dns", baseUrl, a.username)
	authKey := fmt.Sprintf("Bearer %s", a.apiKey)
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

	for _, record := range output.Response.Dns {
		id, subErr := strconv.Atoi(record.Id)
		if subErr != nil {
			return records, err
		}

		ttl, subErr := strconv.Atoi(record.TTL)
		if subErr != nil {
			return records, err
		}

		canonical := dnsRecord{
			Id:        id,
			Type:      record.Type,
			Name:      record.Name,
			Data:      record.Data,
			TTL:       ttl,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}
		records = append(records, canonical)
	}
	return
}

func getRecord(a auth, id int) (dnsRecord, error) {
	records, err := getRecords(a)
	if err != nil {
		return dnsRecord{}, err
	}

	for _, record := range records {
		if id == record.Id {
			return record, nil
		}
	}

	return dnsRecord{}, fmt.Errorf("unable to find DNS record with ID %s", id)
}

func createRecord(a auth, record dnsRecord) (dnsRecord, error) {
	client := http.Client{Timeout: timeout}
	endpoint := fmt.Sprintf("%s/address/%s/dns", baseUrl, a.username)
	authKey := fmt.Sprintf("Bearer %s", a.apiKey)
	reqUrl, err := url.Parse(endpoint)
	if err != nil {
		return record, err
	}

	m, err := json.Marshal(record)
	if err != nil {
		return record, err
	}
	bodyReader := ioutil.NopCloser(strings.NewReader(string(m)))

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
		respBody, subErr := ioutil.ReadAll(resp.Body)
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

	record.Id = output.Response.ResponseReceived.Data.Id
	return record, nil
}

func doDeleteRecord(a auth, id int) error {
	client := http.Client{Timeout: timeout}
	endpoint := fmt.Sprintf("%s/address/%s/dns/%d", baseUrl, a.username, id)
	authKey := fmt.Sprintf("Bearer %s", a.apiKey)
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
		respBody, subErr := ioutil.ReadAll(resp.Body)
		if subErr != nil {
			return subErr
		}
		return fmt.Errorf("error deleting record: %s", respBody)
	}

	return nil
}

func deleteRecord(a auth, id int) error {
	records, err := getRecords(a)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Id == id {
			subErr := doDeleteRecord(a, id)
			if subErr != nil {
				return subErr
			}
		}
	}

	return nil
}
