import AppKit

let randomFilename = UUID().uuidString
let fullPath = URL(fileURLWithPath: "daniel")
// print(fullPath)
do {
    let archiver = NSKeyedArchiver(requiringSecureCoding: false)

    archiver.outputFormat = .xml
    archiver.encode("test")
    archiver.finishEncoding()
    print("data")
    print(String(bytes: archiver.encodedData, encoding: .ascii)!)
    let data = try NSKeyedArchiver.archivedData(withRootObject: "test", requiringSecureCoding: false)
    try data.write(to: fullPath)
} catch {
    print("Couldn't write file")
}

/*

 do {
     if let loadedStrings = try NSKeyedUnarchiver.unarchiveTopLevelObjectWithData(data) as? [String] {
         savedArray = loadedStrings
     }
 } catch {
     print("Couldn't read file.")
 }
 **/
