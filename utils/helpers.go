package utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// ParseXML is used to unmarshal an XML schema into a Go structure.
func ParseXML(path string, data interface{}) error {
	var err error

	source, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open the file: %s", path)
	}

	sourceData, err := ioutil.ReadAll(source)
	if err != nil {
		return fmt.Errorf("could not read the file: %s", path)
	}

	err = xml.Unmarshal(sourceData, data)
	if err != nil {
		return fmt.Errorf("could not unmarshal the XML: %s", path)
	}

	return nil
}
