package declaration_extractor

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sync"
)

//Regexes
var reName = regexp.MustCompile(`[_a-zA-Z][_a-zA-Z0-9]*`)
var reEnumInit = regexp.MustCompile(`const\s+\(`)
var reEnumDec = regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)(?:\s+([_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*)){0,1}(?:\s*\=\s*[\s\S]+\S+){0,1}`)
var reNotSpace = regexp.MustCompile(`\S+`)

// Func Declaration regex for name extraction and syntax checking
// Components:
// func - func keyword
// (?:\s*\(\s*[_a-zA-Z]\w*\s+\*{0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*\)\s*)|\s+) -  [space/no space] [([reciever object]) [name] | [space]]
// ([_a-zA-Z]\w*) - name of func
// (?:\s*\(\s*(?:(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1})\s*\)){1,2} - [[space/no space] ([params])]{1,2}
//		(?:(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1}) - [[param name] [data type] [,]]{0,1} [param name] [data type] | [param name] [data type]
// 			(?:[_a-zA-Z]\w*\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)* - [param name] [*]{0,1} [\[[1-9][0-9]+|[0-9]\][*]{0,1}]{0,1} [word] [[.][word]]*
//
// First, finds the func keyword
// Second, finds out whether the function has a receiver object or not and extracts the func name
// Third, finds whether there's one or two pairs of parenthesis
// Forth, finds whether there's one or more params in the parenthesis
var reFuncDec = regexp.MustCompile(`(func(?:(?:\s*\(\s*[_a-zA-Z]\w*\s+\*{0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*\)\s*)|\s+)([_a-zA-Z]\w*)(?:\s*\(\s*(?:(?:[_a-zA-Z]\w*(?:\s*\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*|\s+)[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*\s*,\s*)+[_a-zA-Z]\w*(?:\s*\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}|\s+)\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*|(?:[_a-zA-Z]\w*(?:\s*\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*|\s+)[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*){0,1})\s*\)){1,2})(?:\s*{){0,1}`)
var reGlobalName = regexp.MustCompile(`var\s+([_a-zA-Z]\w*)\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*(?:\s*\=\s*[\s\S]+\S+){0,1}`)
var reStruct = regexp.MustCompile(`type\s+[_a-zA-Z][_a-zA-Z0-9]*\s+struct`)
var reStructHeader = regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+struct\s*{`)
var reStructField = regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)\s+\*{0,1}\s*(?:\[(?:[1-9]\d+|[0-9]){0,1}\]\*{0,1}){0,1}\s*[_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*`)
var reTypeDefinition = regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+([_a-zA-Z]\w*(?:\.[_a-zA-Z]\w*)*)`)

func ReDeclarationCheck(Import []ImportDeclaration, Glbl []GlobalDeclaration, Enum []EnumDeclaration, TypeDef []TypeDefinitionDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) error {

	// Checks for the first declaration redeclared
	// in the order:
	// Import -> Global -> Enum -> Type Definition -> Struct -> Func

	for i := 0; i < len(Import); i++ {
		for j := i + 1; j < len(Import); j++ {
			if Import[i].ImportName == Import[j].ImportName && Import[i].PackageID == Import[j].PackageID && Import[i].FileID == Import[j].FileID {
				return fmt.Errorf("%v:%v: redeclaration error: import: %v", filepath.Base(Import[j].FileID), Import[j].LineNumber, Import[i].ImportName)
			}
		}
	}

	for i := 0; i < len(Glbl); i++ {
		for j := i + 1; j < len(Glbl); j++ {
			if Glbl[i].GlobalVariableName == Glbl[j].GlobalVariableName && Glbl[i].PackageID == Glbl[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: global: %v", filepath.Base(Glbl[j].FileID), Glbl[j].LineNumber, Glbl[i].GlobalVariableName)
			}
		}
	}

	for i := 0; i < len(Enum); i++ {
		for j := i + 1; j < len(Enum); j++ {
			if Enum[i].EnumName == Enum[j].EnumName && Enum[i].PackageID == Enum[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: enum: %v", filepath.Base(Enum[j].FileID), Enum[j].LineNumber, Enum[i].EnumName)
			}
		}
	}

	for i := 0; i < len(TypeDef); i++ {
		for j := i + 1; j < len(TypeDef); j++ {
			if TypeDef[i].TypeDefinitionName == TypeDef[j].TypeDefinitionName && TypeDef[i].PackageID == TypeDef[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: type definition: %v", filepath.Base(TypeDef[j].FileID), TypeDef[j].LineNumber, TypeDef[i].TypeDefinitionName)
			}
		}
	}

	for i := 0; i < len(Strct); i++ {

		StructFields := Strct[i].StructFields
		for m := 0; m < len(StructFields); m++ {
			for n := m + 1; n < len(StructFields); n++ {
				if StructFields[m].StructFieldName == StructFields[n].StructFieldName {
					return fmt.Errorf("%v:%v: redeclaration error: struct field: %v", filepath.Base(Strct[i].FileID), StructFields[n].LineNumber, StructFields[n].StructFieldName)
				}
			}
		}

		for j := i + 1; j < len(Strct); j++ {
			if Strct[i].StructName == Strct[j].StructName && Strct[i].PackageID == Strct[j].PackageID {
				return fmt.Errorf("%v:%v: redeclaration error: struct: %v", filepath.Base(Strct[j].FileID), Strct[j].LineNumber, Strct[i].StructName)
			}
		}
	}

	for i := 0; i < len(Func); i++ {
		for j := i + 1; j < len(Func); j++ {
			if Func[i].FuncName == Func[j].FuncName && Func[i].PackageID == Func[j].PackageID {
				funcMethod1, err := ExtractMethod(Func[i])
				if err != nil {
					return err
				}
				funcMethod2, err := ExtractMethod(Func[j])
				if err != nil {
					return err
				}
				if funcMethod1 == funcMethod2 {
					return fmt.Errorf("%v:%v: redeclaration error: func: %v", filepath.Base(Func[j].FileID), Func[j].LineNumber, Func[i].FuncName)
				}
			}
		}
	}

	return nil
}

func GetDeclarations(source []byte, Glbls []GlobalDeclaration, Enums []EnumDeclaration, TypeDefs []TypeDefinitionDeclaration, Strcts []StructDeclaration, Funcs []FuncDeclaration) []string {

	var declarations []string

	for _, glbl := range Glbls {
		declarations = append(declarations, string(source[glbl.StartOffset:glbl.StartOffset+glbl.Length]))
	}

	for _, enum := range Enums {
		declarations = append(declarations, string(source[enum.StartOffset:enum.StartOffset+enum.Length]))
	}

	for _, typeDef := range TypeDefs {
		declarations = append(declarations, string(source[typeDef.StartOffset:typeDef.StartOffset+typeDef.Length]))
	}

	for _, strct := range Strcts {
		declarations = append(declarations, string(source[strct.StartOffset:strct.StartOffset+strct.Length]))
	}

	for _, fun := range Funcs {
		declarations = append(declarations, string(source[fun.StartOffset:fun.StartOffset+fun.Length]))
	}

	return declarations
}

func ExtractAllDeclarations(sourceCodeStrings []string, fileNames []string) ([]ImportDeclaration, []GlobalDeclaration, []EnumDeclaration, []TypeDefinitionDeclaration, []StructDeclaration, []FuncDeclaration, error) {

	//Variable declarations
	var Imports []ImportDeclaration
	var Globals []GlobalDeclaration
	var Enums []EnumDeclaration
	var TypeDefinitions []TypeDefinitionDeclaration
	var Structs []StructDeclaration
	var Funcs []FuncDeclaration

	//Channel declarations
	importChannel := make(chan []ImportDeclaration, len(sourceCodeStrings))
	globalChannel := make(chan []GlobalDeclaration, len(sourceCodeStrings))
	enumChannel := make(chan []EnumDeclaration, len(sourceCodeStrings))
	typeDefinitionChannel := make(chan []TypeDefinitionDeclaration, len(sourceCodeStrings))
	structChannel := make(chan []StructDeclaration, len(sourceCodeStrings))
	funcChannel := make(chan []FuncDeclaration, len(sourceCodeStrings))
	errorChannel := make(chan error, len(sourceCodeStrings))

	var wg sync.WaitGroup

	// concurrent extractions start
	for i, sourceCode := range sourceCodeStrings {

		wg.Add(1)

		go func(sourceCode string, fileName string, globalChannel chan<- []GlobalDeclaration, enumChannel chan<- []EnumDeclaration, typeDefinition chan<- []TypeDefinitionDeclaration, structChannel chan<- []StructDeclaration, funcChannel chan<- []FuncDeclaration, errorChannel chan<- error, wg *sync.WaitGroup) {

			defer wg.Done()

			srcBytes := []byte(sourceCode)
			replaceComments := ReplaceCommentsWithWhitespaces(srcBytes)
			replaceStringContents, err := ReplaceStringContentsWithWhitespaces(replaceComments)
			if err != nil {
				errorChannel <- fmt.Errorf("%v:%v", filepath.Base(fileName), err)
				return
			}

			wg.Add(6)

			go func(importChannel chan<- []ImportDeclaration, replaceComments []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				imports, err := ExtractImports(replaceComments, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				importChannel <- imports

			}(importChannel, replaceComments, fileName, wg)

			go func(globalChannel chan<- []GlobalDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				globals, err := ExtractGlobals(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				globalChannel <- globals

			}(globalChannel, replaceStringContents, fileName, wg)

			go func(enumChannel chan<- []EnumDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				enums, err := ExtractEnums(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				enumChannel <- enums

			}(enumChannel, replaceStringContents, fileName, wg)

			go func(typeDefinitionChannel chan<- []TypeDefinitionDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				typeDefinitions, err := ExtractTypeDefinitions(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				typeDefinitionChannel <- typeDefinitions

			}(typeDefinitionChannel, replaceStringContents, fileName, wg)

			go func(structChannel chan<- []StructDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				structs, err := ExtractStructs(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				structChannel <- structs

			}(structChannel, replaceStringContents, fileName, wg)

			go func(funcChannel chan<- []FuncDeclaration, replaceStringContents []byte, fileName string, wg *sync.WaitGroup) {

				defer wg.Done()

				funcs, err := ExtractFuncs(replaceStringContents, fileName)

				if err != nil {
					errorChannel <- err
					return
				}

				funcChannel <- funcs

			}(funcChannel, replaceStringContents, fileName, wg)

		}(sourceCode, fileNames[i], globalChannel, enumChannel, typeDefinitionChannel, structChannel, funcChannel, errorChannel, &wg)
	}

	wg.Wait()

	// Close all channels for reading
	close(importChannel)
	close(globalChannel)
	close(enumChannel)
	close(typeDefinitionChannel)
	close(structChannel)
	close(funcChannel)
	close(errorChannel)

	//Read from channels concurrently
	wg.Add(6)

	go func() {

		for imprt := range importChannel {
			Imports = append(Imports, imprt...)
		}

		wg.Done()

	}()

	go func() {

		for global := range globalChannel {
			Globals = append(Globals, global...)
		}

		wg.Done()

	}()

	go func() {

		for enum := range enumChannel {
			Enums = append(Enums, enum...)
		}

		wg.Done()
	}()

	go func() {

		for typeDef := range typeDefinitionChannel {
			TypeDefinitions = append(TypeDefinitions, typeDef...)
		}

		wg.Done()

	}()

	go func() {

		for strct := range structChannel {
			Structs = append(Structs, strct...)
		}

		wg.Done()

	}()

	go func() {

		for fun := range funcChannel {
			Funcs = append(Funcs, fun...)
		}

		wg.Done()

	}()

	wg.Wait()

	// there's an error, return values with first error
	if err := <-errorChannel; err != nil {
		return Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, err
	}

	reDeclarationCheck := ReDeclarationCheck(Imports, Globals, Enums, TypeDefinitions, Structs, Funcs)

	// there's declaration redeclared return values with error
	if reDeclarationCheck != nil {
		return Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, reDeclarationCheck
	}

	return Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, nil
}
