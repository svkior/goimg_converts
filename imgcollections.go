package main
import (
	"gopkg.in/qml.v1"
	"log"
	"net/url"
	"runtime"
	"path/filepath"
	"strings"
)

type Images struct {
	list []string
	Len int
}

func (images *Images) Add(fileName string){
	images.list = append(images.list, fileName)
	images.Len = len(images.list)
	qml.Changed(images, &images.Len)
}

func (images *Images) Image(index int) string {
	return images.list[index]
}

func (images *Images) ImageName(index int) string {

	_, myfile := filepath.Split(images.list[index])

	extension := filepath.Ext(myfile)
	name := myfile[0:len(myfile)-len(extension)]

	return name
}

func (images *Images) Clear(){
	images.list = []string{}
	images.Len = len(images.list)
	qml.Changed(images, &images.Len)
}


func (images *Images) AddDir(dir string, ext string) {
	files, _ := filepath.Glob(filepath.Join(dir, ext))
	for _, f := range files {
		images.Add(f)
	}
}


func (images *Images) Scan(dir string) {
	log.Printf("DIR: %s", dir)
	u, err := url.Parse(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("URL path: %v",u.Path)
	normalPath := filepath.FromSlash(u.Path)

	if runtime.GOOS == "windows" {
		normalPath = strings.TrimPrefix(normalPath, "\\")
	}

	log.Printf("Normal path: %v", normalPath)

	images.AddDir(normalPath, "*.jp*")
	images.AddDir(normalPath, "*.JP*")
	images.AddDir(normalPath, "*.Jp*")
	images.AddDir(normalPath, "*.pn*")
	images.AddDir(normalPath, "*.Pn*")
	images.AddDir(normalPath, "*.PN*")
}