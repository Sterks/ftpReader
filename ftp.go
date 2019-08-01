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
	// Переменные необходимые для функции recReadAllFiles
	//connect := connect()
	//pathRoot := "/fcs_regions/Moskva/notifications"
	//from, _ := time.Parse("2006-01-02 15:04:05", "2019-08-01 00:00:00")
	//to, _ := time.Parse("2006-01-02 15:04:05", "2019-08-01 16:00:00")
	//var allSize int64
	//var quantityFiles int32
}

func connect() *goftp.Client {
	config := goftp.Config{
		User:               "free",
		Password:           "free",
		ConnectionsPerHost: 10,
		Timeout:            2000 * time.Second,
		Logger:             os.Stderr,
	}
	ftp, err := goftp.DialConfig(config, "ftp.zakupki.gov.ru:21")
	if err != nil {
		_ = fmt.Errorf("Блок - 1 %v", err)
	}
	return ftp
}

//Функция для рекурсивного перебора
func recReadAllFiles(pathRoot string, connect *goftp.Client, from time.Time, to time.Time, quantityFiles int32, allSize int64) {
	err := Walk(connect, pathRoot, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			// no permissions is okay, keep walking
			if err.(goftp.Error).Code() == 550 {
				return nil
			}
			_ = fmt.Errorf("%v", err)
			return err
		}
		allSize = allSize + info.Size()
		quantityFiles++
		fmt.Printf("%v | %v | %v \n", fullPath, info.Size(), info.ModTime())
		return nil
	}, from, to)
	if err != nil {
		_ = fmt.Errorf("%v", err)
	}
	fmt.Println("Общий размер файлов: ", allSize)
	fmt.Printf("Кол-во файлов за период с %v по %v составляет: %v \n", from.Format("2006-01-02"), to.Format("2006-01-02"), quantityFiles)
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
