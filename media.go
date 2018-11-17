package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//FileUploadHandler handles uploading naming and resizing images
func FileUploadHandler(s Server) func(w http.ResponseWriter, r *http.Request) {

	// stmt, err := s.DB.Prepare(`SELECT id, created_at, name, email FROM users`)
	// if err != nil {
	// 	log.Fatal("Error preparing sql statement: " + err.Error())
	// }

	return func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//TODO make sure we can't do a full path or ../ or something sneaky like that
		out := filepath.Join("media", fileHeader.Filename)

		f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.Copy(f, file)
		f.Close()
		file.Close()

		w.Write([]byte(out))
	}
}
