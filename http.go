package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	var defaultPath string = filepath.Dir(ex)
	var defaultPort int = 3000

	dir := flag.String("dir", defaultPath, "http Folder")
	port := flag.Int("port", defaultPort, "http port number")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Printf("Please input which directory what you want to share,\ndefault: \"" + defaultPath + "\":\n")
		n, err := fmt.Scanf("%s\n", dir)
		if err != nil {
			if n == 0 {
				dir = &defaultPath
			} else {
				fmt.Println(n, err)
			}
			err = nil
		} else {
			checkDir(defaultPath, *dir)
		}

		fmt.Printf("Please input port Number\ndefault: " + strconv.Itoa(*port) + ": \n")
		m, err := fmt.Scanf("%d\n", port)
		if err != nil {
			if m == 0 {
				port = &defaultPort
			} else {
				fmt.Println(m, err)
			}
			err = nil
		}
	}

	if *dir == defaultPath {
		log.Println("http folder:", *dir)
	} else {
		log.Println("http folder:", filepath.Join(defaultPath, *dir))
		checkDir(defaultPath, *dir)
	}
	log.Println("listening...", "localhost:"+strconv.Itoa(*port))

	panic(http.ListenAndServe(":"+strconv.Itoa(*port), &logServer{
		hdl: http.FileServer(http.Dir(*dir)),
	}))
}

type logServer struct {
	hdl http.Handler
}

func (l *logServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print(r.RemoteAddr, r.URL.Path)
	l.hdl.ServeHTTP(w, r)
}

func checkDir(defaultPath, dir string) {
	if _, err := os.Stat(filepath.Join(defaultPath, dir)); err != nil {
		if os.IsNotExist(err) {
			log.Println("creates a directory named ", dir)
			err := os.MkdirAll(filepath.Join(defaultPath, dir), 0777)
			if err != nil {
				log.Fatal("MkdirAll : %s", err)
			}
		} else {
			log.Fatal("other error")
			// other error
		}
	}
}
