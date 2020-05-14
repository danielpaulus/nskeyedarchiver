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
		expected []interface{}
	}{
		"simple boolean": {"boolean", []interface{}{true}},
	}

	for _, tc := range testCases {
		dat, err := ioutil.ReadFile("fixtures/" + tc.filename + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tc.expected, archiver.UnarchiveXML(string(dat)))

		dat, err = ioutil.ReadFile("fixtures/" + tc.filename + ".bin")
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tc.expected, archiver.UnarchiveBin(dat))
	}
}
