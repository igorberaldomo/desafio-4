package main

import (
	"encoding/json"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"unicode"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}
type Weather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func main() {
	http.HandleFunc("/", FindTempHandler)
	http.ListenAndServe(":8080", nil)
}

func FindTempHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "No page found", http.StatusNotFound)
		return
	}

	cep := r.URL.Query().Get("cep")

	if len(cep) != 8 {
		http.Error(w, "invalid cep", http.StatusUnprocessableEntity)
		os.Exit(1)
		return
	}

	for _, r := range cep {
		if !unicode.IsNumber(r) {
			http.Error(w, "invalid cep", http.StatusUnprocessableEntity)
			os.Exit(1)
			return
		}
	}
	if len(cep) == 8 {
		// request para pegar localidade
		if len(cep) == 8 {
			req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
			if err != nil {
				http.Error(w, "error in requisition via CEP", http.StatusInternalServerError)
				return
			}
			switch req.StatusCode {
			case http.StatusBadRequest:
				{
					http.Error(w, "bad request", http.StatusBadRequest)
					os.Exit(1)
					return
				}
			case http.StatusNotFound:
				{
					http.Error(w, "cannot find zipcode", http.StatusNotFound)
					os.Exit(1)
					return
				}
			case http.StatusUnprocessableEntity:
				{
					http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
					os.Exit(1)
					return
				}
			}

			defer req.Body.Close()

			res, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, "error in reading the body via CEP", http.StatusInternalServerError)
				return
			}
			var data ViaCEP

			err = json.Unmarshal(res, &data)
			if err != nil {
				http.Error(w, "error in unmarshal via CEP", http.StatusInternalServerError)
				return
			}
			local := data.Localidade

			url := "http://api.weatherapi.com/v1/current.json?key=18525c8de5ac479f994185201250303&q=" + neturl.QueryEscape(local) + "&aqi=no"

			// novo request para pegar a temperatura
			req2, err2 := http.Get(url)
			if err2 != nil {
				http.Error(w, "error in requisition via WeatherAPI", http.StatusInternalServerError)
				return
			}
			defer req2.Body.Close()

			res2, err2 := io.ReadAll(req2.Body)
			if err2 != nil {
				http.Error(w, "error in reading the body via WeatherAPI", http.StatusInternalServerError)
				return
			}
			var data2 Weather
			err2 = json.Unmarshal(res2, &data2)
			if err2 != nil {
				http.Error(w, "error in unmarshal via WeatherAPI", http.StatusInternalServerError)
				return
			}

			tempC := data2.Current.TempC
			tempF := tempC*1.8 + 32
			tempK := tempC + 273

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"temp_C": tempC,
				"temp_F": tempF,
				"temp_K": tempK,
			})
		}
	}
}
