package archiver

import (
	"bytes"
	"encoding/json"
	"fmt"

	plist "howett.net/plist"
)

const (
	archiverKey   = "$archiver"
	archiverValue = "NSKeyedArchiver"
	versionKey    = "$version"
	topKey        = "$top"
	objectsKey    = "$objects"
	versionValue  = 100000
)

type NSKeyedObject struct {
	isPrimitive bool
	primitive   interface{}
}

func ArchiveXML([]NSKeyedObject) string {
	return ""
}
func ArchiveBin([]NSKeyedObject) []byte {
	return make([]byte, 0)
}

func Unarchive(xml []byte) ([]NSKeyedObject, error) {
	plist, err := plistFromBytes(xml)
	if err != nil {
		return nil, err
	}
	nsKeyedArchiverData := plist.(map[string]interface{})

	err = verifyCorrectArchiver(nsKeyedArchiverData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func verifyCorrectArchiver(nsKeyedArchiverData map[string]interface{}) error {
	printAsJSON(nsKeyedArchiverData)
	if val, ok := nsKeyedArchiverData[archiverKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", archiverKey)
	} else {
		if stringValue := val.(string); stringValue != archiverValue {
			return fmt.Errorf("Invalid value: %s for key '%s', expected: '%s'", stringValue, archiverKey, archiverValue)
		}
	}
	if _, ok := nsKeyedArchiverData[topKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", topKey)
	}

	if _, ok := nsKeyedArchiverData[objectsKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", objectsKey)
	}

	if val, ok := nsKeyedArchiverData[versionKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", versionKey)
	} else {
		if stringValue := val.(uint64); stringValue != versionValue {
			return fmt.Errorf("Invalid value: %d for key '%s', expected: '%d'", stringValue, versionKey, versionValue)
		}
	}

	return nil
}

//ToPlist converts a given struct to a Plist using the
//github.com/DHowett/go-plist library. Make sure your struct is exported.
//It returns a string containing the plist.
func ToPlist(data interface{}) string {
	buf := &bytes.Buffer{}
	encoder := plist.NewEncoder(buf)
	encoder.Encode(data)
	return buf.String()
}
func plistFromBytes(plistBytes []byte) (interface{}, error) {
	var test interface{}
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))

	err := decoder.Decode(&test)
	if err != nil {
		return test, err
	}
	return test, nil
}

func printAsJSON(obj interface{}) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
