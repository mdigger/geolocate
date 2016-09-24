package geolocate

type Request struct {
	// The mobile country code stored on the SIM card (100-999).
	HomeMobileCountryCode uint16 `json:"homeMobileCountryCode,omitempty"`
	// The mobile network code stored on the SIM card (0-32767).
	HomeMobileNetworkCode uint16 `json:"homeMobileNetworkCode,omitempty"`
	// The mobile radio type. Supported values are lte, gsm, cdma, and wcdma.
	RadioType string `json:"radioType,omitempty"`
	// The clear text name of the cell carrier / operator.
	Carrier string `json:"carrier,omitempty"`
	// Should the clients IP address be used to locate it, defaults to true.
	ConsiderIp bool `json:"considerIp"`
	// Array of cell towers
	CellTowers []CellTower `json:"cellTowers,omitempty"`
	// Array of wifi access points
	WifiAccessPoints []WifiAccessPoint `json:"wifiAccessPoints,omitempty"`
	// Client IP Address
	IPAddress string `json:"ipaddress,omitempty"`
	// The fallback section is a custom addition to the GLS API.
	Fallbacks *Fallbacks `json:"fallbacks,omitempty"`
}

type Fallbacks struct {
	// If no exact cell match can be found, fall back from exact cell position
	// estimates to more coarse grained cell location area estimates, rather
	// than going directly to an even worse GeoIP based estimate.
	LAC bool `json:"lacf"`
	// If no position can be estimated based on any of the provided data points,
	// fall back to an estimate based on a GeoIP database based on the senders
	// IP address at the time of the query.
	IP bool `json:"ipf"`
}

type CellTower struct {
	// The mobile country code.
	MobileCountryCode uint16 `json:"mobileCountryCode"`
	// The mobile network code.
	MobileNetworkCode uint16 `json:"mobileNetworkCode"`
	// The location area code for GSM and WCDMA networks. The tracking area code
	// for LTE networks.
	LocationAreaCode uint16 `json:"locationAreaCode"`
	// The cell id or cell identity.
	CellId uint32 `json:"cellId"`
	// The signal strength for this cell network, either the RSSI or RSCP.
	SignalStrength int16 `json:"signalStrength,omitempty"`
	// The number of milliseconds since this networks was last detected.
	Age uint32 `json:"age,omitempty"`
	// The timing advance value for this cell network.
	TimingAdvance uint8 `json:"timingAdvance,omitempty"`
}

type WifiAccessPoint struct {
	// The BSSID of the WiFi network.
	MacAddress string `json:"macAddress"`
	// The received signal strength (RSSI) in dBm.
	SignalStrength int16 `json:"signalStrength,omitempty"`
	// The number of milliseconds since this network was last detected.
	Age uint32 `json:"age,omitempty"`
	// The WiFi channel, often 1 - 13 for networks in the 2.4GHz range.
	Channel uint8 `json:"channel,omitempty"`
	// The current signal to noise ratio measured in dB.
	SignalToNoiseRatio uint16 `json:"signalToNoiseRatio,omitempty"`
}

type Response struct {
	// The user’s estimated latitude and longitude, in degrees.
	Location Point `json:"location"`
	// The accuracy of the estimated location, in meters.
	Accuracy float64 `json:"accuracy"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lng"`
}
