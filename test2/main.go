package main

import (
	"net/http"
	"strings"
	//	"log"

	//	"github.com/datadog-go/statsd"
)

func main() {
	body := strings.NewReader(`{ "series" :
         [{"metric":"ntp.offset",
          "points":[[$currenttime, 20]],
          "type":"gauge",
          "host":"avdps03.geonet.org.nz",
          "tags":["environment:test"]}
        ]
    }`)
	req, err := http.NewRequest("POST", "https://app.datadoghq.com/api/v1/series?api_key=6abee9c15db54a3ec292bddd28965f53", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	/*
		c, err := statsd.New("127.0.0.1:8125")
		if err != nil {
			log.Fatal(err)
		}
		// prefix every metric with the app name
		c.Namespace = "flubber."
		// send the EC2 availability zone as a tag with every metric
		c.Tags = append(c.Tags, "us-east-1a")
		err = c.Gauge("request.duration", 1.2, nil, 1)
		// ...
	*/
}
