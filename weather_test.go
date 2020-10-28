package netatmo

import (
	"fmt"
	"testing"
	"time"
)

type NetatmoMQTTStatus struct {
	Location          string  `json:"location"`
	Name              string  `json:"name"`
	BatteryPercentage int32   `json:"battery_percentage,omitempty"`
	WifiStatus        int32   `json:"wifi_status,omitempty"`
	RFStatus          int32   `json:"rf_status,omitempty"`
	Temperature       float32 `json:"temperature,omitempty"`
	Humidity          int32   `json:"humidity,omitempty"`
	Pressure          float32 `json:"pressure,omitempty"`
	AbsolutePressure  float32 `json:"absolute_pressure,omitempty"`
	Co2               int32   `json:"co2,omitempty"`
	Noise             int32   `json:"noise,omitempty"`
	LastUpdate        int64   `json:"lastupdate"`
}

func TestWeather(t *testing.T) {
	netatmo, err := New()
	if err != nil {
		panic(err)
	}

	dc, err := netatmo.Read()
	if err != nil {
		panic(err)
	}

	ct := time.Now().UTC().Unix()

	for _, station := range dc.Stations() {
		fmt.Printf("Station : %s\n", station.StationName)

		for _, module := range station.Modules() {
			fmt.Printf("\tModule : %s\n", module.ModuleName)

			update := NetatmoMQTTStatus{}

			update.Location = station.StationName
			update.Name = module.ModuleName

			{
				ts, data := module.Info()
				for dataName, value := range data {
					switch dataName {
					case "BatteryPercent":
						update.BatteryPercentage = value.(int32)
					case "WifiStatus":
						update.WifiStatus = value.(int32)
					case "RFStatus":
						update.RFStatus = value.(int32)
					}
					secs := ct - ts
					update.LastUpdate = secs
					fmt.Printf("\t\t%s : %v (updated %ds ago)\n", dataName, value, ct-ts)
				}
			}

			{
				ts, data := module.Data()
				for dataName, value := range data {
					switch dataName {
					case "Temperature":
						update.Temperature = value.(float32)
					case "Humidity":
						update.Humidity = value.(int32)
					case "Pressure":
						update.Pressure = value.(float32)
					case "AbsolutePressure":
						update.AbsolutePressure = value.(float32)
					case "CO2":
						update.Co2 = value.(int32)
					case "Noise":
						update.Noise = value.(int32)
					}
					secs := ct - ts
					update.LastUpdate = secs
					fmt.Printf("\t\t%s : %v (updated %ds ago)\n", dataName, value, ct-ts)
				}
			}

			fmt.Printf("update: %+v", update)
		}
	}

}
