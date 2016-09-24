package geolocate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

var (
	// the maximum waiting time for a response
	RequestTimeout = time.Second * 10
	// do not use definition of by IP address
	IgnoreIPMethod = true
)

// Error returned when a data request is standard service of geolocation.
var (
	// invalid data format of the request or a bad key
	ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
	// over the limit of queries
	ErrForbidden = errors.New(http.StatusText(http.StatusForbidden))
	// information not found
	ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))
)

// The URL of the location services.
const (
	Mozilla = "https://location.services.mozilla.com/v1/geolocate"
	Google  = "https://www.googleapis.com/geolocation/v1/geolocate"
	Yandex  = "http://api.lbs.yandex.net/geolocation"
)

// Locator describes the interface supported by all types of location services.
type Locator interface {
	Get(req Request) (*Response, error)
}

// base describes information about the geolocation service, using standard
// types of queries, such as Mozilla and Google Locator.
type base struct {
	serviceUrl string       // адрес для запроса сервиса
	client     *http.Client // HTTP-клиент
}

// New returns a new initialized the geolocation service.
func New(serviceUrl, apiKey string) (locator Locator, err error) {
	if serviceUrl == Yandex { // для Яндекса возвращаем отдельный обрабочик
		return &yandex{
			apiKey: apiKey,
			client: &http.Client{
				Timeout: RequestTimeout,
			},
		}, nil
	}
	if apiKey != "" {
		serviceUrl += "?key=" + url.QueryEscape(apiKey)
	}
	if _, err := url.ParseRequestURI(serviceUrl); err != nil {
		return nil, err
	}
	return &base{
		serviceUrl: serviceUrl,
		client: &http.Client{
			Timeout: RequestTimeout,
		},
	}, nil
}

// Get transmits data to the server geolocation and returns it parsed from the
// response or the error.
func (l *base) Get(req Request) (*Response, error) {
	req.ConsiderIp = !IgnoreIPMethod
	if IgnoreIPMethod {
		req.Fallbacks = &Fallbacks{
			LAC: false,
			IP:  false,
		}
	}
	if req.RadioType == "" {
		req.RadioType = "gsm" // Mozilla does not find the data if not specified
	}
	req.IPAddress = "" // not used in the queries
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := l.client.Post(l.serviceUrl, "application/json",
		bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200: // all is well — data obtained
	case 400: // invalid data format of the request or a bad key
		return nil, ErrBadRequest
	case 403: // over the limit of queries
		return nil, ErrForbidden
	case 404: // information not found
		return nil, ErrNotFound
	default: // another bad error
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &response, nil
}
