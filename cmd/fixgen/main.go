package main

import (
	"flag"
	"fmt"
	"github.com/b2broker/simplefix-go/generator"
	"github.com/b2broker/simplefix-go/utils"
	"os"
	"path/filepath"
)

func main() {
	var err error

	outputDir := flag.String("o", "./fix44/", "output directory")
	typesMappingPath := flag.String("t", "./source/types.xml", "path to XML file with types mapping")
	sourceXMLPath := flag.String("s", "./source/fix44.xml", "path to main XML file")

	flag.Parse()

	doc := &generator.Doc{}
	if err = utils.ParseXml(*sourceXMLPath, doc); err != nil {
		panic(fmt.Errorf("could not make Doc XML: %s", err))
	}

	config := &generator.Config{}
	if err = utils.ParseXml(*typesMappingPath, config); err != nil {
		panic(fmt.Errorf("could not make Doc XML: %s", err))
	}

	g := generator.NewGenerator(doc, config, filepath.Base(*outputDir))

	err = os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = g.Execute(*outputDir)
	if err != nil {
		panic(err)
	}
}
