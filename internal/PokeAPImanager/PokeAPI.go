package PokeAPImanager

import (
	"encoding/json"
	"io"
	"net/http"
)

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(URL string) ([]byte, error) {
	res, err := http.Get(URL)

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func UnmarshalJson(data []byte) ([]Locations, error) {

	var locations []Locations

	if err := json.Unmarshal(data, &locations); err != nil {
		return []Locations{}, err
	}

	return locations, nil
}
