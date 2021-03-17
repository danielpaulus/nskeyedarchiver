# archived in favor of https://github.com/danielpaulus/go-ios which has a fully working DTX implementation and can run UI tests on linux

# nskeyedarchiver
A Golang based implementation of Swift/ObjectiveCs NSKeyedArchiver/NSKeyedUnarchiver

Unarchive extracts NSKeyedArchiver Plists, either in XML or Binary format, and returns an array of the archived objects converted to usable Go Types.
- Primitives will be extracted just like regular Plist primitives (string, float64, int64, []uint8 etc.).
- NSArray, NSMutableArray, NSSet and NSMutableSet will transformed into []interface{}
- NSDictionary and NSMutableDictionary will be transformed into map[string] interface{}. I might add non string keys later.

Todos: 
- Add custom object support (anything that is not an array, set or dictionary)
- Add archiving/encoding support


Thanks howett.net/plist for your awesome Plist library :-) 
