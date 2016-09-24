package geolocate

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// yandex describes the  handler location requests for Yandex.
type yandex struct {
	apiKey string       // the key to retrieve the data
	client *http.Client // HTTP client
}

// Get transmits data to Yandex and returns a parsed response.
func (l *yandex) Get(req Request) (*Response, error) {
	cells := make([]yandexCell, len(req.CellTowers))
	for i, cell := range req.CellTowers {
		cells[i] = yandexCell{
			MobileCountryCode: cell.MobileCountryCode,
			MobileNetworkCode: cell.MobileNetworkCode,
			CellId:            cell.CellId,
			LocationAreaCode:  cell.LocationAreaCode,
			SignalStrength:    cell.SignalStrength,
			Age:               cell.Age,
		}
	}
	wifis := make([]yandexWiFi, len(req.WifiAccessPoints))
	for i, wifi := range req.WifiAccessPoints {
		wifis[i] = yandexWiFi{
			MacAddress:     wifi.MacAddress,
			SignalStrength: wifi.SignalStrength,
			Age:            wifi.Age,
		}
	}
	request := yandexRequest{
		Common: yandexCommon{
			Version: "1.0",
			ApiKey:  l.apiKey,
		},
		GSMCells:     cells,
		WiFiNetworks: wifis,
	}
	if !IgnoreIPMethod && req.IPAddress != "" {
		request.IP = yandexIP{AddressV4: req.IPAddress}
	}
	data, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	// for Yandex it is transmitted in the form of a query of the form
	params := make(url.Values, 1)
	params.Set("json", string(data))
	resp, err := l.client.PostForm(Yandex, params)
	if err != nil {
		return nil, err
	}
	// according to the documentation Yandex does not return error codes, so it
	// is left just in case
	switch resp.StatusCode {
	case 200:
	case 400:
		return nil, ErrBadRequest
	case 403:
		return nil, ErrForbidden
	case 404:
		return nil, ErrNotFound
	default:
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	// decode response
	var yresp yandexRespose
	err = json.NewDecoder(resp.Body).Decode(&yresp)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// check the error text, and treated similarly to a public method of all the
	// locators
	switch yresp.Error {
	case "": // all is well — no error — ignore
	case "JSON request is invalid":
		return nil, ErrBadRequest
	case "invalid api_key":
		return nil, ErrForbidden
	case "Location not found":
		return nil, ErrNotFound
	default: // other text errors
		return nil, errors.New(yresp.Error)
	}
	// check that the method definitions are not by IP
	if IgnoreIPMethod && yresp.Position.Type == "ip" {
		return nil, ErrNotFound
	}
	response := Response{
		Location: Point{
			Lat: yresp.Position.Latitude,
			Lon: yresp.Position.Longitude,
		},
		Accuracy: yresp.Position.Precision,
	}
	return &response, nil
}

type yandexRequest struct {
	Common       yandexCommon `json:"common"`
	GSMCells     []yandexCell `json:"gsm_cells,omitempty"`
	WiFiNetworks []yandexWiFi `json:"wifi_networks,omitempty"`
	IP           yandexIP     `json:"ip,omitempty"`
}

type yandexIP struct {
	AddressV4 string `json:"address_v4"`
}

type yandexCell struct {
	MobileCountryCode uint16 `json:"countrycode"`
	MobileNetworkCode uint16 `json:"operatorid"`
	CellId            uint32 `json:"cellid"`
	LocationAreaCode  uint16 `json:"lac"`
	SignalStrength    int16  `json:"signal_strength,omitempty"`
	Age               uint32 `json:"age,omitempty"`
}
type yandexWiFi struct {
	MacAddress     string `json:"mac"`
	SignalStrength int16  `json:"signal_strength,omitempty"`
	Age            uint32 `json:"age,omitempty"`
}
type yandexCommon struct {
	Version string `json:"version"`
	ApiKey  string `json:"api_key"`
}

type yandexPosition struct {
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Altitude          float64 `json:"altitude"`
	Precision         float64 `json:"precision"`
	AltitudePrecision float64 `json:"altitude_precision"`
	Type              string  `json:"type"` // Method positioning: gsm, wifi, ip
}

type yandexRespose struct {
	Position yandexPosition `json:"position"`
	Error    string         `json:"error"`
}
