package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
	"github.com/secsy/goftp"
)

// FileInfo ...
type FileInfo struct {
	ARId          int64
	nameFile      string
	filepath      string
	size          int64
	modeTime      time.Time
	area          string
	hash          string
	saveFile      bool
	ext           string
	localFilePath string
	arparent      int
	mode          os.FileMode
	sys           interface{}
	checkDir      bool
}

// Waiter что-то
type Waiter struct {
}

// FileInfoSQL ...
type FileInfoSQL struct {
	ARId          int64
	nameFile      sql.NullString
	filepath      sql.NullString
	size          sql.NullInt64
	modeTime      time.Time
	area          sql.NullString
	hash          sql.NullString
	saveFile      bool
	ext           sql.NullInt64
	localFilePath sql.NullString
	arparent      sql.NullInt64
}

var err error

func main() {
	fmt.Println("Основной процесс запущен, через 2 мин запуститься задача")
	// gocron.Every(2).Minutes().Do(mainProccess)
	// <-gocron.Start()
	fmt.Println("Основной процесс завершен")
	mainProccess()
}

func mainProccess() {
	var infoFileMass []FileInfo
	fmt.Println("Запуск программы ...")
	start := time.Now()
	client := connect()

	now := time.Now()
	y, m, d := now.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Now()
	listInfoFiles := GetFiles(client, from, to, "notifications", infoFileMass)
	fmt.Println(len(listInfoFiles))
	for _, info := range listInfoFiles {
		if FindHash(info.hash) == false {
			SaveFiles(client, "./Files", info)
		} else {
			fmt.Printf("Запись в базе уже существует - %s - %s - %d\n", info.filepath, info.hash, info.size)
		}
	}

	listNotArchive := FindNotUnArch()
	openFile := "./Open/" + DateTimeNowString() + "/"
	for _, value := range listNotArchive {
		UnArchive(value.localFilePath.String, openFile)
	}
	fmt.Println(time.Since(start))
	fmt.Println("Выполнение завершено ...")
}

//GetFiles Извещения по 44
func GetFiles(client *goftp.Client, from time.Time, to time.Time, typeDocum string, infoFileMass []FileInfo) []FileInfo {
	rootPath := "/fcs_regions"
	listFiles, err := client.ReadDir(rootPath)
	if err != nil {
		log.Panic(err)
	}
	var lister []os.FileInfo
	for _, value := range listFiles {
		if value.IsDir() == true {
			lister = append(lister, value)
		}
	}
	var infoFileMassS []FileInfo
	for _, value2 := range lister {
		Walk(client, rootPath+"/"+value2.Name()+"/"+typeDocum, func(fullPath string, info os.FileInfo, err error) error {
			if err != nil {
				// no permissions is okay, keep walking
				if err.(goftp.Error).Code() == 550 {
					return nil
				}
				return err
			}
			// go func() {
			var infoFile FileInfo

			buf := new(bytes.Buffer)
			err = client.Retrieve(fullPath, buf)
			if err != nil {
				fmt.Println(err)
			}
			var hasher = sha256.New()
			_, err = io.Copy(hasher, buf)
			if err != nil {
				log.Panic(err)
			}
			hhh := hex.EncodeToString(hasher.Sum(nil))
			infoFile.nameFile = info.Name()
			infoFile.modeTime = info.ModTime()
			infoFile.size = info.Size()
			infoFile.hash = hhh
			infoFile.area = value2.Name()
			infoFile.filepath = fullPath
			hash := fmt.Sprintf("%s - %x", fullPath, hasher.Sum(nil))
			fmt.Println(hash)
			infoFileMassS = append(infoFileMassS, infoFile)
			return nil
		}, from, to)
	}
	fmt.Println(len(infoFileMassS))
	return infoFileMassS
}

// Walk ...
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
				if file.ModTime().After(from) && file.ModTime().Before(to) && file.IsDir() == false {
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

// Создание директории
func (f *FileInfo) createFolder(storeFiles string) string {
	t := time.Now().Local()
	s := t.Format("2006-01-02")
	path := storeFiles + "/" + s
	os.MkdirAll(path, 0755)
	return path
}

// Загрузка файлов
func (f *FileInfo) download(connect *goftp.Client, value FileInfo, storeFiles string) string {
	var hasher = md5.New()
	filename := storeFiles + "/" + value.nameFile
	fileName, _ := os.Create(filename)
	defer fileName.Close()
	//writer := bufio.NewWriter(fileName)
	hash := fmt.Sprintf("%s %x", filename, hex.EncodeToString(hasher.Sum(nil)))
	connect.Retrieve(value.filepath, fileName)
	return hash
}

func connect() *goftp.Client {
	config := goftp.Config{
		User:               "free",
		Password:           "free",
		ConnectionsPerHost: 10,
		// Timeout:            0 * time.Second,
		//Logger:             os.Stderr,
	}
	ftp, err := goftp.DialConfig(config, "ftp.zakupki.gov.ru:21")
	if err != nil {
		_ = fmt.Errorf("Блок - 1 %v", err)
	}
	return ftp
}

func checkFolder(connect *goftp.Client, infoFileMass []FileInfo, from time.Time, to time.Time, pathRoot string, rem string) []FileInfo {
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
				infoFile.area = rem
				infoFileMass = append(infoFileMass, infoFile)
			}
		} else {
			pd := pathRoot + "/" + value.Name()
			infoFileMass = checkFolder(connect, infoFileMass, from, to, pd, rem)
		}
	}
	return infoFileMass
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
					infoFileMass = append(infoFileMass, infoFile)
				}
			} else {
				pd := pathDir + "/" + value.Name()
				infoFileMass = checkFolder(connect, infoFileMass, from, to, pd, rem.Name())
			}
		}
	}
	return infoFileMass
}
