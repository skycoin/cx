package constants

// CorePackages ...
var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"al", "gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx", "json", "regexp", "cipher", "tcp",
}

// IsCorePackage ...
func IsCorePackage(ident string) bool {
	for _, core := range CorePackages {
		if core == ident {
			return true
		}
	}
	return false
}

