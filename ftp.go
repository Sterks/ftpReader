package main

import (
	"fmt"
	"github.com/secsy/goftp"
)

func main(){
	connect := connect()
	path := "/fcs_regions/Moskva/notifications"
	recFiles(connect, path)
}

func connect() *goftp.Client {
	config := goftp.Config{User: "free", Password: "free"}
	ftp, err := goftp.DialConfig(config,"ftp.zakupki.gov.ru:21")
	if err != nil{
		_ = fmt.Errorf("Блок - 1 %v", err)
	}
	return ftp
}

func recFiles(connect *goftp.Client, path string ){

	list, err := connect.ReadDir(path)
	if err != nil {
		_ = fmt.Errorf("recFiles %v", err)
	}
	for _, value := range list {
		if value.IsDir() == false {

		}
	}
}
