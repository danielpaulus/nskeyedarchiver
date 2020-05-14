package archiver_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/nskeyedarchiver/archiver"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	testCases := map[string]struct {
		filename string
		encode   []archiver.NSKeyedObject
	}{
		"simple boolean": {"boolean", []archiver.NSKeyedObject{}},
	}

	for _, tc := range testCases {
		dat, err := ioutil.ReadFile("fixtures/" + tc.filename + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, string(dat), archiver.ArchiveXML(tc.encode))

		dat, err = ioutil.ReadFile("fixtures/" + tc.filename + ".bin")
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, dat, archiver.ArchiveXML(tc.encode))
	}
}
