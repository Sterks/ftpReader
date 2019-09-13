package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/secsy/goftp"
)

const connection string = "postgres://postgres:596run49@localhost/postgres?sslmode=disable"

/*Парсинг XML*/
func UseNotification(pathFolder string) []string {
	files, err := ioutil.ReadDir(pathFolder)
	if err != nil {
		fmt.Println(err)
	}
	var mass []string
	for _, file := range files {
		name := file.Name()
		res := strings.Contains(name, "fcsNotificationEA")
		if res == true {
			pathFile := pathFolder + "/" + name
			mass = append(mass, pathFile)
		}
	}
	return mass
}

//Определить разрешение файла
func FileExt(path string) string {
	g := filepath.Ext(path)
	return g
}

//Дата в формате
func DateTimeNowString() string {
	t := time.Now().Local()
	s := t.Format("2006-01-02")
	return s
}

//Сохранение файлов на диск
func SaveFiles(connect *goftp.Client, pathSave string, value FileInfo) string {
	os.MkdirAll(pathSave+"/"+DateTimeNowString(), 0755)
	pathLocalFile := pathSave + "/" + DateTimeNowString() + "/" + value.nameFile
	fmt.Println(pathLocalFile)
	filePath, err := os.Create(pathLocalFile)
	if err != nil {
		fmt.Println(err)
	}
	defer filePath.Close()
	err = connect.Retrieve(value.filepath, filePath)
	if err != nil {
		fmt.Println(err)
	}
	// c <- pathLocalFile
	// if pathLocalFile == "" {
	// 	close(c)
	// }
	value.localFilePath = pathLocalFile
	NewFileInfo(value, true, FileExt(pathLocalFile), "N")
	return pathLocalFile
}

//Определяем Hash файла
func HashFiles(connect *goftp.Client, pathFile string, value FileInfo) string {
	hasher := md5.New()
	buf := new(bytes.Buffer)
	err = connect.Retrieve(pathFile, buf)
	if err != nil {
		fmt.Println(err)
	}
	io.Copy(hasher, buf)
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

//Создание новой записи
func NewFileInfo(value FileInfo, saveFile bool, ext string, unarch string) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStatement := `insert into "ArchFiles" ("nameFile", "area", "filepath", "size", "modeTime", "hash", "saveFile", "ext", unarch, "localPath") values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Exec(sqlStatement, value.nameFile, value.area, value.filepath, value.size, value.modeTime, value.hash, saveFile, ext, unarch, value.localFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

//Найти запись
func FindHash(hash string) bool {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var fiID int64
	_ = db.QueryRow(`select "ARId" from "ArchFiles" where "hash" = $1`, hash).Scan(&fiID)

	if fiID > 0 {
		return true
	} else {
		return false
	}
}

func ListIndex() {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select * from \"ArchFiles\"")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}

func FindNotUnArch() []*FileInfo {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`select * from "ArchFiles" where ext = '.zip' and unarch = 'N'`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	listFiles := make([]*FileInfo, 0)
	for rows.Next() {
		file := new(FileInfo)
		err := rows.Scan(&file.ARId, &file.nameFile, &file.filepath, &file.size, &file.modeTime, &file.area, &file.hash, &file.saveFile, &file.ext, &file.unarch, &file.localFilePath)
		if err != nil {
			log.Fatal(err)
		}
		listFiles = append(listFiles, file)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return listFiles
}

// func NeedFileOpen() {
// 	// db, err := sql.Open("postgres", connection)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// db.QueryRow(`select * from \"ArchFiles\" where `)
// }

func UnArchive(srcPath string, dstPasth string, out chan FileInfo) ([]string, error) {
	var filenames []string
	file, err := zip.OpenReader(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, f := range file.File {
		fpath := filepath.Join(dstPasth, f.Name)
		if filepath.Ext(fpath) != ".sig" {
			filenames = append(filenames, fpath)
			if f.FileInfo().IsDir() {
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			rc, err := f.Open()
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()

			if err != nil {
				return filenames, err
			}
		}
	}
	return filenames, nil
}
