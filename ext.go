package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
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

// UseNotification ...
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

//FileExt Определить разрешение файла
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
	NewFileInfo(value, true, FileExt(pathLocalFile))
	return pathLocalFile
}

//HashFiles Определяем Hash файла
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

//NewFileInfo Создание новой записи
func NewFileInfo(value FileInfo, saveFile bool, ext string) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var TFId int
	_ = db.QueryRow(`select "TF_Id" from "TypeFile" where "TF_Extension" = $1`, ext).Scan(&TFId)
	if err != nil {
		log.Panic(err)
	}

	sqlStatement := `insert into "ArchFiles" 
					("AR_Name", "AR_Area", "AR_Filepath", "AR_Size", "AR_ModeTime", "AR_Hash", "AR_SaveFile", "AR_Ext", "AR_LocalPath")
					values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = db.Exec(sqlStatement, value.nameFile, value.area, value.filepath, value.size, value.modeTime, value.hash, saveFile, TFId, value.localFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

// NewFileInfoOS ...
func NewFileInfoOS(value *zip.File, fi FileInfo) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var TFId int
	_ = db.QueryRow(`select "TF_Id" from "TypeFile" where "TF_Extension" = $1`, fi.ext).Scan(&TFId)
	if err != nil {
		log.Panic(err)
	}

	sqlStatement := `insert into "ArchFiles"
	("AR_Name", "AR_Size", "AR_ModeTime", "AR_Hash", "AR_Ext", "AR_LocalPath")
	values ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(sqlStatement, fi.nameFile, fi.size, fi.modeTime, fi.hash, TFId, fi.localFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

//FindHash Найти запись
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

// ListIndex ...
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

//FindNotUnArch Поиск по ID
func FindNotUnArch() []*FileInfoSQL {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`select * from "ArchFiles" where "AR_Ext" = 1`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	listFiles := make([]*FileInfoSQL, 0)
	for rows.Next() {
		// var filesql *FileInfoSQL
		file := new(FileInfoSQL)
		err := rows.Scan(&file.ARId, &file.nameFile, &file.filepath, &file.size, &file.modeTime, &file.area, &file.hash, &file.saveFile,
			&file.ext, &file.localFilePath, &file.arparent)
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

// NewNullString ...
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// UnArchive ...
func UnArchive(srcPath string, dstPasth string) ([]FileInfo, error) {
	var filenames []FileInfo
	file, err := zip.OpenReader(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var fi FileInfo
	for _, f := range file.File {
		fpath := filepath.Join(dstPasth, f.Name)
		// f := os.FileInfo(fpath)
		if filepath.Ext(fpath) != ".sig" {

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

			h, err := f.Open()
			l := sha256.New()
			if _, err := io.Copy(l, h); err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("%x - %s\n", l.Sum(nil), f.Name)
			fi.hash = hex.EncodeToString(l.Sum(nil))
			fi.localFilePath = fpath
			fi.nameFile = f.Name
			fi.size = int64(f.UncompressedSize64)
			fi.ext = filepath.Ext(fpath)
			fi.modeTime = f.Modified
			h.Close()

			rc, err := f.Open()
			if err != nil {
				return filenames, err
			}
			defer rc.Close()

			_, err = io.Copy(outFile, rc)
			NewFileInfoOS(f, fi)
			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()
			filenames = append(filenames, fi)
			if err != nil {
				return filenames, err
			}
		}
	}
	return filenames, nil
}
