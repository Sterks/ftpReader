package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/secsy/goftp"
)

const connection string = "postgres://postgres:596run49@localhost/postgres?sslmode=disable"

//Сохранение файлов на диск
func SaveFiles(connect *goftp.Client, pathSave string, value FileInfo) {
	os.MkdirAll(pathSave+"/"+dateTimeNowString(), 0755)
	pathLocalFile := pathSave + "/" + dateTimeNowString() + "/" + value.nameFile
	fmt.Println(pathLocalFile)
	filePath, err := os.Create(pathLocalFile)
	if err != nil {
		fmt.Println(err)
	}
	defer filePath.Close()
	connect.Retrieve(value.filepath, filePath)
	if err != nil {
		fmt.Println(err)
	}

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
func NewFileInfo(value FileInfo, hash string, saveFile bool) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `insert into "ArchFiles" ("nameFile", "area", "filepath", "size", "modeTime", "hash", "saveFile") values ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(sqlStatement, value.nameFile, value.area, value.filepath, value.size, value.modeTime, hash, saveFile)
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

func NeedFileOpen() {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	// db.QueryRow(`select * from \"ArchFiles\" where `)
}

func UnArchive(srcPath string, dstPasth, out chan string) {
	file, err := zip.OpenReader(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, f := range file.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.CopyN(os.Stdout, rc, 68)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		fmt.Println()
	}

}
