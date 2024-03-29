package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ViaCepResponse struct {
	Location string `json:"localidade"`
}

type Temperature struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

var viaCepAPI = "https://viacep.com.br/ws/%s/json/"
var weatherAPI = "https://api.weatherapi.com/v1/current.json?key=2abdcba66a8b4196b4402638242702&q=%s"

func getLocation(zip string) (*ViaCepResponse, error) {
	log.Printf("Getting location for ZIP code: %s", zip)

	resp, err := http.Get(fmt.Sprintf(viaCepAPI, zip))
	if err != nil {
		log.Printf("Error occurred while making HTTP request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("HTTP request successful")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error occurred while reading response body: %v", err)
		return nil, err
	}

	log.Println("Response body read successfully")

	var viaCepResponse ViaCepResponse
	if err := json.Unmarshal(body, &viaCepResponse); err != nil {
		log.Printf("Error occurred while unmarshalling JSON: %v", err)
		return nil, err
	}

	log.Println("JSON unmarshalled successfully")

	if viaCepResponse.Location == "" {
		var errorResp map[string]string
		if err := json.Unmarshal(body, &errorResp); err != nil {
			log.Printf("Error occurred while unmarshalling error response JSON: %v", err)
			return nil, err
		}
		if val, ok := errorResp["erro"]; ok && val == "true" {
			log.Println("CEP not found")
			return nil, errors.New("CEP not found")
		}
	}

	log.Println("Location retrieved successfully")

	return &viaCepResponse, nil
}

func getTemperature(location string) (*WeatherAPIResponse, error) {

	encodedLocation := url.QueryEscape(location)

	resp, err := http.Get(fmt.Sprintf(weatherAPI, encodedLocation))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherApiResponse WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherApiResponse); err != nil {
		return nil, err
	}

	return &weatherApiResponse, nil
}

func convertTemperatures(tempC float64) Temperature {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return Temperature{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	zip := r.URL.Query().Get("zip")
	if len(zip) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	cepResponse, err := getLocation(zip)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	temperatureResponse, err := getTemperature(cepResponse.Location)
	if err != nil {
		http.Error(w, "can not find weather information", http.StatusNotFound)
		return
	}

	temperatures := convertTemperatures(temperatureResponse.Current.TempC)

	_ = json.NewEncoder(w).Encode(temperatures)
}

func main() {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.HandleFunc("/weather", HandleRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
