package main

import (
	"net/http"
	"log"
	"fmt"
	"io"
	"strings"
	"bytes"
	"time"
	"strconv"
	"io/ioutil"
    "html/template"
	"path"
	"os"
)

type UploadResult struct {
  SubmitTime    string
  Name 			string
  FileSize 		string
}


func main(){	
	// ex.NewTestStorage implements the "osin.Storage" interface

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/saveResult", saveResult)
	http.HandleFunc("/readlog", readlog)
	
	fmt.Println("init...") 
	log.Fatal(http.ListenAndServe(":9000", nil))
}
	
func upload (w http.ResponseWriter, r *http.Request) {
	
    if r.Method == "GET" {
    	
    	fp := path.Join("templates", "upload.html")
		
		tmpl, err := template.ParseFiles(fp)
	    
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		    return
	    }
		
		if err := tmpl.Execute(w, "no data needed"); err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		}
    } else {
	    //	start := time.Now()
		//r.ParseMultipartForm(32 << 20) // limit your max input length!
		
		r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512)
		
		var buf bytes.Buffer
		
		// in your case file would be fileupload
		file, header, err := r.FormFile("file")
		
		submitTime := r.FormValue("submitTime")		
		
		if err != nil {
		    panic(err)
		}
		 
		defer file.Close()
		
		name := strings.Split(header.Filename, ".")
		
		size := int64(header.Size)
		
		fmt.Println("File size... %s\n", size)
		
		fmt.Println("File name %s\n", name[0])
		
		// Copy the file data to my buffer
		io.Copy(&buf, file)
		
		buf.Reset()
		
		stInt, err := strconv.ParseInt(submitTime,10,64)
		
	   	fmt.Println("SubmitTime %s", time.Unix(0, stInt * int64(1000000)).Format("2006-01-02 03:04:05 PM"))

		uploadResult := UploadResult{submitTime, name[0], strconv.FormatInt(size, 10)}
		
	    fp := path.Join("templates", "upload_success.html")
	    
	    tmpl, err := template.ParseFiles(fp)
	    
	    if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		    return
	    }
		
		if err := tmpl.Execute(w, uploadResult); err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		}
    }
}

//Save REsult and save in temp file 
func saveResult (w http.ResponseWriter, r *http.Request) {
	
	f, err := os.OpenFile("log_time.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	
	defer f.Close()
	
	s:=
	"\n---------------------------------------\n"+ 
	"File name: " + r.FormValue("fileName") + 		
	"\n Submit Time: " + r.FormValue("subTime") + 		
	"\n File Size: " + r.FormValue("fileSize") + 		
	"\n Binomial took: " + r.FormValue("binomialTook") + " ms";	
	
	
	
	if _, err := f.WriteString(s); err != nil {
		log.Println(err)
	}

}


func readlog (w http.ResponseWriter, r *http.Request) {
	 content, err := ioutil.ReadFile("log_time.log")

     if err != nil {
          log.Fatal(err)
     }
	
	w.Write([]byte(content))

}
