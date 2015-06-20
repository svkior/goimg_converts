package main

import (
	"fmt"
	"log"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"github.com/disintegration/imaging"
	//"image/color"
	"image"
	"os/exec"
	"math"
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


func rgb2l(r,g,b uint32) (uint32){
	var l, rf, gf, bf, max, min float64

	rf = math.Max(math.Min(float64(r)/65535,1), 0)
	gf = math.Max(math.Min(float64(g)/65535,1), 0)
	bf = math.Max(math.Min(float64(b)/65535,1), 0)

	max = math.Max(rf, math.Max(gf,bf))
	min = math.Min(rf, math.Min(gf,bf))
	l = (max + min) / 2
	return uint32(l*100)
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
	if i > 1 {
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

	expWidth := fmt.Sprintf("--export-width=%d", int(float32(width) * 0.6) )
	//expHeight := fmt.Sprintf("--export-height=%d",height)

	out, err := exec.Command(
		"C:\\Program Files\\Inkscape\\inkscape.exe",
		"TTS-watermark-white.svg",
		"--export-png=out.png",
		expWidth,
		//expHeight,
	).Output()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stdout: %s", out)


	waterImg, err := imaging.Open("out.png")
	if err != nil {
		return
	}

	// Нужно расчитать куда выводить ватермарку
	wmWidth := waterImg.Bounds().Dx()
	wmHeight := waterImg.Bounds().Dy()

	wmBeginX := (width - wmWidth)/2
	wmBeginY := (height - wmHeight)/2

	wmBeginY2 := int(float32(height)* 0.2) - wmHeight/2
	wmBeginY3 := int(float32(height)* 0.8) - wmHeight/2


	// Делаем Blur
	bluredImg := imaging.Blur(waterImg, 10)

	// Определеяем какого цвета будем делать Watermark
	grImg := imaging.Grayscale(srcImage)

	Summ := uint32(0)
	Count := uint32(0)

	for idx :=0; idx < grImg.Bounds().Dx(); idx++ {
		for idy := 0; idy < grImg.Bounds().Dy(); idy++ {
			oldPixel := grImg.At(idx, idy)
			r, g, b, _ := oldPixel.RGBA()
			//log.Printf("%d %d %d %d", r, g, b, a)
			Summ += rgb2l(r,g,b)
			Count ++
		}
	}
	//log.Printf("Summ: %d, Count: %d, mean: %d", Summ, Count, Summ/Count)
	meann := Summ/Count



	waterMarked := imaging.Overlay(srcImage, bluredImg, image.Pt(wmBeginX, wmBeginY), 0.15)

	mMk2 := imaging.Overlay(waterMarked, bluredImg, image.Pt(wmBeginX, wmBeginY2), 0.15)

	mMk3 := imaging.Overlay(mMk2, bluredImg, image.Pt(wmBeginX, wmBeginY3), 0.15)

	expWidth2 := fmt.Sprintf("--export-width=%d", int(float32(width) * 0.2) )
	out, err = exec.Command(
		"C:\\Program Files\\Inkscape\\inkscape.exe",
		"TTS-watermark-white.svg",
		"--export-png=out2.png",
		expWidth2,
		//expHeight,
	).Output()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stdout: %s", out)
	waterImg, _ = imaging.Open("out2.png")


	wmBeginX4 := int(float64(width)*0.9) - waterImg.Bounds().Dx()
	wmBeginY4 := int(float64(height)*0.9) - waterImg.Bounds().Dy()


	if(meann > 50){
		log.Printf("Mean: %d", meann)
		waterImg = imaging.Invert(waterImg)
		mMk4 := imaging.Overlay(mMk3, waterImg, image.Pt(wmBeginX4, wmBeginY4), 1.0)
		imaging.Save(mMk4, "test.jpg")
	} else {
		log.Printf("Mean: %d", meann)
		mMk4 := imaging.Overlay(mMk3, waterImg, image.Pt(wmBeginX4, wmBeginY4), 1.0)
		imaging.Save(mMk4, "test.jpg")
	}



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
