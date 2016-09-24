package geolocate_test

import (
	"fmt"
	"log"

	"github.com/mdigger/geolocate"
)

// Your API Keys
const (
	googleAPIKey = "<API Key>"
	yandexAPIKey = "<API Key>"
)

func Example_google() {
	request := geolocate.Request{
		CellTowers: []geolocate.CellTower{
			{250, 2, 7743, 22517, -78, 0, 0},
			{250, 2, 7743, 39696, -81, 0, 0},
			{250, 2, 7743, 22518, -91, 0, 0},
			{250, 2, 7743, 27306, -101, 0, 0},
			{250, 2, 7743, 29909, -103, 0, 0},
			{250, 2, 7743, 22516, -104, 0, 0},
			{250, 2, 7743, 20736, -105, 0, 0},
		},
		WifiAccessPoints: []geolocate.WifiAccessPoint{
			{"2:18:E4:C8:38:30", -22, 0, 0, 0},
		},
	}
	google, err := geolocate.New(geolocate.Google, googleAPIKey)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := google.Get(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Google:", resp)
}

func Example_mozilla() {
	request := geolocate.Request{
		CellTowers: []geolocate.CellTower{
			{250, 2, 7743, 22517, -78, 0, 0},
			{250, 2, 7743, 39696, -81, 0, 0},
			{250, 2, 7743, 22518, -91, 0, 0},
			{250, 2, 7743, 27306, -101, 0, 0},
			{250, 2, 7743, 29909, -103, 0, 0},
			{250, 2, 7743, 22516, -104, 0, 0},
			{250, 2, 7743, 20736, -105, 0, 0},
		},
		WifiAccessPoints: []geolocate.WifiAccessPoint{
			{"2:18:E4:C8:38:30", -22, 0, 0, 0},
		},
	}
	mozilla, err := geolocate.New(geolocate.Mozilla, "test")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := mozilla.Get(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mozilla:", resp)
}

func Example_yandex() {
	request := geolocate.Request{
		CellTowers: []geolocate.CellTower{
			{250, 2, 7743, 22517, -78, 0, 0},
			{250, 2, 7743, 39696, -81, 0, 0},
			{250, 2, 7743, 22518, -91, 0, 0},
			{250, 2, 7743, 27306, -101, 0, 0},
			{250, 2, 7743, 29909, -103, 0, 0},
			{250, 2, 7743, 22516, -104, 0, 0},
			{250, 2, 7743, 20736, -105, 0, 0},
		},
		WifiAccessPoints: []geolocate.WifiAccessPoint{
			{"2:18:E4:C8:38:30", -22, 0, 0, 0},
		},
	}
	yandex, err := geolocate.New(geolocate.Yandex, yandexAPIKey)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := yandex.Get(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Yandex:", resp)
}
