# CX Package Format
This is how we’re representing the file. Every file has a length, a name, and a SHA256 hash. Every package is a list of files – file structs. And then, we have a list of the package structs. And then, we can serialize it and hash that to get the ID for the whole program.

## File Struct
- Name string
- Length uint32
- SHA256Hash UUID

## Package Struct
- Name string
- Files []FileStruct

## CXProgram Struct
- Packages []PackageID

1. We start with making package structs.
2. Each package struct will have a name and a list of file struct. 
3. Each file struct will have a filename, a UINT32 length, and a SHA256 hash of the file content.
4. The file content will be stored in a key-value store with its SHA256 hash as its key. 
5. We will sort the files in the order of the SHA256 hash. 
6. And then, we will serialize, using [SkyEncoder](https://github.com/skycoin/skyencoder), the Package struct and then we hash the serialization with SHA256, so we get a Package ID.
7. The Package struct will be stored in a key-value store with its SHA256 hash(Package ID) as its key. 
8. The CXProgram struct will then be composed of a list of packageID in which the first packageID is the main package.
9. And then we will serialize, using [SkyEncoder](https://github.com/skycoin/skyencoder), the CXProgram struct and then we hash the serialization with SHA256, so we get a CXProgram ID.
10. The CXProgram struct will be stored in a key-value store with its SHA256 hash(CXProgram ID) as its key. 

Then for retrieving the packages and files to run the cx program
1. Using only the CXProgram ID, retrieve the serialized CXProgram struct and deserialize it.
2. Using the deserialized CXProgram struct, retrieve the serialized packages and deserialize it.
3. Retrieve the package CX files using their SHA256 hash.
4. Run the CX Program.

This allows us to store the packages in our key-value store, and it allows us to store the program in our key-value store, so when we ask for a program, we just give the SHA256 hash, and the web server responds with the program. This is how we’re going to package our games, how we’re going to package our CX apps, how we’re going to package our libraries so that they can be put on a key-value store in Redis initially, and eventually distributed peer-to-peer over DHT several years from now when we need to do that.

