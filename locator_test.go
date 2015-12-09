package geolocate

import (
	"fmt"
	"testing"
)

// Your API Keys
const (
	googleAPIKey = "<API Key>"
	yandexAPIKey = "<API Key>"
)

func TestLocator(t *testing.T) {
	google, err := New(Google, googleAPIKey)
	if err != nil {
		t.Fatal(err)
	}
	mozilla, err := New(Mozilla, "test")
	if err != nil {
		t.Fatal(err)
	}
	yandex, err := New(Yandex, yandexAPIKey)
	if err != nil {
		t.Fatal(err)
	}
	request := Request{
		CellTowers: []CellTower{
			{250, 2, 22517, 7743, -78, 0, 0},
			{250, 2, 39696, 7743, -81, 0, 0},
			{250, 2, 22518, 7743, -91, 0, 0},
			{250, 2, 27306, 7743, -101, 0, 0},
			{250, 2, 29909, 7743, -103, 0, 0},
			{250, 2, 22516, 7743, -104, 0, 0},
			{250, 2, 20736, 7743, -105, 0, 0},
		},
		WifiAccessPoints: []WifiAccessPoint{
			{"2:18:E4:C8:38:30", -22, 0, 0, 0},
		},
	}
	resp, err := google.Get(request)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Google:", resp)

	resp, err = mozilla.Get(request)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Mozilla:", resp)

	resp, err = yandex.Get(request)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Yandex:", resp)

}
