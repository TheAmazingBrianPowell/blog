package main

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/resources/", resourceHandler)
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src none; script-src 'self'; img-src 'self'; media-src 'self'; style-src-elem 'self';")
	if fileExists("public/" + r.URL.Path) {
		http.ServeFile(w, r, "public/"+r.URL.Path)
	} else {
		write404(w)
	}
	/*var (
		f *os.File
		err error
	)
	f, err = os.Open("public/" + r.URL.Path)
	defer f.Close()
	if err != nil {
		write404(w)
	} else {
		w.Header().Set("Content-Type", getFileContentType(r.URL.Path))
		b := make([]byte, 4)
		var out string
		for {
			readTotal, err := f.Read(b)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			out += string(b[:readTotal])
		}
		fmt.Fprintf(w, out)
	}*/
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src none; script-src 'self'; img-src 'self'; media-src 'self'; style-src-elem 'self';")
	var (
		f   *os.File
		err error
	)
	if r.URL.Path != "/" {
		//if nameLoc := strings.Index(r.URL.Path, "."); nameLoc != -1 {
		//	f, err = os.Open(r.URL.Path[:nameLoc] + ".html")
		//} else {
		f, err = os.Open("public/" + r.URL.Path + ".html")
		//}
	} else {
		f, err = os.Open("public/index.html")
	}
	defer f.Close()
	if err != nil {
		write404(w)
	} else {
		b := make([]byte, 4)
		var out string
		for {
			readTotal, err := f.Read(b)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			out += string(b[:readTotal])
		}
		fmt.Fprintf(w, out)
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src none; script-src 'self'; img-src 'self'; media-src 'self'; style-src-elem 'self';")
	out, err := interpretToHTML("public" + r.URL.Path)
	if err != nil {
		write404(w)
	}
	fmt.Fprintf(w, out)
}

func interpretToHTML(filename string) (output string, err error) {
	var f *os.File

	f, err = os.Open(filename + ".txt")

	if err != nil {
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)
	var title string
	isCode := 0
	isInCode := false
	started := false
	var lastChar string
	firstLine := true
	for scanner.Scan() {
		if str := scanner.Text(); str == "`" {
			if isInCode {
				isCode--
			} else {
				isCode++
			}
			if isCode == 3 && !isInCode {
				isInCode = true
				output += "<pre>"
			} else if isCode == 0 && isInCode {
				isInCode = false
				lastPre := strings.LastIndex(output, "<pre>") + 5
				output = output[:lastPre] + html.EscapeString(output[lastPre:])
				output += "</pre>"
			}
			lastChar = "`"
		} else if str == "\n" && !isInCode && lastChar != "\n" {
			if firstLine {
				output = "<!DOCTYPE html><html lang='en'><head><meta name='viewport' content='width=device-width, initial-scale=1'><title>" + title + " | The Blob Blog</title><link rel='stylesheet' type='text/css' href='/resources/main.css'></head><body><main><h1>" + title + "</h1>"
				firstLine = false
			}
			if started {
				output += "</p><p>"
				started = true
			} else {
				output += "<p>"
			}
			lastChar = "\n"
		} else if firstLine {
			title += str
		} else if str != "\n" || isInCode {
			output += str
			lastChar = str
		}
	}

	output += "</p></main></body></html>"

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func write404(w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Error 404</title><link rel='stylesheet' type='text/css' href='/resources/main.css'></head><body><h1>Error 404: Page Not Found</h1></body></html>")
}
