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
		objects, err := archiver.UnarchiveXML(dat)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, objects)

		dat, err = ioutil.ReadFile("fixtures/" + tc.filename + ".bin")
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tc.expected, archiver.UnarchiveBin(dat))
	}
}

func TestValidation(t *testing.T) {

	testCases := map[string]struct {
		filename string
	}{
		"$archiver key is missing":         {"missing_archiver"},
		"$archiver is not nskeyedarchiver": {"wrong_archiver"},
		"$top key is missing":              {"missing_top"},
		"$objects key is missing":          {"missing_objects"},
		"$version key is missing":          {"missing_version"},
		"$version is wrong":                {"wrong_version"},
		"plist is invalid":                 {"broken_plist"},
	}

	for _, tc := range testCases {
		dat, err := ioutil.ReadFile("fixtures/" + tc.filename + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		_, err = archiver.UnarchiveXML(dat)
		assert.Error(t, err)
	}
}
