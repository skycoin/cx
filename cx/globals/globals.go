package globals

//Path is only used by os module and only to get working directory
//Path and working directory should not be hard coded into program struct (etc, when serialized)
//Working directory is property of executable and can be retrieved with golang library
var CxProgramPath string = ""