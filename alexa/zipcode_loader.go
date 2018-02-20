package alexa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// ErrUnknownZipcode is set when the Alexa API returns 204 staus code saying
// that the device has no zipcode stored in their system
var ErrUnknownZipcode = errors.New("unknown zipcode")

// ErrNotUS is an error indicating that the requeting user does not
// reside in the USA
var ErrNotUS = errors.New("this API only works for devices in the US")

// Client is an api for interacting with the Alexa HTTP Client
type Client struct {
	defaultZip string
	client     *http.Client
}

// API loads zipcodes for devices
type API interface {
	DeviceZip(ctx context.Context, deviceID, apiAccessToken string) (string, error)
}

// NewAPI creates a new API that connects to the Alexa API for retrieving
// zipcodes for deviceIDs provided with requests
func NewAPI() *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 2,
		},
	}
}

// DeviceZip gets the zipcode from the amazon alexa request
// using the API instructions found here
// https://developer.amazon.com/docs/custom-skills/device-address-api.html
//
// Accept: application/json
// Authorization: Bearer {apiAccessToken}
// GET https://api.amazonalexa.com/v1/devices/{deviceId}/settings/address/countryAndPostalCode
func (api *Client) DeviceZip(ctx context.Context, deviceID, apiAccessToken string) (string, error) {
	url := fmt.Sprintf("https://api.amazonalexa.com/v1/devices/%s/settings/address/countryAndPostalCode", deviceID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate HTTP request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiAccessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := api.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "failed to send request")
	}

	switch resp.StatusCode {
	case 200:
		break
	case 204:
		return "", ErrUnknownZipcode
	default:
		body, _ := ioutil.ReadAll(resp.Body)
		return "", errors.Errorf("received bad status code (%d) from alexa API: %s", resp.StatusCode, body)
	}

	var res zipcodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", errors.Wrap(err, "failed to parse response body from alexa API")
	}

	if res.CountryCode != "US" {
		return "", ErrNotUS
	}

	return res.PostalCode, nil
}

type zipcodeResponse struct {
	CountryCode string `json:"countryCode"`
	PostalCode  string `json:"postalCode"`
}
