package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const contentSecurityPolicyValue = "default-src none; script-src 'self'; img-src 'self'; media-src 'self'; style-src-elem 'self'; style-src 'self';"
const contentSecurityPolicy = "Content-Security-Policy"
var lastIP string

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/resources/", resourceHandler)
	port := os.Getenv("PORT")
	if port == "" {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentSecurityPolicy, contentSecurityPolicyValue)
	ip := getIP(r)
	if lastIP != ip {
		fmt.Println(" ")
		lastIP = ip
	}
	if fileExists("public/" + r.URL.Path) {
		http.ServeFile(w, r, "public/" + r.URL.Path)
		fmt.Println(ip + strings.Repeat(" ", min(21 - len(ip), 1)) + " | " + r.URL.Path + strings.Repeat(" ", min(40 - len(r.URL.Path), 0)) + " | " + time.Now().Format("2006/01/02 03:04:05.000000 PM") + " | " + "200")
	} else {
		write404(w)
		fmt.Println(ip + strings.Repeat(" ", min(21 - len(ip), 1)) + " | " + r.URL.Path + strings.Repeat(" ", min(40 - len(r.URL.Path), 0)) + " | " + time.Now().Format("2006/01/02 03:04:05.000000 PM") + " | " + "404")
	}
}

/*func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentSecurityPolicy, contentSecurityPolicyValue)
	var (
		f   *os.File
		err error
	)
	if r.URL.Path != "/" {
		f, err = os.Open("public/" + r.URL.Path + ".html")
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
}*/

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		resourceHandler(w, r)
		return
	}
	w.Header().Set(contentSecurityPolicy, contentSecurityPolicyValue)
	out, err := interpretToHTML("public" + r.URL.Path)
	ip := getIP(r)
	if lastIP != ip {
		fmt.Println(" ")
		lastIP = ip
	}
	if err != io.EOF {
		write404(w)
		fmt.Println(ip + strings.Repeat(" ", min(21 - len(ip), 0)) + " | " + r.URL.Path + strings.Repeat(" ", min(40 - len(r.URL.Path), 0)) + " | " + time.Now().Format("2006/01/02 03:04:05.000000 PM") + " | " + "404")
	} else {
		fmt.Println(ip + strings.Repeat(" ", min(21 - len(ip), 0)) + " | " + r.URL.Path + strings.Repeat(" ", min(40 - len(r.URL.Path), 0)) + " | " + time.Now().Format("2006/01/02 03:04:05.000000 PM") + " | " + "200")
		fmt.Fprintf(w, out)
	}
}

func interpretToHTML(filename string) (output string, err error) {
	var f *os.File
	if filename == "public/" {
		f, err = os.Open("public/index.txt")
	} else if filename == "public/index" {
		err = os.ErrNotExist
	} else {
		f, err = os.Open(filename + ".txt")
	}
	defer f.Close()

	if err != nil {
		return
	}
	b := make([]byte, 1)
	var isCode int8 = 0
	isInCode := false
	var inCommand int8
	command := ""
	arg := ""
	inParagraph := false
	lineBreak := false

	for {
		var readTotal int
		readTotal, err = f.Read(b)
		if err != nil {
			if err != io.EOF {
				return
			}
			break
		}
		if str := string(b[:readTotal]); str == "`" {
			if isInCode {
				isCode--
			} else {
				isCode++
			}
			if isCode == 3 && !isInCode {
				if inParagraph {
					inParagraph = false
					output += "</p>"
				}
				isInCode = true
				output += "<pre>"
				lineBreak = false
			} else if isCode == 0 && isInCode {
				isInCode = false
				lastPre := strings.LastIndex(output, "<pre>") + 5
				output = output[:lastPre] + html.EscapeString(output[lastPre:])
				output += "</pre>"
				lineBreak = true
			}
		} else if inCommand == 2 && str == "\n" || (inCommand == 1 && command == "title" && str == "\n") {
			if inParagraph {
				output += "</p>"
				inParagraph = false
			} else if command == "title" {
				if arg != "" {
					output = "<!DOCTYPE html><html lang='en'><head><meta name='viewport' content='width=device-width, initial-scale=1'><title>" + arg + " | The Blob Blog</title><link rel='stylesheet' type='text/css' href='/resources/main.css'/><script defer src = '/resources/main.js'></script></head><body><header><nav><a href='/' id = 'home'><img src='favicon.ico' alt='A Blob'><span>Home</span></a><a href='/projects'>Projects</a><a href='/tutorials'>Coding Tutorials</a><a href='/music'>Music</a><form><input placeholder='Search...' aria-label='Search'></form><button aria-label='Close'><div id='bar1'></div><div id='bar2'></div><div id='bar3'></div><div id='bar4'></div></button></nav></header><main>"
				} else {
					output = "<!DOCTYPE html><html lang='en'><head><meta name='viewport' content='width=device-width, initial-scale=1'><title>The Blob Blog</title><link rel='stylesheet' type='text/css' href='/resources/main.css'/><script defer src = '/resources/main.js'></script></head><body><header><nav><a href='/' id = 'home'><img src='favicon.ico' alt='A Blob'><span>Home</span></a><a href='/projects'>Projects</a><a href='/tutorials'>Coding Tutorials</a><a href='/music'>Music</a><form><input placeholder='Search...' aria-label='Search'></form><button aria-label='Close'><div id='bar1'></div><div id='bar2'></div><div id='bar3'></div><div id='bar4'></div></button></nav></header><main>"
				}
			} else if command == "img" {
				output += "<img src = '" + arg + "'>"
			} else if command == "h1" {
				output += "<h1>" + arg + "</h1>"
			} else if command == "h2" {
				output += "<h2>" + arg + "</h2>"
			}
			inCommand = 0
			command = ""
			arg = ""
			lineBreak = true
		} else if inCommand == 2 {
			arg += str
		} else if inCommand == 1 && str == " " {
			inCommand = 2
			lineBreak = false
		} else if inCommand == 1 {
			command += str
			lineBreak = false
		} else if str == "#" {
			inCommand = 1
			lineBreak = false
		} else if str != "\n" || isInCode {
			output += str
			lineBreak = false
		} else if lineBreak && str == "\n" {
			if inParagraph {
					output += "</p><p>"
			} else {
				output += "<p>"
			}
			inParagraph = true
			lineBreak = false
		} else if str == "\n" {
			lineBreak = true
		}
	}
	output += "</p></main><footer><p>Copyright © 2020 <strong>Brian E. Powell</strong></p></footer></body></html>"
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

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func min(value int, minimum int) int {
	if value < minimum {
		return minimum
	}
	return value
}
