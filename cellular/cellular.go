package main

import (
	"fmt"
	"log"
	"strings"
//        "net"
)

import (
        "github.com/ozym/zone"
)

func main() {

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
//        for _, e := range list {
//               log.Println(e)
//        }

        // all Trimbles ...
        fmt.Println("--- Trimble ---")
        list, err = z.MatchByModel("^Hongdian")
        if err != nil {
                log.Fatal(err)
        }
        for _, e := range list {
//                fmt.Println(e)
//				fmt.Println(e.Name)
//				fmt.Println(e.Place)
				s := (strings.Replace(e.Place, " ", "-", -1))
//				fmt.Println(e.IP,"\n")
				fmt.Printf("%s %s #\n", e.IP, s)
        }
/*
        log.Println("--- Dump gps-waihorastation.wan.geonet.org.nz. ---")
        if d != nil {
                log.Printf("Name: %s\n", d.Name)
                log.Printf("Place: %s\n", d.Place)
                log.Printf("IP: %s\n", d.IP)
                log.Printf("Model: %s\n", d.Model)
                log.Printf("Code: %s\n", d.Code)
                log.Printf("Latitude: %g\n", d.Latitude)
                log.Printf("Longitude: %g\n", d.Longitude)
                log.Printf("Height: %g\n", d.Height)
        }
*/
}