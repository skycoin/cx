package base

type byFnName []*CXFunction
type byTypName []*CXType
type byModName []*CXModule
type byDefName []*CXDefinition
type byStrctName []*CXStruct
type byFldName []*CXField
type byParamName []*CXParameter

/*
  Lens
*/

func (s byFnName) Len() int {
    return len(s)
}
func (s byTypName) Len() int {
    return len(s)
}
func (s byModName) Len() int {
    return len(s)
}
func (s byDefName) Len() int {
    return len(s)
}
func (s byStrctName) Len() int {
    return len(s)
}
func (s byFldName) Len() int {
    return len(s)
}
func (s byParamName) Len() int {
    return len(s)
}

/*
  Swaps
*/

func (s byFnName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byTypName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byModName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byDefName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byStrctName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byFldName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byParamName) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

/*
  Lesses
*/

func (s byFnName) Less(i, j int) bool {
    return concat(s[i].Module.Name, ".", s[i].Name) < concat(s[j].Module.Name, ".", s[j].Name)
}
func (s byTypName) Less(i, j int) bool {
    return s[i].Name < s[j].Name
}
func (s byModName) Less(i, j int) bool {
    return s[i].Name < s[j].Name
}
func (s byDefName) Less(i, j int) bool {
    return concat(s[i].Module.Name, ".", s[i].Name) < concat(s[j].Module.Name, ".", s[j].Name)
}
func (s byStrctName) Less(i, j int) bool {
    return concat(s[i].Module.Name, ".", s[i].Name) < concat(s[j].Module.Name, ".", s[j].Name)
}
func (s byFldName) Less(i, j int) bool {
    return s[i].Name < s[j].Name
}
func (s byParamName) Less(i, j int) bool {
    return s[i].Name < s[j].Name
}
