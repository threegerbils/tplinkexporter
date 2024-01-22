package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/thelastguardian/tplinkexporter/clients"
	"github.com/thelastguardian/tplinkexporter/collectors"
)

func main() {
	var (
		host       = kingpin.Flag("host", "Host of target tplink easysmart switch.").Required().String()
		username   = kingpin.Flag("username", "Username for switch GUI login").Default("admin").String()
		password   = kingpin.Flag("password", "Password for switch GUI login").Required().String()
		port       = kingpin.Flag("port", "Metrics port to listen on for prometheus scrapes").Default("9717").Int()
		switchType = kingpin.Flag("type", "Host of target tplink easysmart switch.").Default("1016DE").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	var trafficCollector *collectors.TrafficCollector = nil

	if *switchType == "1016DE" {
		tplinkSwitch := clients.NewTPLink1016DESwitch(*host, *username, *password)
		trafficCollector = collectors.NewTrafficCollector("tplinkexporter", tplinkSwitch)
	} else {
		tplinkSwitch := clients.NewTPLinkSwitch(*host, *username, *password)
		trafficCollector = collectors.NewTrafficCollector("tplinkexporter", tplinkSwitch)
	}
	prometheus.MustRegister(trafficCollector)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Beginning to serve on port :" + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
