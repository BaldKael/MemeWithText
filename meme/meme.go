/**
 * Author: BaldKael
 * Email:  baldkael@sina.com
 * Date:   2020/8/15 15:53
 */

package meme

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"github.com/nfnt/resize"
)

const (
	width  = 300
	height = 300
)

type Meme struct {
	SourcePath   string
	TargetPath   string
	TextContent  string
	TextColor    []int // RGB
	TextFont     string
	PicSplit     int // if pic split to 5 * 5, the text position should in [0. 25)
	TextPosition int // [0, 25]
	TextSize     int
	Img          image.Image
	tmpPath      string
}

/*
	* * * * *
	* * * * *
	* * * * *
	* * √ * *
	* * * * *

	a pic will split into 25 block
	if position == x, the text will write on the x block
	default position is 17, which is the √ position
*/

func New(sourcePath, targetPath, textFont, textContent string, textPosition, textSize, picSplit int, textColor []int) *Meme {
	// TODO: if http in sourcePath, will download this pic
	meme := new(Meme)

	if textFont == "" {
		meme.TextFont = "default"
	} else {
		meme.TextFont = textFont
	}
	if textColor == nil {
		meme.TextColor = []int{0, 0, 0}
	} else {
		meme.TextColor = textColor
	}
	meme.TextPosition = textPosition
	meme.PicSplit = picSplit
	if textPosition < 0 || picSplit < 1 || picSplit > 10 {
		meme.TextPosition = 7
		meme.PicSplit = 5
	} else {
		if textPosition >= picSplit*picSplit || textPosition < 0 {
			meme.TextPosition = 7
			meme.PicSplit = 5
		}
	}
	if textPosition < 0 && picSplit > 1 && picSplit <= 10 {
		// calculate the recommend position
		x := (picSplit + 1)/2 - 1
		//y := (picSplit + 1) / 2
		//if picSplit > 6 {
		//	y++
		//}
		y := picSplit - 1

		meme.TextPosition = y*picSplit + (x + 1)
		//log.Println(picSplit, "recommend:", x, y, meme.TextPosition)
		meme.PicSplit = picSplit
	}
	if textSize < 0 || textSize > 30 {
		meme.TextSize = 15
	} else {
		meme.TextSize = textSize
	}
	meme.SourcePath = sourcePath
	meme.TargetPath = targetPath
	meme.TextContent = textContent

	return meme
}

// resize the pic
func (meme *Meme) Resize() *Meme {
	if meme.Img == nil {
		file, err := os.Open(meme.SourcePath)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = file.Close()
		}()

		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}

		meme.Img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	} else {
		meme.Img = resize.Resize(uint(width), uint(height), meme.Img, resize.Lanczos3)
	}
	log.Println("resize done")
	return meme
}

// add text for the pic
func (meme *Meme) AddText() *Meme {
	meme.generateTextPic()
	meme.combinePic()
	log.Println("add text done")
	return meme
}

func (meme *Meme) generateTextPic() {
	prefixNum := rand.Intn(10000)
	textPicFileName := fmt.Sprintf("tmp/font_%s.png", strconv.Itoa(prefixNum))
	meme.tmpPath = textPicFileName

	imgFile, err := os.Create(textPicFileName)
	defer func() {
		_ = imgFile.Close()
	}()
	if err != nil {
		panic(err)
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// RGBA(A is transparency)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 0})
		}
	}

	// read font ttf
	fontFile := "source/fonts/WenQuanYiZenHei.ttf"
	if meme.TextFont != "default" {
		fontFile = "source/fonts/" + meme.TextFont + ".ttf"
	}
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		panic(err)
	}

	// parse font
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	f := freetype.NewContext()
	f.SetFont(font)
	f.SetFontSize(float64(meme.TextSize))
	f.SetClip(img.Bounds())
	f.SetDst(img)

	var colorCheck = true
	if len(meme.TextColor) != 3 {
		colorCheck = false
	}
	for _, v := range meme.TextColor {
		if v > 255 || v < 0 {
			colorCheck = false
			break
		}
	}
	if !colorCheck {
		f.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255}))
	} else {
		f.SetSrc(image.NewUniform(color.RGBA{uint8(meme.TextColor[0]), uint8(meme.TextColor[1]), uint8(meme.TextColor[2]), 255}))
	}

	y := meme.TextPosition / meme.PicSplit
	x := (meme.TextPosition - y * meme.PicSplit) - 1

	//log.Println(meme.TextPosition, meme.PicSplit, x, y)
	//log.Println(meme.PicSplit, 300 * x / meme.PicSplit, 300 * y / meme.PicSplit)

	position := freetype.Pt(300 * (x + 1) / meme.PicSplit - len(meme.TextContent) / 8 * meme.TextSize , 300 * (y + 1) / meme.PicSplit - meme.TextSize)

	_, err = f.DrawString(meme.TextContent, position)
	if err != nil {
		panic(err)
	}

	err = png.Encode(imgFile, img)
	if err != nil {
		panic(err)
	}

	log.Println("generate font pic done")
}

func (meme *Meme) combinePic() {
	log.Println(meme.tmpPath)
	pic2, err := os.Open(meme.tmpPath)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	img2, err := png.Decode(pic2)
	if err != nil {
		log.Println("combine err:", err)
		//panic(err)
	}

	b := meme.Img.Bounds()
	m := image.NewRGBA(b)

	draw.Draw(m, b, meme.Img, image.ZP, draw.Src)
	draw.Draw(m, img2.Bounds(), img2, image.ZP, draw.Over)

	meme.Img = m
	log.Println("combine done")
}

func (meme *Meme) Save() *Meme {
	if img, err := os.Create(meme.TargetPath); err != nil {
		panic(err)
	} else {
		defer func() {
			_ = img.Close()
		}()

		if err := png.Encode(img, meme.Img); err != nil {
			panic(err)
		}
	}
	return meme
}

func (meme *Meme) Clean() *Meme {
	if err := os.Remove(meme.tmpPath); err != nil {
		log.Println(err)
	}
	return meme
}
