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
func NewFileInfo(nameFile string, area string, filepath string, size int64, modeTime time.Time, hash string) {
	db, err := sql.Open("postgres", "postgres://postgres:596run49@localhost/EFtest?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `insert into "FilesInfo" ("namefile", "area", "filepath", "size", "modetime", "hash") values ($1, $2, $3, $4, $5, $6)`
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
	var fiID int64
	_ = db.QueryRow(`select "file_id" from "FilesInfo" where "hash" = $1`, hash).Scan(&fiID)

	if fiID > 0 {
		return true
	} else {
		return false
	}
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
