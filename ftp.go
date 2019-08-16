package main

import (
	"fmt"
	"os"
	"time"

	"github.com/secsy/goftp"
)

type FileInfo struct {
	nameFile string
	area     string
	filepath string
	size     int64
	checkDir bool
	modeTime time.Time
	mode     os.FileMode
	sys      interface{}
}

var infoFileMass []FileInfo

func main() {
	now := time.Now()
	y, m, d := now.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Now()
	connect := connect()
	// pathRoot := "/fcs_regions"
	t := readNotification44(connect, infoFileMass, from, to)
	fmt.Println(len(t))
}

func checkFolder(connect *goftp.Client, infoFileMass []FileInfo, from time.Time, to time.Time, pathRoot string) []FileInfo {
	var infoFile FileInfo
	files, _ := connect.ReadDir(pathRoot)
	for _, value := range files {
		if value.IsDir() == false {
			if value.ModTime().After(from) && value.ModTime().Before(to) {
				fullpathFile := pathRoot + "/" + value.Name()
				fmt.Println(fullpathFile)
				infoFile.nameFile = value.Name()
				infoFile.filepath = fullpathFile
				infoFile.size = value.Size()
				infoFile.checkDir = value.IsDir()
				infoFile.modeTime = value.ModTime()
				infoFile.mode = value.Mode()
				infoFile.sys = value.Sys()
				infoFileMass = append(infoFileMass, infoFile)
			}
		} else {
			pd := pathRoot + "/" + value.Name()
			infoFileMass = checkFolder(connect, infoFileMass, from, to, pd)
		}
	}
	return infoFileMass
}

func connect() *goftp.Client {
	config := goftp.Config{
		User:               "free",
		Password:           "free",
		ConnectionsPerHost: 10,
		Timeout:            2000 * time.Second,
		//Logger:             os.Stderr,
	}
	ftp, err := goftp.DialConfig(config, "ftp.zakupki.gov.ru:21")
	if err != nil {
		_ = fmt.Errorf("Блок - 1 %v", err)
	}
	return ftp
}

//Извещения по 44 с загрузкой
func readNotification44(connect *goftp.Client, infoFileMass []FileInfo, from time.Time, to time.Time) []FileInfo {
	pathR := "/fcs_regions"
	listAllFiles, _ := connect.ReadDir(pathR)
	var listFolder []os.FileInfo
	for _, value := range listAllFiles {
		if value.IsDir() == true {
			listFolder = append(listFolder, value)
		}
	}
	var infoFile FileInfo
	for _, rem := range listFolder {
		pathDir := pathR + "/" + rem.Name() + "/" + "notifications"
		listFiles, _ := connect.ReadDir(pathDir)
		for _, value := range listFiles {
			if value.IsDir() == false {
				if value.ModTime().After(from) && value.ModTime().Before(to) {
					fullpathFile := pathDir + "/" + value.Name()
					infoFile.nameFile = value.Name()
					infoFile.filepath = fullpathFile
					infoFile.size = value.Size()
					infoFile.checkDir = value.IsDir()
					infoFile.modeTime = value.ModTime()
					infoFile.mode = value.Mode()
					infoFile.sys = value.Sys()
					infoFile.area = rem.Name()
					fmt.Println(infoFile)
					infoFileMass = append(infoFileMass, infoFile)
				}
			} else {
				pd := pathDir + "/" + value.Name()
				infoFileMass = checkFolder(connect, infoFileMass, from, to, pd)
			}
		}
	}
	return infoFileMass
}
