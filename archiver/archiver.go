package archiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	plist "howett.net/plist"
)

const (
	archiverKey     = "$archiver"
	nsKeyedArchiver = "NSKeyedArchiver"
	versionKey      = "$version"
	topKey          = "$top"
	objectsKey      = "$objects"
	nsObjects       = "NS.objects"
	nsKeys          = "NS.keys"
	class           = "$class"
	className       = "$classname"
	versionValue    = 100000
)

const (
	nsArray        = "NSArray"
	nsMutableArray = "NSMutableArray"
	nsSet          = "NSSet"
	nsMutableSet   = "NSMutableSet"
)

const (
	nsDictionary        = "NSDictionary"
	nsMutableDictionary = "NSMutableDictionary"
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

func Unarchive(xml []byte) ([]interface{}, error) {
	plist, err := plistFromBytes(xml)
	if err != nil {
		return nil, err
	}
	nsKeyedArchiverData := plist.(map[string]interface{})

	err = verifyCorrectArchiver(nsKeyedArchiverData)
	if err != nil {
		return nil, err
	}
	return extractObjectsFromTop(nsKeyedArchiverData[topKey].(map[string]interface{}), nsKeyedArchiverData[objectsKey].([]interface{}))

}

func extractObjectsFromTop(top map[string]interface{}, objects []interface{}) ([]interface{}, error) {
	objectCount := len(top)
	objectRefs := make([]plist.UID, objectCount)
	for i := 0; i < objectCount; i++ {
		objectIndex := top[fmt.Sprintf("$%d", i)].(plist.UID)
		objectRefs[i] = objectIndex
	}
	return extractObjects(objectRefs, objects)
}

func extractObjects(objectRefs []plist.UID, objects []interface{}) ([]interface{}, error) {
	objectCount := len(objectRefs)
	returnValue := make([]interface{}, objectCount)
	for i := 0; i < objectCount; i++ {
		objectIndex := objectRefs[i]
		objectRef := objects[objectIndex]
		if object, ok := isPrimitiveObject(objectRef); ok {
			returnValue[i] = object
			continue
		}
		if object, ok := isArrayObject(objectRef.(map[string]interface{}), objects); ok {
			extractObjects, err := extractObjects(toUidList(object[nsObjects].([]interface{})), objects)
			if err != nil {
				return nil, err
			}
			returnValue[i] = extractObjects
			continue
		}

		if object, ok := isDictionaryObject(objectRef.(map[string]interface{}), objects); ok {
			dictionary, err := extractDictionary(object, objects)
			if err != nil {
				return nil, err
			}
			returnValue[i] = dictionary
			continue
		}

		objectType := reflect.TypeOf(objectRef).String()
		return nil, fmt.Errorf("Unknown object type:%s", objectType)

	}
	return returnValue, nil
}

func extractDictionary(object map[string]interface{}, objects []interface{}) (map[string]interface{}, error) {
	keyRefs := toUidList(object[nsKeys].([]interface{}))
	keys, err := extractObjects(keyRefs, objects)
	if err != nil {
		return nil, err
	}

	valueRefs := toUidList(object[nsObjects].([]interface{}))
	values, err := extractObjects(valueRefs, objects)
	if err != nil {
		return nil, err
	}
	mapSize := len(keys)
	result := make(map[string]interface{}, mapSize)
	for i := 0; i < mapSize; i++ {
		result[keys[i].(string)] = values[i]
	}

	return result, nil
}

func toUidList(list []interface{}) []plist.UID {
	l := len(list)
	result := make([]plist.UID, l)
	for i := 0; i < l; i++ {
		result[i] = list[i].(plist.UID)
	}
	return result
}

func isDictionaryObject(object map[string]interface{}, objects []interface{}) (map[string]interface{}, bool) {
	className, err := resolveClass(object[class], objects)
	if err != nil {
		return nil, false
	}
	if className == nsDictionary || className == nsMutableDictionary {
		return object, true
	}
	return object, false
}

func isArrayObject(object map[string]interface{}, objects []interface{}) (map[string]interface{}, bool) {
	className, err := resolveClass(object[class], objects)
	if err != nil {
		return nil, false
	}
	if className == nsArray || className == nsMutableArray || className == nsSet || className == nsMutableSet {
		return object, true
	}
	return object, false
}

func resolveClass(classInfo interface{}, objects []interface{}) (string, error) {
	printAsJSON(reflect.TypeOf(classInfo).String())
	if v, ok := classInfo.(plist.UID); ok {
		classDict := objects[v].(map[string]interface{})
		return classDict[className].(string), nil
	}
	return "", fmt.Errorf("Could not find class for %s", classInfo)
}

func isPrimitiveObject(object interface{}) (interface{}, bool) {
	if v, ok := object.(uint64); ok {
		return v, ok
	}
	if v, ok := object.(float64); ok {
		return v, ok
	}
	if v, ok := object.(bool); ok {
		return v, ok
	}
	if v, ok := object.(string); ok {
		return v, ok
	}
	if v, ok := object.([]uint8); ok {
		return v, ok
	}
	return object, false
}

func verifyCorrectArchiver(nsKeyedArchiverData map[string]interface{}) error {
	if val, ok := nsKeyedArchiverData[archiverKey]; !ok {
		return fmt.Errorf("Invalid NSKeyedAchiverObject, missing key '%s'", archiverKey)
	} else {
		if stringValue := val.(string); stringValue != nsKeyedArchiver {
			return fmt.Errorf("Invalid value: %s for key '%s', expected: '%s'", stringValue, archiverKey, nsKeyedArchiver)
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
