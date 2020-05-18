package archiver_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/nskeyedarchiver/archiver"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	testCases := map[string]struct {
		filename string
		expected string
	}{
		"test one value":       {"onevalue", "[true]"},
		"test all primitives":  {"primitives", "[1,1,1,1.5,\"YXNkZmFzZGZhZHNmYWRzZg==\",true,\"Hello, World!\",\"Hello, World!\",\"Hello, World!\",false,false,42]"},
		"test arrays and sets": {"arrays", "[[1,1,1,1.5,\"YXNkZmFzZGZhZHNmYWRzZg==\",true,\"Hello, World!\",\"Hello, World!\",\"Hello, World!\",false,false,42],[true,\"Hello, World!\",42],[true],[42,true,\"Hello, World!\"]]"},
		"test nested arrays":   {"nestedarrays", "[[true],[42,true,\"Hello, World!\"]]"},
	}

	for _, tc := range testCases {
		dat, err := ioutil.ReadFile("fixtures/" + tc.filename + ".xml")
		if err != nil {
			log.Fatal(err)
		}
		objects, err := archiver.Unarchive(dat)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, convertToJSON(objects))

		dat, err = ioutil.ReadFile("fixtures/" + tc.filename + ".bin")
		if err != nil {
			log.Fatal(err)
		}
		objects, err = archiver.Unarchive(dat)
		assert.Equal(t, tc.expected, convertToJSON(objects))
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
		_, err = archiver.Unarchive(dat)
		assert.Error(t, err)
	}
}

func convertToJSON(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}
