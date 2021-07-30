package generator

import (
	"fmt"
	"github.com/b2broker/simplefix-go/utils"
	"os"
	"testing"
)

var generator *Generator

func TestMain(m *testing.M) {
	var err error
	doc := &Doc{}
	if err = utils.ParseXml("./testdata/fix.4.4.xml", doc); err != nil {
		panic(fmt.Errorf("could not make Doc XML: %s", err))
	}

	config := &Config{}
	if err = utils.ParseXml("./testdata/types.xml", config); err != nil {
		panic(fmt.Errorf("could not make Doc XML: %s", err))
	}

	generator = NewGenerator(doc, config, "fix")

	m.Run()
	os.Exit(0)
}
