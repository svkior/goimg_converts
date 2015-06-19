package main

import (
	"fmt"
	"log"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"github.com/disintegration/imaging"
	//"image/color"
	"image"
	"os/exec"
)

func main() {
	mw := &MyMainWindow{
		model: NewImgModel(),
	}

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Walk ListBox Example",
		MinSize:  Size{240, 320},
		Size:     Size{400, 600},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					ListBox{
						AssignTo: &mw.lb,
						Model:    mw.model,
						OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
						OnItemActivated:       mw.lb_ItemActivated,
					},
					HSplitter{
						Children: []Widget{
							ImageView{
								AssignTo: &mw.iw,
							},
						},

					},

				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

type MyMainWindow struct {
	*walk.MainWindow
	model *ImgModel
	lb    *walk.ListBox
	iw    *walk.ImageView
	iw2   *walk.ImageView
}

func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	i := mw.lb.CurrentIndex()
	if i > 0 {
		return
	}

	if i < 0 {
		return
	}
	item := &mw.model.items[i]

//	mw.te.SetText(item.path)


	srcImage, err := imaging.Open(item.path)
	if err != nil {
		return
	}

	width := srcImage.Bounds().Dx()
	height := srcImage.Bounds().Dy()

	expWidth := fmt.Sprintf("--export-width=%d", width)
	expHeight := fmt.Sprintf("--export-height=%d",height)

	out, err := exec.Command(
		"C:\\Program Files\\Inkscape\\inkscape.exe",
		"TTS-watermark-white.svg",
		"--export-png=out.png",
		expWidth,
		expHeight,
	).Output()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stdout: %s", out)


	waterImg, err := imaging.Open("out.png")
	if err != nil {
		return
	}

	// Делаем Blur
	bluredImg := imaging.Blur(waterImg, 10)

	waterMarked := imaging.Overlay(srcImage, bluredImg, image.Pt(0, 0), 0.20)

	imaging.Save(waterMarked, "test.jpg")

	img, err := walk.NewImageFromFile("test.jpg")
	if err != nil {
		return
	}

	mw.iw.SetImage(img)

	fmt.Println("CurrentIndex: ", i)
	fmt.Println("CurrentEnvVarName: ", item.name)
}

func (mw *MyMainWindow) lb_ItemActivated() {
	value := mw.model.items[mw.lb.CurrentIndex()].path

	walk.MsgBox(mw, "Path", value, walk.MsgBoxIconInformation)
}

/* Модель списка изображений */

type ImgItem struct {
	name string
	path string
}

type ImgModel struct {
	walk.ListModelBase
	items []ImgItem
}

func NewImgModel() *ImgModel {
	files, _ := ioutil.ReadDir(".\\from")

	m := &ImgModel{items: make([]ImgItem, len(files))}

	for i,f := range files {
		fmt.Println(f.Name())

		name := f.Name()
		path := ".\\from\\" + name
		m.items[i] = ImgItem{name, path}
	}
	return m
}

func (m *ImgModel) ItemCount() int {
	return len(m.items)
}

func (m *ImgModel) Value(index int) interface{} {
	return m.items[index].name
}
