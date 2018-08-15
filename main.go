package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jeremytorres/rawparser"
)

func main() {
	input := flag.String("i", "", "file or directory to convert")
	output := flag.String("o", "", "output location")
	flag.Parse()

	info, err := os.Stat(*input)
	if err != nil {
		log.Fatal(err)
	}

	if info.IsDir() {
		err = convertDirectory(*input, *output)
	} else {
		err = convertFile(*input, *output)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func convertFile(file, output string) error {
	if !strings.HasSuffix(strings.ToLower(file), ".nef") {
		return errors.New("input not a nef file")
	}
	parser, _ := rawparser.NewNefParser(true)
	info := &rawparser.RawFileInfo{
		File:    file,
		Quality: 100,
		DestDir: output,
	}
	_, err := parser.ProcessFile(info)
	return err
}

func convertDirectory(dir, output string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(strings.ToLower(f.Name()), ".nef") {
			continue
		}
		if err = convertFile(path.Join(dir, f.Name()), output); err != nil {
			return err
		}
	}
	return nil
}
