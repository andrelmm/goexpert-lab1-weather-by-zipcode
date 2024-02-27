## Weather by Zip Code

This Go application receives a valid Brazilian ZIP code (CEP), identifies the city, and returns the current weather temperature in Celsius, Fahrenheit, and Kelvin. It utilizes the ViaCEP API to find the location based on the ZIP code and the WeatherAPI to retrieve the current weather temperature.

### How to Use

To use this application, make a GET request to the following endpoint:

```http request
https://lab-1-weather-by-zipcode-qdowlklj5q-uc.a.run.app/weather?zip={zip_code}
```
Replace `{zip_code}` with the desired Brazilian ZIP code (CEP). For example:

```http request
https://lab-1-weather-by-zipcode-qdowlklj5q-uc.a.run.app/weather?zip=78175000
```

The response will be in the following format:

```json
{"temp_C":25.3,"temp_F":77.53999999999999,"temp_K":298.45}
```

Where:

- temp_C is the temperature in Celsius.
- temp_F is the temperature in Fahrenheit.
- temp_K is the temperature in Kelvin.

### Accessing the Published Version
You can access the published version of this application using the following URL:

[https://lab-1-weather-by-zipcode-qdowlklj5q-uc.a.run.app](https://lab-1-weather-by-zipcode-qdowlklj5q-uc.a.run.app)
