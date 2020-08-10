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
)

func main(){	
	// ex.NewTestStorage implements the "osin.Storage" interface

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/readlog", readlog)
	
	fmt.Println("init...") 
	log.Fatal(http.ListenAndServe(":9000", nil))
}
	
func upload (w http.ResponseWriter, r *http.Request) {
	
	
    if r.Method == "GET" {
    	
    	w.Header().Set("diffResult", r.Header.Get("diffResult"))

        http.ServeFile(w, r, "upload/upload.html")
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
		
		subTimeObj := time.Unix(0, stInt * int64(1000000))

		elapsed := time.Since(subTimeObj)
	   
	   	fmt.Println("SubmitTime %s", time.Unix(0, stInt * int64(1000000)).Format("2006-01-02 03:04:05 PM"))

		fmt.Println("Binomial took %s", elapsed.Milliseconds())
		
		w.Header().Set("diffResult", submitTime)
		
		http.ServeFile(w, r, "upload/upload.html")
    }
}

func readlog (w http.ResponseWriter, r *http.Request) {
	 content, err := ioutil.ReadFile("out.log")

     if err != nil {
          log.Fatal(err)
     }
	
	w.Write([]byte(content))

}
