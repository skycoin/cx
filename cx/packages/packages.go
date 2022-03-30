package packages

// DefaultPackages ...
var DefaultPackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	//"al", "gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx", "json", "regexp", "cipher", "tcp",
	"al", "gl", "glfw", "time", "os", "gltext", "cx", "json", "cipher", "tcp",
}

// IsDefaultPackage ...
func IsDefaultPackage(ident string) bool {
	for _, core := range DefaultPackages {
		if core == ident {
			return true
		}
	}
	return false
}
