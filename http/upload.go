package http

import (
	"github.com/uadmin/uadmin/blueprint/auth/services"
	"io"
	"net/http"
	"os"
	"strings"
	sessionsmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
)

// UploadImageHandler handles files sent from Tiny MCE's photo uploader
func UploadImageHandler(w http.ResponseWriter, r *http.Request, session *sessionsmodel.Session) {
	r.ParseMultipartForm(32 << 20)

	for _, f := range r.MultipartForm.File["file"] {
		src, _ := f.Open()
		folderPath := "./media/htmlimages/" + services.GenerateBase64(24) + "/"
		for {
			if _, err := os.Stat(folderPath); os.IsNotExist(err) {
				break
			}
			folderPath = "./media/htmlimages/" + services.GenerateBase64(24) + "/"
		}
		os.MkdirAll(folderPath, 0744)

		fileName := strings.Replace(f.Filename, "/", " ", -1)

		dst, _ := os.Create(folderPath + fileName)
		io.Copy(dst, src)
		src.Close()
		dst.Close()
		res := `{ "location" : "` + strings.TrimPrefix(folderPath+fileName, ".") + `" }`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res))
	}
}

