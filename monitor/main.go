package main

import (
	"github.com/go-gps/rinex"
	"fmt"
	"flag"
	"os"
)
func main () {
	var site string
	flag.StringVar(&site, "site", "auck", "base site code")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Program to find the IP address of a GPS Site\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()
	var y string
	y = "2016"
	check := rinex.RinexCheck(site, y)
	fmt.Println(check)
}