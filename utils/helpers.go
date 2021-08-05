package utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// ParseXML helper for unmarshal XML into struct
func ParseXML(path string, data interface{}) error {
	var err error

	source, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %s", path)
	}

	sourceData, err := ioutil.ReadAll(source)
	if err != nil {
		return fmt.Errorf("could not read file: %s", path)
	}

	err = xml.Unmarshal(sourceData, data)
	if err != nil {
		return fmt.Errorf("could not unmarshal xml %s", path)
	}

	return nil
}
