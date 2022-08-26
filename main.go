package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/fogleman/gg"
)

const url = "https://www.vzsar.ru/news/2022/08/24/v-saratove-chistyat-cvetyschyu-volgy-s-pomoschu-amfibii.html?utm_source=yxnews&utm_medium=desktop"
const picBaseUrl = "https://www.vzsar.ru"
const picDir = "pics"

func isJpg(s string) (result bool) {
	arr := strings.Split(s, ".")
	result = arr[len(arr)-1] == "jpg"
	return
}
//"../../../../Library/Fonts/Roboto-Regular.ttf"
func drawOnPicture(path, text string, x, y int) {
	const S = 1024
	im, err := gg.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(x, y)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./assets/fonts/Roboto-Regular.ttf", 20); err != nil {
		log.Fatal(err)
	}
	//dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)

	dc.DrawRoundedRectangle(0, 0, 512, 512, 0)
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored(text,  float64(x / 2), float64(y / 2), 0.5, 0.5)
	dc.Clip()
	dc.SavePNG(path)
}

func main() {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var arr []string

	doc.Find(`img`).Each(func(i int, s *goquery.Selection) {
		res, _ := s.Attr(`src`)
		if isJpg(res) {
			arr = append(arr, res)
		}
	})

	for _, imageUrl := range arr {
		res, err := http.Get(picBaseUrl + "/" + imageUrl)
		if err != nil {
			log.Fatalf("http.Get -> %v", err)
		}
		data, e := io.ReadAll(res.Body)
		if e != nil {
			log.Fatalf("http.Get -> %v", e)
		}
		res.Body.Close()
		if _, err := os.Stat(picDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(picDir, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		fileName := strings.Split(imageUrl, "/")
		file, err := os.Create(picDir + "/" + fileName[len(fileName)-1])
		if err != nil {
			log.Fatal(err)
		}
		file.Write(data)
		file, _ = os.Open(picDir + "/" + fileName[len(fileName)-1])
		img, _, _ := image.Decode(file)
		b := img.Bounds()
		label := "width:" + strconv.Itoa(b.Max.X) + " height:" + strconv.Itoa(b.Max.Y)
		fmt.Println(label)
		drawOnPicture(picDir + "/" + fileName[len(fileName)-1], label, b.Max.X, b.Max.Y)
	}
}

