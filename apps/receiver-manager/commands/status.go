package commands

import (
	"net/http"
    "io/ioutil"
)


//type Status struct {
//	State string
//}



func Status(server string) (state string) {
	response, err := http.Get(server+"/prog/Show?InstallFirmwareStatus")
    if err != nil {
		return ""
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
			return ""
        }
		//result := strings.Split(string(contents),"\n")
	return string(contents)
	}
}
