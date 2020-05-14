package archiver

func Run() {}

type NSKeyedObject struct {
}

func ArchiveXML([]NSKeyedObject) string {
	return ""
}
func ArchiveBin([]NSKeyedObject) []byte {
	return make([]byte, 0)
}

func UnarchiveXML(xml string) []NSKeyedObject {
	return nil
}
func UnarchiveBin(data []byte) []NSKeyedObject {
	return nil
}
