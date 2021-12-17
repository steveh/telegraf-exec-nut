package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	nut "github.com/robbiet480/go.nut"
)

type void struct{}

const (
	measurement = "nut"

	precision = time.Nanosecond
)

var (
	host     string
	port     int
	username string
	password string

	entry void

	desiredTags = map[string]void{
		"battery.type":  entry,
		"device.mfr":    entry,
		"device.model":  entry,
		"device.serial": entry,
		"device.type":   entry,
	}

	desiredValues = map[string]void{
		"battery.charge":           entry,
		"battery.charge.low":       entry,
		"battery.runtime":          entry,
		"input.frequency":          entry,
		"input.transfer.high":      entry,
		"input.transfer.low":       entry,
		"input.voltage":            entry,
		"output.frequency":         entry,
		"output.frequency.nominal": entry,
		"output.voltage":           entry,
		"output.voltage.nominal":   entry,
		"ups.beeper.status":        entry,
		"ups.delay.shutdown":       entry,
		"ups.delay.start":          entry,
		"ups.firmware":             entry,
		"ups.load":                 entry,
		"ups.power":                entry,
		"ups.power.nominal":        entry,
		"ups.realpower":            entry,
		"ups.status":               entry,
		"ups.timer.shutdown":       entry,
		"ups.timer.start":          entry,
	}
)

func init() {
	flag.StringVar(&host, "host", "localhost", "NUT host")
	flag.IntVar(&port, "port", 3493, "NUT port")
	flag.StringVar(&username, "username", "", "NUT username")
	flag.StringVar(&password, "password", "", "NUT password")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	client, err := nut.Connect(host, port)
	check(err)

	defer func() {
		_, err := client.Disconnect()
		check(err)
	}()

	if username != "" {
		_, err = client.Authenticate(username, password)
		check(err)
	}

	timestamp := time.Now().UTC()

	upsList, err := client.GetUPSList()
	check(err)

	for _, ups := range upsList {
		tags := map[string]string{
			"name": ups.Name,
		}

		values := make(map[string]interface{}, len(desiredValues))

		for _, v := range ups.Variables {
			if _, ok := desiredTags[v.Name]; ok {
				tags[v.Name] = v.Value.(string)
			} else if _, ok := desiredValues[v.Name]; ok {
				values[v.Name] = v.Value
			}
		}

		p := influxdb2.NewPoint(measurement, tags, values, timestamp)

		fmt.Print(write.PointToLineProtocol(p, precision))
	}
}
