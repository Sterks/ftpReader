package controller

import (
	"github.com/Sterks/ftpReader/ftp"
	"net/http"
)

// DownloadHandler ...
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	go ftp.MainProccess()
}
