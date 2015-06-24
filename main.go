package main

import (
	"fmt"
	"log"
	"github.com/disintegration/imaging"
	"image"
	"os/exec"
	"math"
	"gopkg.in/qml.v1"
	"os"
	"net/url"
	"sync"
	"image/color"
	"runtime"
	//"github.com/llgcode/draw2d"
	"path/filepath"
	"strings"
)

//go:generate genqrc assets

var wg sync.WaitGroup

var www chan int

func generatePng(){

	for {
		select {
			case width := <-www:
				log.Printf("Pinging: %d", width)
				expWidth := fmt.Sprintf("--export-width=%d", width )

				log.Printf("Getting wd")
				wd, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				log.Printf("Current wd: %s", wd)



				fileName := wd + filepath.FromSlash("/img/TTS-watermark-white.svg")
				log.Printf("FileName: %s", fileName)

				outFileName := "--export-png="+wd + filepath.FromSlash("/img/out.png")
				log.Println("Executing")
				out, err := exec.Command(
					INKSCAPE,
					"-z",
					fileName,
					outFileName,
					expWidth,
					//expHeight,
				).Output()
				log.Println("End of executing")

				if err != nil {
					log.Println(out)
					log.Fatal(err)
				}

				log.Println(out)
				wg.Done()
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	wg = sync.WaitGroup{}

	www = make(chan int)

	go generatePng()
	//wg.Add(1)
	//www <- 1600;
	//log.Println("Waiting:")
	//wg.Wait()
	//log.Println("Done.")

	if err := qml.Run(run); err != nil{
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func run() error {
	engine := qml.NewEngine()
	engine.AddImageProvider("pwd", func(id string, w, h int) image.Image{

		deadImage := imaging.New(800, 600, color.RGBA{0, 0, 0, 255})

		u, err := url.Parse(id)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ID from QT: %s", id)
		log.Printf("URL path: %v",u.Path)

		normalPath := filepath.FromSlash(u.Path)

		if runtime.GOOS == "windows" {
			normalPath = strings.TrimPrefix(normalPath, "\\")
		}

		log.Printf("FilePath: %s", normalPath)

		mydir, myfile := filepath.Split(normalPath)
		myext := filepath.Ext(myfile)

		// FIXME: DOWN!!!!!

		outDir := mydir + filepath.FromSlash("/out")

		if ok, _ := exists(outDir ); ok {

		} else {
			os.MkdirAll(outDir, os.ModeDir | os.ModePerm)
		}

		r := strings.NewReplacer(myext, "-wm" + myext)
		newfiles := outDir + r.Replace(myfile)

		log.Printf("NewFile %s", newfiles)


		srcImage, err := imaging.Open(normalPath) // FIXME: Error Under Windows
		if err != nil{
			log.Printf("Error open image file: %s", normalPath)
			log.Printf("Number of error: %v", err)
			return deadImage
		}

		waterImg, err := imaging.Open(filepath.FromSlash("./img/out.png"))
		if err != nil {
			return deadImage
		}

		width := srcImage.Bounds().Dx()
		height := srcImage.Bounds().Dy()

		origWIw := waterImg.Bounds().Dx()
		origWIh := waterImg.Bounds().Dy()

		// Это коэфициент, чтобы развиджить на 100%
		widthRatio := float64(width)/float64(origWIw)

		// С учетом 60% ширины
		destWidth := int(float64(origWIw) * widthRatio * 0.6)
		destHeight := int(float64(origWIh) * widthRatio  * 0.6)

		//log.Printf("Width: %d, wi Width: %d, Ratio: %f destW: %d", width, origWIw, widthRatio, destWidth )

		wImgFitted := imaging.Fit(waterImg, destWidth, destHeight, imaging.Lanczos)

		//log.Printf("DestWidth: %d", wImgFitted.Bounds().Dx())

		// Нужно расчитать куда выводить ватермарку
		wmWidth := wImgFitted.Bounds().Dx()
		wmHeight := wImgFitted.Bounds().Dy()

		wmBeginX := (width - wmWidth)/2
		wmBeginY := (height - wmHeight)/2

		wmBeginY2 := int(float32(height)* 0.2) - wmHeight/2
		wmBeginY3 := int(float32(height)* 0.8) - wmHeight/2

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

		var cleanImage *image.NRGBA

		cleanImage = imaging.New(width, height,color.RGBA{0, 0, 0, 0})

		waterMarked := imaging.Overlay(cleanImage, wImgFitted, image.Pt(wmBeginX, wmBeginY), 1.0)

		mMk2 := imaging.Overlay(waterMarked, wImgFitted, image.Pt(wmBeginX, wmBeginY2), 1.0)

		mMk3 := imaging.Overlay(mMk2, wImgFitted, image.Pt(wmBeginX, wmBeginY3), 1.0)

		// Делаем Blur
		bluredImg := imaging.Blur(mMk3, 5)

		destWidth = int(float64(origWIw) * widthRatio * 0.2)
		destHeight = int(float64(origWIh) * widthRatio  * 0.2)

		//log.Printf("Width: %d, wi Width: %d, Ratio: %f destW: %d", width, origWIw, widthRatio, destWidth )

		waterImg2 := imaging.Fit(waterImg, destWidth, destHeight, imaging.Lanczos)

		//log.Printf("DestWidth: %d", wImgFitted.Bounds().Dx())

		wmBeginX4 := int(float64(width)*0.95) - waterImg2.Bounds().Dx()
		wmBeginY4 := int(float64(height)*0.95) - waterImg2.Bounds().Dy()

		overlayedImg := imaging.Overlay(srcImage, bluredImg, image.Pt(0,0), 0.15) // FIXME: Нужно точно знать сколько нужно

		var im *image.NRGBA

		//log.Printf("X : %d", wmBeginX4)
		//log.Printf("Y : %d", wmBeginY4)

		if(meann > 50){
			log.Printf("Mean: %d", meann)
			waterImg2 = imaging.Invert(waterImg2)
			mMk4 := imaging.Overlay(overlayedImg, waterImg2, image.Pt(wmBeginX4, wmBeginY4), 1.0)
			//imaging.Save(mMk4, "test.jpg")
			im = mMk4
		} else {
			log.Printf("Mean: %d", meann)
			mMk4 := imaging.Overlay(overlayedImg, waterImg2, image.Pt(wmBeginX4, wmBeginY4), 1.0)
			//imaging.Save(mMk4, "test.jpg")
			im = mMk4
		}

		//im = bluredImg

		imaging.Save(im, newfiles)

		return im
	})

	component, err := engine.LoadFile("qrc:///assets/imgprovider.qml")
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()
	return nil
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
