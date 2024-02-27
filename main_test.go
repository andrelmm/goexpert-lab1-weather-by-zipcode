package main

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocation(t *testing.T) {
	cases := []struct {
		Zip             string
		ExpectedError   error
		ExpectedLogLine string
	}{
		{"12345678", errors.New("can not find zipcode"), ""},
		{"", errors.New("invalid zipcode"), ""},
		{"17011067", nil, ""},
	}

	for _, tc := range cases {
		t.Run(tc.Zip, func(t *testing.T) {
			logger := &testLogger{}

			respBody := `{"localidade": "São Paulo"}`
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(respBody))
			})
			ts := httptest.NewServer(h)
			defer ts.Close()

			viaCepAPI = ts.URL + "/ws/%s/json/"

			loc, err := getLocation(tc.Zip)

			if logger.lastLogLine != tc.ExpectedLogLine {
				t.Errorf("Expected log line: %s, got: %s", tc.ExpectedLogLine, logger.lastLogLine)
			}

			if err != nil && err.Error() != tc.ExpectedError.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.ExpectedError, err)
			}

			if loc != nil && loc.Location != "São Paulo" {
				t.Errorf("Expected location: São Paulo, got: %s", loc.Location)
			}
		})
	}
}

func TestGetTemperature(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current": {"temp_c": 25.3}}`))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()
	weatherAPI = ts.URL + "/weather?key=2abdcba66a8b4196b4402638242702&q=%s"

	temp, err := getTemperature("São Paulo")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if temp == nil || temp.Current.TempC != 25.3 {
		t.Errorf("Expected temperature 25.3, got: %v", temp.Current.TempC)
	}
}

type testLogger struct {
	lastLogLine string
}

func (l *testLogger) Printf(format string, v ...interface{}) {
	l.lastLogLine = format
}

func TestConvertTemperatures(t *testing.T) {
	temp := convertTemperatures(25.3)

	expectedFahrenheit := 77.54
	expectedKelvin := 298.45

	if round(temp.TempF, 2) != expectedFahrenheit {
		t.Errorf("Expected temperature in Fahrenheit: %f, got: %f", expectedFahrenheit, temp.TempF)
	}

	if round(temp.TempK, 2) != expectedKelvin {
		t.Errorf("Expected temperature in Kelvin: %f, got: %f", expectedKelvin, temp.TempK)
	}
}

func round(num float64, decimalPlaces int) float64 {
	precision := math.Pow(10, float64(decimalPlaces))
	return math.Round(num*precision) / precision
}
