package rinex

import (
	"os"
	"io/ioutil"
	"path/filepath"
)

func RinexCheck(site string, y string)  string {
	basedir := "/archive/gps/nrtdata/2016"
	dayfiles, _ := ioutil.ReadDir(basedir)
    for _, d := range dayfiles {
		hourfiles, _ := ioutil.ReadDir(filepath.Join(basedir,d.Name()))
		for _, h := range  hourfiles {
			filename := site+d.Name()+"a.16d.Z"
			fullpath := filepath.Join("/archive/gps/nrtdata",y,d.Name(),h.Name(),filename)
			if _, err := os.Stat(fullpath); err != nil {
				return "file not found "+fullpath
			}
		}
    }
	return "Exit "

/*
	if _, err := os.Stat("/archive/gps/nrtdata/2016/033/00"); os.IsNotExist(err) {
		return "Dir Not Found"
	}
	if _, err := os.Stat("/archive/gps/nrtdata/2016/033/00/auck033a.16d.Z"); err != nil {
		return "file not found"
	}
*/	

}
