package main

import (
	"fmt"
	"github.com/secsy/goftp"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"
)

func main() {
	connect := connect()
	pathRoot := "/fcs_regions"
	//recFiles(connect, pathRoot)
	from, _ := time.Parse("2006-01-02 15:04:05", "2019-03-01 00:00:00")
	to, _ := time.Parse("2006-01-02 15:04:05", "2019-04-10 00:00:00")
	_ = Walk(connect, pathRoot, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			// no permissions is okay, keep walking
			if err.(goftp.Error).Code() == 550 {
				return nil
			}
			return err
		}

		fmt.Println(info.ModTime())

		return nil
	}, from, to)
}

func connect() *goftp.Client {
	config := goftp.Config{
		User:               "free",
		Password:           "free",
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
	}
	ftp, err := goftp.DialConfig(config, "ftp.zakupki.gov.ru:21")
	if err != nil {
		_ = fmt.Errorf("Блок - 1 %v", err)
	}
	return ftp
}

func recFiles(connect *goftp.Client, path string) {

	list, err := connect.ReadDir(path)
	if err != nil {
		_ = fmt.Errorf("recFiles %v", err)
	}
	for _, value := range list {
		if value.IsDir() == false {
			from, _ := time.Parse("2006-01-02 15:04:05", "2019-03-01 00:00:00")
			to, _ := time.Parse("2006-01-02 15:04:05", "2019-04-10 00:00:00")
			if value.ModTime().After(from) && value.ModTime().Before(to) {

			}
		}
	}
}

func Walk(client *goftp.Client, root string, walkFn filepath.WalkFunc, from time.Time, to time.Time) (ret error) {
	dirsToCheck := make(chan string, 100)

	var workCount int32 = 1
	dirsToCheck <- root

	for dir := range dirsToCheck {
		go func(dir string) {
			files, err := client.ReadDir(dir)

			if err != nil {
				if err = walkFn(dir, nil, err); err != nil && err != filepath.SkipDir {
					ret = err
					close(dirsToCheck)
					return
				}
			}

			for _, file := range files {
				if file.ModTime().After(from) && file.ModTime().Before(to) {
					if err = walkFn(path.Join(dir, file.Name()), file, nil); err != nil {
						if file.IsDir() && err == filepath.SkipDir {
							continue
						}
						ret = err
						close(dirsToCheck)
						return
					}
				}

				if file.IsDir() {
					atomic.AddInt32(&workCount, 1)
					dirsToCheck <- path.Join(dir, file.Name())
				}
			}

			atomic.AddInt32(&workCount, -1)
			if workCount == 0 {
				close(dirsToCheck)
			}
		}(dir)
	}

	return ret
}
