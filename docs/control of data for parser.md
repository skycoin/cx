

function parseProgram take source code as input 

task I
create program 
goto
actions.PRGRM = cxcore.MakeProgram()

task II 
ParseSourceCode
cxgo.ParseSourceCode(sourceCode, fileNames)

goto
ParseSourceCode 

goto
lexerStep0
// lexerStep0 performs a first pass for the CX parser. Globals, packages and
// custom types are added to `cxgo0.PRGRM0`

cxgo0.PRGRM0 has all data into Packages []*CXPackage varaiable.