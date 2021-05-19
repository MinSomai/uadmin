package http

import (
	imageapi "github.com/uadmin/uadmin/blueprint/image/api"
	sessionmodel "github.com/uadmin/uadmin/blueprint/sessions/models"
	"github.com/uadmin/uadmin/model"
	"github.com/uadmin/uadmin/preloaded"
	"net/http"
)

func dAPIUpload(w http.ResponseWriter, r *http.Request, schema *model.ModelSchema) (map[string]string, error) {
	fileList := map[string]string{}

	if r.MultipartForm == nil {
		return fileList, nil
	}

	for k := range r.MultipartForm.File {
		// Process File
		// Check if the file is type file or image
		var field *model.F
		for i := range schema.Fields {
			if schema.Fields[i].ColumnName == k[1:] {
				field = &schema.Fields[i]
				r.MultipartForm.File[k[1:]] = r.MultipartForm.File[k]
				break
			}
		}
		if field == nil {
			// @todo, redo
			// utils.Trail(utils.WARNING, "dAPIUpload received a file that has no field: %s", k)
			continue
		}

		s := r.Context().Value(preloaded.CKey("session"))
		var session *sessionmodel.Session
		if s != nil {
			session = s.(*sessionmodel.Session)
		}

		fileName := imageapi.ProcessUpload(r, field, schema.ModelName, session, schema)
		if fileName != "" {
			fileList[field.ColumnName] = fileName
		}
	}
	return fileList, nil
}

