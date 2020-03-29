package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	files := groupFilesByMovement(filepath.Join(".", "character", "Dead"))
	numberOfMovements, numberOfImagesPerMovement, imageWidth, imageHeight := getFilesCount(files)

	r := image.Rectangle{Max: image.Point{X: numberOfImagesPerMovement * imageWidth, Y: numberOfMovements * imageHeight}}
	rgba := image.NewRGBA(r)

	i := 0
	for _, movement := range files {
		for imageKey, currentImage := range movement {
			rect := image.Rectangle{Min: image.Point{X: imageKey * imageWidth, Y: i * imageHeight}, Max: image.Point{X: imageKey*imageWidth + imageWidth, Y: i*imageHeight + imageHeight}}
			draw.Draw(rgba, rect, currentImage, image.Point{}, draw.Src)
		}
		i++
	}
	out, err := os.Create("./output.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(out, rgba)
}

func getFilesCount(files map[string][]image.Image) (int, int, int, int) {
	firstRowKey := getFirstRow(files)

	referenceFile := files[firstRowKey][0]

	return len(files), len(files[firstRowKey]), referenceFile.Bounds().Dx(), referenceFile.Bounds().Dy()

}

func groupFilesByMovement(dir string) map[string][]image.Image {
	files := scanDir(dir)
	m := make(map[string][]image.Image)

	for _, file := range files {
		i := strings.Index(file.Name(), " ")
		movement := file.Name()[0:i]

		if _, ok := m[movement]; ok {
			m[movement] = append(m[movement], openFile(file, dir))
		} else {
			m[movement] = []image.Image{openFile(file, dir)}
		}
	}

	return m
}
func sortName(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	// split numeric suffix
	i := len(name) - 1
	for ; i >= 0; i-- {
		if '0' > name[i] || name[i] > '9' {
			break
		}
	}
	i++
	// string numeric suffix to uint64 bytes
	// empty string is zero, so integers are plus one
	b64 := make([]byte, 64/8)
	s64 := name[i:]
	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}
	// prefix + numeric-suffix + ext
	return name[:i] + string(b64) + ext
}

func scanDir(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("before Sort")
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			fmt.Println(fi.Name())
		}
	}

	sort.Slice(
		files,
		func(i, j int) bool {
			println(sortName(files[i].Name()))
			println(sortName(files[j].Name()))

			return sortName(files[i].Name()) < sortName(files[j].Name())
		},
	)
	fmt.Println("After Sort")

	for _, fi := range files {
		if fi.Mode().IsRegular() {
			fmt.Println(fi.Name())
		}
	}

	return files
}

func openFile(file os.FileInfo, path string) image.Image {
	imageFile, err := os.Open(filepath.Join(path, file.Name()))
	if err != nil {
		log.Fatal(err)
	}

	img1, _, err := image.Decode(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	return img1
}

func getFirstRow(inputMap map[string][]image.Image) string {
	for key, _ := range inputMap {
		return key
	}

	return ""
}
