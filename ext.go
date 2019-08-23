package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/secsy/goftp"
)

//Сохранение файлов на диск
func SaveFiles(pathSave string, value FileInfo, buf *bytes.Buffer) {
	os.MkdirAll(pathSave+"/"+dateTimeNowString(), 0755)
	pathLocalFile := pathSave + "/" + dateTimeNowString() + "/" + value.nameFile
	filePath, err := os.Create(pathLocalFile)
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(buf, filePath)
	if err != nil {
		fmt.Println(err)
	}
}

//Определяем Hash файла
func HashFiles(connect *goftp.Client, pathFile string, value FileInfo) (string, *bytes.Buffer) {
	hasher := md5.New()
	buf := new(bytes.Buffer)
	connect.Retrieve(pathFile, buf)
	io.Copy(hasher, buf)
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, buf
}

//Создание новой записи
func NewFileInfo(nameFile string, area string, filepath string, size int64, modeTime time.Time, hash string) {
	db, err := sql.Open("postgres", "postgres://postgres:596run49@localhost/EFtest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `insert into "FilesInfo" ("nameFile", "area", "filepath", "size", "modeTime", "hash") values ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(sqlStatement, nameFile, area, filepath, size, modeTime, hash)
	if err != nil {
		log.Fatal(err)
	}
}

//Найти запись
func FindHash(hash string) bool {
	db, err := sql.Open("postgres", "postgres://postgres:596run49@localhost/EFtest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var fi_id int64
	e := db.QueryRow(`select "file_id" from "FilesInfo" where "hash" = $1`, hash).Scan(&fi_id)
	if e != nil {
		fmt.Println(e)
	}
	if fi_id > 0 {
		return true
	}
	return false
}

func ListIndex() {
	db, err := sql.Open("postgres", "postgres://postgres:596run49@localhost/EFtest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select * from \"Users\"")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
