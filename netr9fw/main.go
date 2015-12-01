package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"github.com/delta/meta"
	"github.com/zone"
	"log"
    "io/ioutil"
	"net/http"
	"net"	
	"strings"
//	"bytes"
)

//type Fw interface {
//    Firmware() string
//}
//type FirmwareVersion struct {
//	FirmWare	string		`json:"FirmwareVersion"`      // 
//	Version     string      `json:"version"`        // 
//	Date 	    string      `json:"date"`   //
//}
//func (f FirmwareVersion) Firmware() string {
//	return "6000"
//}

func gv(ip net.IP) string {
	response, err := http.Get("http://admin:cuseeme@"+ip.String()+"/prog/Show?FirmwareVersion")
    if err != nil {
		return ""
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
			return ""
        }
		result := strings.Split(string(contents),"\n")
	    for i := range result {
			str := "version="
//			if strings.Contains(result[i],str) {
				x := strings.Split(result[i], " ")
				for j := range x {
//					str := "volts="
					if strings.Contains(x[j],str) {
					return string(x[j])
					}
				}
			}
		}
return ""
}


func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var data string
	flag.StringVar(&data, "data", "/home/davec/go/src/github.com/delta/data", "base data directory")

	var site string
	flag.StringVar(&site, "site", "TEST", "base site code")

	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "An example program to examine ...\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Parse()
	// load antenna details into a map ...
	antennamap := make(map[string]meta.InstalledAntenna)
	{
		var antennas meta.InstalledAntennas
		if err := meta.LoadList(filepath.Join(data, "antennas/trimble/installs.csv"), &antennas); err != nil {
			panic(err)
		}

		for _, n := range antennas {
			antennamap[n.Mark] = n
		}
	}
	// load station details
	markmap := make(map[string]meta.Mark)
	{
		var marks meta.Marks
		if err := meta.LoadList(filepath.Join(data, "marks.csv"), &marks); err != nil {
			panic(err)
		}

		for _, m := range marks {
			markmap[m.Code] = m
		}
	}
	//fmt.Println(netmap)
///////////////////////////////////////////////////////
	// sort the akeys on output
	var akeys []string
	for k, _ := range antennamap {
		akeys = append(akeys, k)
	}
	sort.Strings(akeys)
	// a simple loop and print
	for _, k := range akeys {
		v, ok := antennamap[k]
		if !ok {
			panic("invalid mark key: " + k)
		}
		if v.Mark == site {
			{
			n, ok := antennamap[v.Mark]
			if !ok {
				panic("unable to find network: " + v.Mark)
			}
			j, err := json.MarshalIndent(n, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
			}
		}
	}
/////_///////////////////////////////////////////////////
	// sort the keys on output
	var keys []string
	for k, _ := range markmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// a simple loop and print
	for _, k := range keys {
		v, ok := markmap[k]
		if !ok {
			panic("invalid mark key: " + k)
		}
		if v.Code == site {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
		}
	}
//###############
        // default geonet connection ....
        z := zone.Equipment{
                Zone:   "wan.geonet.org.nz.",
                Server: "rhubarb.geonet.org.nz",
                Port:   "53",
        }

        // get all equipment ...
        list, err := z.List()
        if err != nil {
                log.Fatal(err)
        }
        // all Trimble NetR9's
        list, err = z.MatchByModelAndCode("^Trimble",site)
        if err != nil {
                log.Fatal(err)
        }
		for _, e := range list {
		fw := gv(e.IP)
		fmt.Println(fw)
		fmt.Println(e.IP)
        }
//###############


}