# CX Package Loader Format
This is how we’re representing the file. Every file has a length, a name, and a blake2 hash. Every package is a list of files – file structs. And then, we have a list of the package structs. And then, we can hash that to get the ID for the whole program.

## File Struct
- FileName string
- Length uint32
- Content []byte
- Blake2Hash UUID

## Package Struct
- PackageName string
- Files []UUID

## PackageList Struct
- Packages []UUID

### Compiler Frontend
1. Check all folders/directory, there should only be one package per folder/directory. Otherwise, give an error.
2. Create a `PackageList struct` for the program. 
3. Start with main package's source files. 
4. Create a `Package struct` for that package. Add its package name.
5. Extract the filename, length, content, and blake2 hash of that package's 
source files and create a `File struct` for it. 
6. Add the blake2 hash of the `File struct` to the `Files` of `Package struct`.
7. Generate blake2 hash of the `Package struct`.
8. Then add the blake2 hash of `Package struct` to the `Packages` of `PackageList struct` if it doesnt exist.
9. Check for imports of each source files of the package. Make a list of all the imports of that package. Sort the list and remove duplicate.
10. Then with the list, for every packages, do step 4 to 10 again. If path of package doesnt exist, give an error.
11. Use the PackageList as the input for the compiler.

### Package Loader Encoder
* Make a function to write the whole program, source files, etc to disk using the PackageList.
* Make a function to check and validate to verify all the hashes match.

### Objectives
This allows us to store the packages in our key-value store, and it allows us to store the program in our key-value store, so when we ask for a program, we just give the blake2 hash, and the web server responds with the program. This is how we’re going to package our games, how we’re going to package our CX apps, how we’re going to package our libraries so that they can be put on a key-value store in Redis initially, and eventually distributed peer-to-peer over DHT several years from now when we need to do that.