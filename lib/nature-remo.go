package mpnatureremo

import (
	"flag"
	"fmt"
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
	natureremo "github.com/papix/go-nature-remo/cloud"
)

type NatureRemoPlugin struct {
	Prefix      string
	AccessToken string
	Client      *natureremo.Client
}

func (nr NatureRemoPlugin) GraphDefinition() map[string]mp.Graphs {
	devices, err := nr.Client.GetDevices()
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

func (nr NatureRemoPlugin) FetchMetrics() (map[string]float64, error) {
	ret := map[string]float64{}

	devices, err := nr.Client.GetDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		ret[fmt.Sprintf("temperature.%s", device.Name)] = float64(device.NewestEvents.Temperature.Value)
		ret[fmt.Sprintf("humidity.%s", device.Name)] = float64(device.NewestEvents.Humidity.Value)
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

	client := natureremo.NewClient(*optAccessToken)
	nr := NatureRemoPlugin{
		Prefix:      *optPrefix,
		AccessToken: *optAccessToken,
		Client:      client,
	}

	plugin := mp.NewMackerelPlugin(nr)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
