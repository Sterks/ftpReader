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

// //Получение директорий первого уровня
// func readFolder(connect *goftp.Client, pathRoot string) []os.FileInfo {

// 	listFiles, err := connect.ReadDir(pathRoot)
// 	var folder []os.FileInfo
// 	if err != nil {
// 		_ = fmt.Errorf("Ошибка в функции ReadDir - %v", err)
// 	}

// 	for _, value := range listFiles {
// 		if value.IsDir() == true {
// 			folder = append(folder, value)
// 		}
// 	}
// 	return folder
// }

// //Функция для рекурсивного перебора
// func recReadAllFiles(pathRoot string, connect *goftp.Client, from time.Time, to time.Time, quantityFiles int32, allSize int64, downloadFile bool, countFiles int, c chan int) {
// 	err := Walk(connect, pathRoot, func(fullPath string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			// no permissions is okay, keep walking
// 			if err.(goftp.Error).Code() == 550 {
// 				return nil
// 			}
// 			_ = fmt.Errorf("%v", err)
// 			return err
// 		}
// 		allSize = allSize + info.Size()
// 		quantityFiles++
// 		countFiles++
// 		if downloadFile == true && info.IsDir() == false {
// 			dateFolder := time.Now().Format("2006-01-02")
// 			_ = os.MkdirAll("./Files/"+dateFolder, 0777)
// 			destFile, _ := os.Create("./Files" + "/" + dateFolder + "/" + info.Name())
// 			_ = connect.Retrieve(fullPath, destFile)
// 		}
// 		return nil
// 	}, from, to)
// 	if err != nil {
// 		_ = fmt.Errorf("%v", err)
// 	}
// 	c <- int(quantityFiles)
// 	fmt.Println(pathRoot)
// 	fmt.Println("Общий размер файлов: ", allSize)
// 	fmt.Printf("Кол-во файлов за период с %v по %v составляет: %v \n", from.Format("2006-01-02"), to.Format("2006-01-02"), quantityFiles)
// 	fmt.Println("----------------------")
// }

// //Общая функция для обработки файлов
// func Walk(client *goftp.Client, root string, walkFn filepath.WalkFunc, from time.Time, to time.Time) (ret error) {
// 	dirsToCheck := make(chan string, 100)
// 	var workCount int32 = 1
// 	dirsToCheck <- root
// 	for dir := range dirsToCheck {
// 		go func(dir string) {
// 			files, err := client.ReadDir(dir)
// 			if err != nil {
// 				if err = walkFn(dir, nil, err); err != nil && err != filepath.SkipDir {
// 					ret = err
// 					close(dirsToCheck)
// 					return
// 				}
// 			}
// 			for _, file := range files {
// 				if file.ModTime().After(from) && file.ModTime().Before(to) {
// 					if err = walkFn(path.Join(dir, file.Name()), file, nil); err != nil {
// 						if file.IsDir() && err == filepath.SkipDir {
// 							continue
// 						}
// 						ret = err
// 						close(dirsToCheck)
// 						return
// 					}
// 				}

// 				if file.IsDir() {
// 					atomic.AddInt32(&workCount, 1)
// 					dirsToCheck <- path.Join(dir, file.Name())
// 				}
// 			}
// 			atomic.AddInt32(&workCount, -1)
// 			if workCount == 0 {
// 				close(dirsToCheck)
// 			}
// 		}(dir)
// 	}
// 	return ret
// }
