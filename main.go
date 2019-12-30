package main

import (
	"fmt"

	"github.com/Sterks/ftpReader/ftp"

	"github.com/Sterks/ftpReader/router"
	"github.com/jasonlvhit/gocron"
)

func main() {
	//go scheduler()
	go ftp.MainProccess()
	router.Start()
}

func scheduler() {
	gocron.Every(10).Seconds().Do(task)
	<-gocron.Start()
}

func task() {
	fmt.Println("text")
}

func init() {

}
