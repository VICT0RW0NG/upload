package main

import (
	"net/http"
	"log"
	"strings"
	"bytes"
	"fmt"
	"io"
)

func main(){	
	// ex.NewTestStorage implements the "osin.Storage" interface
	http.HandleFunc("/upload", upload)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
	
func upload (w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	var buf bytes.Buffer
	
	_ = "breakpoint"
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	
	if err != nil {
	    panic(err)
	}
	
	defer file.Close()
	
	name := strings.Split(header.Filename, ".")
	
	fmt.Printf("File name %s\n", name[0])
	
	// Copy the file data to my buffer
	io.Copy(&buf, file)
	
	contents := buf.String()

	fmt.Println(contents)

	buf.Reset()

}