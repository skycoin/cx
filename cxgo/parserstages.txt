before lexerstep0 cxgo0.Parse() function call 
Program
0.- Package: main
	Globals
		0.- Global: i 
	Structs
		0.- Struct: Point

after lexerstep0 cxgo0.Parse() function call 
Program
0.- Package: main
	Imports
		0.- Import: http
	Globals
		0.- Global: i i32
	Structs
		0.- Struct: Point
			0.- Field: x i32
			1.- Field: y i32
	Functions
		0.- Function: main () ()
		1.- Function: sayHi () ()

after lexerStep0
Program
0.- Package: main
	Imports
		0.- Import: http
	Globals
		0.- Global: i i32
	Structs
		0.- Struct: Point
			0.- Field: x i32
			1.- Field: y i32
	Functions
		0.- Function: main () ()
		1.- Function: sayHi () ()

before cxgo.Parse(cxgo.NewLexer(b) 
Program
0.- Package: main
	Imports
		0.- Import: http
	Globals
		0.- Global: i i32
	Structs
		0.- Struct: Point
			0.- Field: x i32
			1.- Field: y i32
	Functions
		0.- Function: main () ()
		1.- Function: sayHi () ()

after cxgo.Parse(cxgo.NewLexer(b) 
Program
0.- Package: main
	Imports
		0.- Import: http
	Globals
		0.- Global: i i32
	Structs
		0.- Struct: Point
			0.- Field: x i32
			1.- Field: y i32
	Functions
		0.- Function: main () ()
			0.- Declaration: loco i32
			1.- Expression: str.print(&*lit str)
			2.- Expression: str.print(&*lit str)
			3.- Expression: str.print(&*lit str)
			4.- Expression: *tmp_0 i32 = mul(*lit i32, *lit i32)
			5.- Expression: *tmp_1 i32 = add(*lit i32, *tmp_0 i32)
			6.- Expression: i32.print(*tmp_1 i32)
			7.- Expression: str.print(&*lit str)
			8.- Expression: sayHi()
			9.- Declaration: p Point
			10.- Expression: p.x i32 = identity(*lit i32)
			11.- Expression: p.y i32 = identity(*lit i32)
			12.- Expression: str.print(&*lit str)
			13.- Expression: i32.print(p.x i32)
			14.- Expression: i32.print(p.y i32)
			15.- Expression: i i32 = identity(*lit i32)
			16.- Expression: str.print(&*lit str)
			17.- Expression: i32.print(i i32)
			18.- Expression: loco i32 = identity(*lit i32)
			19.- Expression: str.print(&*lit str)
			20.- Expression: i32.print(loco i32)
		1.- Function: sayHi () ()
			0.- Expression: str.print(&*lit str)

after ParseSourceCode
Program
0.- Package: main
	Imports
		0.- Import: http
	Globals
		0.- Global: i i32
	Structs
		0.- Struct: Point
			0.- Field: x i32
			1.- Field: y i32
	Functions
		0.- Function: main () ()
			0.- Declaration: loco i32
			1.- Expression: str.print(&*lit str)
			2.- Expression: str.print(&*lit str)
			3.- Expression: str.print(&*lit str)
			4.- Expression: *tmp_0 i32 = mul(*lit i32, *lit i32)
			5.- Expression: *tmp_1 i32 = add(*lit i32, *tmp_0 i32)
			6.- Expression: i32.print(*tmp_1 i32)
			7.- Expression: str.print(&*lit str)
			8.- Expression: sayHi()
			9.- Declaration: p Point
			10.- Expression: p.x i32 = identity(*lit i32)
			11.- Expression: p.y i32 = identity(*lit i32)
			12.- Expression: str.print(&*lit str)
			13.- Expression: i32.print(p.x i32)
			14.- Expression: i32.print(p.y i32)
			15.- Expression: i i32 = identity(*lit i32)
			16.- Expression: str.print(&*lit str)
			17.- Expression: i32.print(i i32)
			18.- Expression: loco i32 = identity(*lit i32)
			19.- Expression: str.print(&*lit str)
			20.- Expression: i32.print(loco i32)
		1.- Function: sayHi () ()
			0.- Expression: str.print(&*lit str)

print string
Hello World!
print 2+3*3
11
call function sayhi()
Hi
print  struct point
100
200
print  global varialble
20
print  local varialble
20
