package system

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type DirInfo struct {
	Count    int        `json:"count"`
	Entities []FileInfo `json:"entities"`
}

type FileInfo struct {
	Filename string `json:"name"`
	Size     int64  `json:"size"`
	UpdateAt string `json:"updated_at"`
	IsDir    bool   `json:"is_dir"`
}

func getFileInfo(f os.FileInfo) (fl FileInfo) {
	return FileInfo{Filename: f.Name(), Size: f.Size(), UpdateAt: f.ModTime().Format("Jan 2, 2006 at 3:04pm (MST)"), IsDir: f.IsDir()}
}

func IsDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

func ListFiles(path string) (d DirInfo) {

	files, _ := ioutil.ReadDir(path)
	counter := 0
	for _, f := range files {
		counter++
		fl := getFileInfo(f)
		d.Entities = append(d.Entities, fl)
	}
	d.Count = counter
	return
}

func MakeDir(filePath string) (string, error) {
	filePath = path.Clean(filePath)
	err := os.MkdirAll(filePath, 0777)
	return filePath, err
}

func Remove(filePath string) error {
	filePath = path.Clean(filePath)
	err := os.RemoveAll(filePath)
	return err
}

// todo: this should stream write
func WriteFile(filePath string, filename string, file io.Reader) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return err
	}
	filePath, err = MakeDir(filePath)

	if err == nil {
		err = ioutil.WriteFile(path.Join(filePath, filename), data, 0777)
	}
	return err
}
