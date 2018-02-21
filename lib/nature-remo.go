package mpnatureremo

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type NatureRemoPlugin struct {
	Prefix      string
	AccessToken string
}

type Device struct {
	Name         string `json:"name"`
	NewestEvents struct {
		Hu struct {
			Val int `json:"val"`
		} `json:"hu"`
		Te struct {
			Val float64 `json:"val"`
		} `json:"te"`
	} `json:"newest_events"`
}

func (nr NatureRemoPlugin) GraphDefinition() map[string]mp.Graphs {
	devices, err := nr.fetchDevices()
	if err != nil {
		return nil
	}

	ret := map[string]mp.Graphs{}
	temperature := make([]mp.Metrics, len(devices))
	humidity := make([]mp.Metrics, len(devices))

	for i, device := range devices {
		temperature[i] = mp.Metrics{
			Name:  fmt.Sprintf("temperature.%s", device.Name),
			Label: "temperature",
		}
		humidity[i] = mp.Metrics{
			Name:  fmt.Sprintf("humidity.%s", device.Name),
			Label: "humidity",
		}
	}
	ret["temperature"] = mp.Graphs{
		Label:   "temperature",
		Unit:    mp.UnitFloat,
		Metrics: temperature,
	}
	ret["humidity"] = mp.Graphs{
		Label:   "humidity",
		Unit:    mp.UnitInteger,
		Metrics: humidity,
	}

	return ret
}

func (nr NatureRemoPlugin) fetchDevices() ([]Device, error) {
	req, err := http.NewRequest("GET", "https://api.nature.global/1/devices", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", nr.AccessToken))

	tr := &http.Transport{
		TLSNextProto: nil,
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var devices []Device
	err = json.Unmarshal(body, &devices)

	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (nr NatureRemoPlugin) FetchMetrics() (map[string]float64, error) {
	ret := map[string]float64{}

	devices, err := nr.fetchDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		ret[fmt.Sprintf("temperature.%s", device.Name)] = float64(device.NewestEvents.Te.Val)
		ret[fmt.Sprintf("humidity.%s", device.Name)] = float64(device.NewestEvents.Hu.Val)
	}

	return ret, nil
}

func (nr NatureRemoPlugin) MetricKeyPrefix() string {
	if nr.Prefix == "" {
		nr.Prefix = "NatureRemo"
	}
	return nr.Prefix
}

func Do() {
	optPrefix := flag.String("metric-key-prefix", "NatureRemo", "Metric key prefix")
	optAccessToken := flag.String("access-token", os.Getenv("NATURE_REMO_ACCESS_TOKEN"), "Access token")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	nr := NatureRemoPlugin{
		Prefix:      *optPrefix,
		AccessToken: *optAccessToken,
	}

	plugin := mp.NewMackerelPlugin(nr)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
