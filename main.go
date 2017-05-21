package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const (
	modelDir = "__models__"
)

type tModel struct {
	name    string // model's file name (will be appended to recovered file)
	header  []byte // initial data of the model file (all until last SOS marker 0xFFDA, excluding the marker)
	sosBloc []byte // SOS marker data bloc (used in case the corrupted file has no SOS marker)
}

var (
	gModels []tModel
	gRegExp = regexp.MustCompile(`\.id_(\d+?)_(.+?)\.onion\._`) // remove ransomware file extension
)

func main() {
	loadModels()

	// jpeg save quality
	var opt jpeg.Options
	opt.Quality = 98

	flag.Parse()
	for i := 0; i < flag.NArg(); i++ {
		filePath := flag.Arg(i)
		fmt.Printf("\n------------------------------------------------------------------------\nFile: %s\n", filepath.Base(filePath))

		// load data from the corrupted file
		fileData := loadFile(filePath)

		for _, model := range gModels {
			fmt.Printf(">> Model: %s\n", model.name)

			// merge model and file data
			img, err := model.appendFileData(fileData)
			if err != nil {
				fmt.Printf(">>>> Error: %v\n\n", err)
				continue
			}

			// create new file
			newFilePath := gRegExp.ReplaceAllString(filePath+model.name, "-")
			hf, err := os.Create(newFilePath)
			if err != nil {
				panic(err)
			}

			// encode new file
			if err := jpeg.Encode(hf, img, &opt); err != nil {
				fmt.Printf(">>>> Error: %v\n", err)
			} else {
				fmt.Printf(">>>> Done: %s\n", filepath.Base(newFilePath))
			}

			fmt.Printf("\n")
			hf.Close()
		}
	}

	fmt.Printf("------------------------------------------------------------------------\nPress enter to exit...")
	fmt.Scanln()
}

func (m *tModel) appendFileData(fileData []byte) (img image.Image, err error) {
	var data []byte
	// if file data already has a SOS marker, then juste merge model header and file data
	if fileData[0] == 0xFF && fileData[1] == 0xDA {
		data = append(m.header, fileData...)
	} else {
		// merge: model header / model SOS bloc / file data
		data = append(m.header, m.sosBloc...)
		data = append(data, fileData...)
	}
	return jpeg.Decode(bytes.NewReader(data)) // try to JPEG decode the new data
}

func loadFile(file string) (data []byte) {
	var err error
	data, err = ioutil.ReadFile(file)
	if err != nil {
		panic(err)
		return
	}
	//data = data[:len(data)-36] // remove the last 36 bytes appended by the ransomware

	// find the last SOS marker (there could be one for thumbnail)
	id := bytes.LastIndex(data, []byte{0xFF, 0xDA})
	if id == -1 || id < 10240 { // if no marker is found, or the marker is withing the encrypted 10kb
		var padding int
		fmt.Printf(">> No SOS (0xFFDA) marker found. first 10kb will be striped.\n>> How much padding do you want in the begining of the stream? ")
		fmt.Scanln(&padding)

		// extract all but the encrypted 10kb
		data = data[10240:]
		if padding > 0 {
			// prepend padding
			data = append(bytes.Repeat([]byte{0}, padding), data...)
		}
	} else {
		// extract SOS bloc
		data = data[id:]
	}
	return
}

func loadModels() {
	selfPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	filepath.Walk(filepath.Join(filepath.Dir(selfPath), modelDir), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		model, err := modelLoad(path)
		if err == nil {
			gModels = append(gModels, model)
		} else {
			fmt.Printf("Model load error: %v\n", err)
		}
		return nil
	})

	if len(gModels) <= 0 {
		panic(fmt.Sprintf("No models found!\nCreate a folder called \"__models__\" aside this executable and fill it with pictures taken with the same cameras as the encrypted ones, with different resolutions, orientation and quality settings. Rename the pictures because the model name will be appended to the recovered file (ex. sony-1080p-paysage.jpg"))
	}
}

func modelLoad(file string) (m tModel, err error) {
	var data []byte
	data, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	id := bytes.LastIndex(data, []byte{0xFF, 0xDA})
	if id == -1 {
		panic(fmt.Sprintf("Model %s is invalid JPEG (no SOS 0xFFDA marker)", filepath.Base(file)))
	}

	m.name = filepath.Base(file)
	m.header = data[:id]

	var sz uint16
	binary.Read(bytes.NewReader(data[id+2:id+4]), binary.BigEndian, &sz)
	m.sosBloc = data[id : id+2+int(sz)]

	return
}
