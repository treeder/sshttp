package system

import (
	"encoding/json"
	"log"
	"net/http"
	"os/user"
)

// Continue to refactor these out...

type Err struct {
	Code    int    `json:"status"`
	Message string `json:"error_message"`
}

func formatError(message string, code int) []byte {
	return doMarshall(Err{code, message})
}

func doMarshall(m interface{}) []byte {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println(err)
	}
	return b
}

// The rest can stay

func parseForm(res http.ResponseWriter, req *http.Request) (string, bool) {
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		return "", false
	}
	if req.FormValue("path") == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal( err )
		}
		log.Println( usr.HomeDir )
		return usr.HomeDir, true
	}
	return req.FormValue("path"), true
}

func DiskHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := parseForm(w, r)
	if !ok {
		return
	}
	w.Write(doMarshall(Disk(s)))
}

func FHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := parseForm(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case "GET":
		dir, err := IsDir(s)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if dir {
				w.Write(doMarshall(ListFiles(s)))
			} else {
				http.ServeFile(w, r, s)
			}
		}
	case "POST":
		log.Println("Posting file to ", s)
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		//if err == nil {
			err = WriteFile(s, handler.Filename, file)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				log.Println("POSTED!", s)
			}
		//} else {
			/*
			I don't think it should make a directory if the file param wasn't found, should be an explicit thing.
			_, err = MakeDir(s)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				log.Println("MADE DIR", s)
			}*/
	//	}
	case "DELETE":
		err := Remove(s)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(doMarshall(Host()))
}

func LoadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(doMarshall(Load()))
}

func RamHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(doMarshall(Ram()))
}

func CpuHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(doMarshall(Cpuinfo()))
}

func ProcessesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(doMarshall(Processes()))
}

func SystemHandler(w http.ResponseWriter, r *http.Request) {
	s, ok := parseForm(w, r)
	if !ok {
		return
	}
	w.Write(System(s))
}

func ShellHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var m interface{}
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}
		S := Shell(r.FormValue("exec"))
		m = ShellStruct{S.Output}

		w.Write(doMarshall(m))
	}
}
