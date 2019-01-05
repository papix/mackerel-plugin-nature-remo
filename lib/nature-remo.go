package mpnatureremo

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	natureremo "github.com/papix/go-nature-remo/cloud"
)

type NatureRemoPlugin struct {
	Prefix      string
	AccessToken string
	Client      *natureremo.Client
}

func (nr NatureRemoPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"temperature": {
			Label: "Temperature",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1"},
			},
		},
		"humidity": {
			Label: "Humidity",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1"},
			},
		},
		"illuminance": {
			Label: "Illluminance",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1"},
			},
		},
	}
}

var reg = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizeName(devName string) string {
	return reg.ReplaceAllString(strings.TrimSpace(devName), "_")
}

func (nr NatureRemoPlugin) FetchMetrics() (map[string]float64, error) {
	ret := map[string]float64{}

	devices, err := nr.Client.GetDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		devName := normalizeName(device.Name)
		evs := device.NewestEvents
		if evs.Temperature.CreatedAt != "" {
			ret[fmt.Sprintf("temperature.%s", devName)] = float64(evs.Temperature.Value)
		}
		if evs.Humidity.CreatedAt != "" {
			ret[fmt.Sprintf("humidity.%s", devName)] = float64(evs.Humidity.Value)
		}
		if evs.Illuminance.CreatedAt != "" {
			ret[fmt.Sprintf("illuminance.%s", devName)] = float64(evs.Illuminance.Value)
		}
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
