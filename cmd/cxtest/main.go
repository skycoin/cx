package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/skycoin/cx/cmd/cxtest/runner"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "cxtest",
		Usage: "cx programs tester",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cxpath",
				Usage:       "cx binary path",
				DefaultText: "../bin/cx",
			},
			&cli.StringFlag{
				Name:        "wdir",
				Usage:       "working directory with *.cx tests",
				DefaultText: "./tests",
			},
			&cli.StringFlag{
				Name:  "log",
				Usage: "Enable logMask set (all, success, stderr, fail, skip, time)",
			},
			&cli.StringFlag{
				Name:  "enable-tests",
				Usage: "Enable test set (all, stable, issue, gui)",
			},
			&cli.StringFlag{
				Name:  "disable-tests",
				Usage: "Disable test set (all, stable, issue, gui)",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Print debug information",
			},
		},
		Action: func(c *cli.Context) error {
			return Execute(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Execute(c *cli.Context) error {
	cxPath := c.String("cxpath")
	if cxPath == "" {
		cxPath = "../bin/cx"
	}

	workingDir := c.String("wdir")
	if workingDir == "" {
		workingDir = "./tests"
	}

	debug := c.Bool("debug")

	var parseBitMask = func(flagName string, bitsMap map[string]runner.Bits, defaultBit runner.Bits) runner.Bits {
		var mask runner.Bits = 0
		flags := strings.Split(c.String(flagName), ",")
		for _, flag := range flags {
			mask = runner.Set(mask, bitsMap[flag])
		}
		if debug {
			fmt.Printf("Parsed bit mask for flag %s%s: %06b\n", flagName, flags, mask)
		}
		return mask
	}

	logMask := parseBitMask("log", runner.LogBits, runner.LogNone)
	enableTestsMaks := parseBitMask("enable-tests", runner.TestBits, runner.TestNone)
	disabledTestsMaks := parseBitMask("disable-tests", runner.TestBits, runner.TestNone)

	if enableTestsMaks == runner.TestAll && disabledTestsMaks == runner.TestAll {
		return errors.New("invalid test flags combination")
	}

	var testsMask runner.Bits = runner.TestAll
	if debug {
		fmt.Printf("Initial test mask: %06b\n", testsMask)
	}
	// turn on only enabled tests
	if enableTestsMaks != runner.TestNone {
		testsMask = testsMask & enableTestsMaks
	}
	// turn off only disabled test
	testsMask = testsMask &^ disabledTestsMaks
	if debug {
		fmt.Printf("Resulting test mask: %06b\n", testsMask)
	}

	tester := runner.NewTestRunner(&runner.Config{
		CxPath:         cxPath,
		WorkingDir:     workingDir,
		TestsMask:      testsMask,
		LogMask:        logMask,
		DefaultTimeout: 10 * time.Second,
	})

	var start = time.Now().Unix()

	fmt.Printf("Running CX tests in dir: '%s'\n", workingDir)
	runTestCases(tester)
	end := time.Now().Unix()

	if runner.Has(logMask, runner.LogTime) {
		fmt.Printf("\nTests finished after %d milliseconds", end-start)
	}

	fmt.Printf("\nA total of %d tests were performed\n", tester.TestCount)
	fmt.Printf("%d were successful\n", tester.TestSuccess)
	fmt.Printf("%d failed\n", tester.TestCount-tester.TestSuccess)
	fmt.Printf("%d skipped\n", tester.TestSkipped)

	if tester.TestCount == 0 || (tester.TestSuccess != tester.TestCount) {
		return errors.New("not all test succeeded")
	}

	return nil
}

func runTestCases(t *runner.TestRunner) {
	// tests
	t.Run("test-i8.cx", runner.CxSuccess, "i32")
	t.Run("test-i16.cx", runner.CxSuccess, "i32")
	t.Run("test-i32.cx", runner.CxSuccess, "i32")
	t.Run("test-i64.cx", runner.CxSuccess, "i64")
	t.Run("test-ui8.cx", runner.CxSuccess, "i32")
	t.Run("test-ui16.cx", runner.CxSuccess, "i32")
	t.Run("test-ui32.cx", runner.CxSuccess, "i32")
	t.Run("test-ui64.cx", runner.CxSuccess, "i64")
	t.Run("test-f32.cx", runner.CxSuccess, "f32")
	t.Run("test-f64.cx", runner.CxSuccess, "f64")
	t.Run("test-bool.cx", runner.CxSuccess, "bool")
	t.Run("test-array.cx", runner.CxSuccess, "array")
	t.Run("test-function.cx", runner.CxSuccess, "function")
	t.Run("test-control-flow.cx", runner.CxSuccess, "control floow")
	t.Run("test-utils.cx test-struct.cx", runner.CxSuccess, "struct")
	t.Run("test-str.cx", runner.CxSuccess, "str")
	t.Run("test-utils.cx test-pointers.cx", runner.CxSuccess, "pointers")
	t.Run("test-slices.cx", runner.CxSuccess, "slices")

	t.Run("--cxpath test-workspace test-workspace-a.cx", runner.CxSuccess, "Testing if CX can set a workspace and then import a library, taking that workspace as the new relative path.")
	t.Run("--cxpath test-workspace test-workspace-b.cx", runner.CxSuccess, "Testing if CX can set a workspace and then import a nested library, taking that workspace as the new relative path.")
	t.Run("--cxpath test-workspace test-workspace-c.cx test-workspace-d.cx", runner.CxSuccess, "Testing if files supplied to the CLI override libraries in the workspace.")
	t.Run("test-slices-index-out-of-range-a.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test index < 0")
	t.Run("test-slices-index-out-of-range-b.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test index >= len")
	t.Run("test-slices-resize-out-of-range-a.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test out of range after resize")
	t.Run("test-slices-resize-out-of-range-b.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test resize with count < 0")
	t.Run("test-slices-insert-out-of-range-a.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test insert with index > len")
	t.Run("test-slices-insert-out-of-range-b.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test insert with index < 0")
	t.Run("test-slices-remove-out-of-range-a.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test remove with index < 0")
	t.Run("test-slices-remove-out-of-range-b.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test remove with index >= len")
	t.Run("test-slices-remove-out-of-range-c.cx", runner.CxRuntimeSliceIndexOutOfRange, "Test remove with index == 0 && len == 0")
	t.Run("test-short-declarations.cx", runner.CxSuccess, "short declarations")
	t.Run("test-parse.cx", runner.CxSuccess, "parse")
	t.Run("test-collection-functions.cx", runner.CxSuccess, "collection functions")
	t.Run("test-scopes.cx", runner.CxSuccess, "ProgramError in scopes.")
	t.RunEx("-heap-initial 0 test-gc.cx", runner.CxSuccess, "Stress-testing the garbage collector", runner.TestIssue, 0)
	t.Run("../lib/json.cx test-json.cx", runner.CxSuccess, "ProgramError in json lib.")
	t.Run("../lib/args.cx test-args.cx", runner.CxSuccess, "ProgramError in args lib.")
	// t.Run("test-regexp-must-compile-fail.cx", runner.CxRuntimeError, "ProgramError in regexp lib - MustCompile should have thrown an error.")
	// t.Run("test-regexp-compile-fail.cx", runner.CxSuccess, "ProgramError in regexp lib - error thrown by regexp.Compile does not matches expected error.")
	// t.Run("test-regexp.cx", runner.CxSuccess, "ProgramError in regexp lib.")
	// t.Run("test-cipher.cx", runner.CxSuccess, "ProgramError in cipher lib.")
	// t.RunEx("test-regexp.cx", CxCompilationError, "Panic when calling gl.BindBuffer with only one argument.", TestGui|TestStable, 0)

	// issues
	t.Run("issue-207.cx", runner.CxCompilationError, "Type casting error not reported.")
	t.RunEx("issue-208.cx", runner.CxCompilationError, "Panic if return value is not used.", runner.TestGui|runner.TestStable, 0)
	t.Run("issue-214.cx", runner.CxSuccess, "String not working across packages")
	t.Run("issue-215.cx issue-215a.cx", runner.CxSuccess, "Order of files matters for structs")
	t.Run("issue-215a.cx issue-215.cx", runner.CxSuccess, "Order of files matters for structs")
	t.RunEx("issue-216.cx", runner.CxCompilationError, "Panic when calling gl.BindBuffer with only one argument.", runner.TestGui|runner.TestStable, 0)
	t.RunEx("issue-217.cx", runner.CxSuccess, "Panic when giving []f32 argument to gl.BufferData", runner.TestGui, 0)
	t.Run("issue-218.cx", runner.CxSuccess, "Struct field crushed")
	t.Run("issue-219.cx", runner.CxSuccess, "Failed to modify value in an array")
	t.Run("issue-220.cx", runner.CxSuccess, "Panic when trying to index (using a var) an array, member of a struct passed as a function argument")
	t.Run("issue-27.cx", runner.CxSuccess, "Failed to use shorthand operator-assign (+=, etc.) for arithmetic statements")
	t.Run("issue-221.cx", runner.CxSuccess, "Can't call method from package")
	t.Run("issue-222.cx", runner.CxSuccess, "Can't call method if it has a parameter")
	t.Run("issue-223.cx", runner.CxSuccess, "Panic when using arithmetic to index an array field of a struct")
	t.Run("issue-224.cx", runner.CxSuccess, "Panic if return value is used in an expression")
	t.Run("issue-225.cx", runner.CxSuccess, "Using a variable to store the return boolean value of a function doesnt work with an if statement")
	t.Run("issue-226.cx", runner.CxSuccess, "Panic when accessing property of struct array passed in as argument to func")
	t.Run("issue-227.cx", runner.CxSuccess, "Unexpected results when accessing arrays of structs in a struct")
	t.Run("issue-230.cx", runner.CxSuccess, "Inline initializations and arrays")
	t.Run("issue-231.cx", runner.CxSuccess, "Slice keeps growing though it's cleared inside the loop")
	t.Run("issue-232.cx", runner.CxSuccess, "Scope not working in loops")
	t.Run("issue-233.cx", runner.CxSuccess, "Interdependant Structs")
	t.Run("issue-234.cx", runner.CxCompilationError, "Panic when trying to access an invalid field.")
	t.Run("issue-235.cx", runner.CxCompilationError, "No compilation error when using an using an invalid identifier")
	t.Run("issue-236a.cx issue-236.cx", runner.CxSuccess, "Silent name clash between packages")
	t.Run("issue-236.cx issue-236a.cx", runner.CxSuccess, "Silent name clash between packages")
	t.Run("issue-237.cx", runner.CxCompilationError, "Invalid implicit cast.")
	t.Run("issue-238.cx", runner.CxCompilationError, "Panic when using +* in an expression")
	t.Run("issue-239.cx", runner.CxCompilationError, "No compilation error when defining a struct with duplicate fields.")
	t.Run("issue-240.cx", runner.CxSuccess, "Can't define struct with a single character identifier.")
	t.Run("issue-241.cx", runner.CxSuccess, "Panic when variable used in if statement without parenthesis.")
	t.Run("issue-242.cx", runner.CxSuccess, "Struct field stomped")
	t.Run("issue-243.cx", runner.CxCompilationError, "No compilation error when indexing an array with a non integral var.")
	t.Run("issue-244a.cx", runner.CxSuccess, "Panic when a field of a struct returned by a function is used in an expression")
	t.Run("issue-244b.cx", runner.CxSuccess, "Panic when a field of a struct returned by a function is used in an expression")
	t.Run("issue-245a.cx issue-245.cx", runner.CxCompilationError, "No compilation error when using var without package qualification.")
	t.Run("issue-246.cx", runner.CxSuccess, "No compilation error when passing *i32 as an i32 arg and conversely")
	t.Run("issue-246a.cx", runner.CxCompilationError, "No compilation error when passing *i32 as an i32 arg and conversely")
	t.Run("issue-247.cx", runner.CxCompilationError, "No compilation error when dereferencing an i32 var.")
	t.Run("issue-248.cx", runner.CxSuccess, "Wrong pointer behaviour.")
	t.Run("issue-249.cx", runner.CxSuccess, "Return from a function doesnt work")
	t.Run("issue-249b.cx", runner.CxCompilationError, "Mismatched number of returning arguments is not throwing an error")
	t.Run("issue-250.cx", runner.CxCompilationError, "No compilation error when var is accessed outside of its declaring scope")
	t.Run("issue-251.cx", runner.CxCompilationError, "Panic when a str var is shadowed by a struct var in another scope")
	t.RunEx("issue-252.cx", runner.CxSuccess, "glfw.GetCursorPos() throws error", runner.TestGui|runner.TestStable, 0)
	t.Run("issue-253.cx", runner.CxSuccess, "Inline field and index 'dereferences' to function calls' outputs")
	t.Run("issue-254.cx", runner.CxCompilationError, "No compilation error when redeclaring a variable")
	t.Run("issue-255.cx", runner.CxSuccess, "Multi-dimensional slices don't work")
	t.Run("issue-256.cx", runner.CxSuccess, "can't prefix a (f32) variable with minus to flip it's signedness")
	t.Run("issue-257.cx", runner.CxCompilationError, "Using int literal 0 where 0.0 was needed gave no error")
	t.Run("issue-258.cx", runner.CxSuccess, "error with sending references of structs to functions")
	t.Run("issue-258b.cx", runner.CxSuccess, "error with references to struct literals")
	t.Run("issue-259.cx", runner.CxCompilationError, "struct identifier (when initializing fields) can be with or without a '&' prefix, with no CX error")
	t.Run("issue-260.cx", runner.CxCompilationError, "can assign to previously undeclared vars with just '='")
	t.Run("issue-261.cx", runner.CxSuccess, "empty code blocks (even if they contain commented-out lines) crash like this")
	t.Run("issue-262.cx", runner.CxSuccess, "increment operator ++ does not work")
	t.Run("issue-263.cx", runner.CxSuccess, "Method does not work")
	t.Run("issue-264.cx", runner.CxSuccess, "Cannot use bool variable in if expression")
	t.Run("issue-265.cx", runner.CxSuccess, "CX Parser does not recognize method")
	t.Run("issue-266.cx", runner.CxSuccess, "Goto not working on windows")
	t.Run("issue-267.cx", runner.CxSuccess, "Methods with pointer receivers don't work")
	t.RunEx("issue-268.cx", runner.CxSuccess, "when using 2 f32 out parameters, only the value of the 2nd gets through", runner.TestGui, 0)
	t.Run("issue-269.cx", runner.CxCompilationError, "Variable redeclaration should not be allowed")
	t.Run("issue-270.cx", runner.CxSuccess, "Short variable declarations are not working with calls to methods or functions")
	t.Run("issue-271.cx", runner.CxCompilationError, "Panic when using equality operator between a bool and an i32")
	t.Run("issue-272.cx", runner.CxSuccess, "String concatenation using the + operator doesn't work")
	t.Run("issue-273.cx", runner.CxSuccess, "Argument list is not parsed correctly")
	t.Run("issue-274.cx", runner.CxSuccess, "Dubious error message when indexing an array with a substraction expression")
	t.Run("issue-275.cx", runner.CxSuccess, "Dubious error message when inline initializing a slice")
	t.Run("issue-276a.cx issue-276.cx", runner.CxSuccess, "Troubles when accessing a global var from another package")
	t.Run("issue-277.cx", runner.CxSuccess, "same func names (but in different packages) collide")
	t.Run("issue-278.cx", runner.CxCompilationError, "can use vars from other packages without a 'packageName.' prefix")
	t.Run("issue-279.cx", runner.CxSuccess, "False positive when detecting variable redeclaration.")
	t.RunEx("issue-279a.cx", runner.CxSuccess, "False positive when detecting variable redeclaration.", runner.TestIssue, 0)
	t.Run("issue-280.cx", runner.CxSuccess, "Problem with struct literals in short variable declarations")
	t.Run("issue-281.cx", runner.CxSuccess, "Panic when using the return value of a function in a short declaration")
	t.Run("issue-282a.cx", runner.CxCompilationError, "Panic when inserting a new line in a string literal")
	t.Run("issue-282b.cx", runner.CxSuccess, "Panic when inserting a new line in a string literal")
	t.Run("issue-283.cx", runner.CxCompilationError, "Panic when declaring a variable of an unknown type")
	t.Run("issue-284.cx", runner.CxCompilationError, "No compilation error when using arithmetic operators on struct instances")
	t.RunEx("issue-285.cx", runner.CxSuccess, "Parser gets confused with `2 -2`", runner.TestStable, 0)
	t.Run("issue-286.cx", runner.CxSuccess, "Panic in when assigning an empty initializer list to a []i32 variable")
	t.Run("issue-287.cx", runner.CxSuccess, "Cx stack overflow when appending to a slice passed by address")
	t.Run("issue-288.cx", runner.CxCompilationError, "Panic when trying to assign return value of a function returning void")
	t.Run("issue-289.cx", runner.CxCompilationError, "Panic when using a function declared in another package without importing the package")
	t.Run("issue-290.cx", runner.CxSuccess, "Cx memory stomped")
	t.Run("issue-291.cx", runner.CxSuccess, "Invalid offset calculation of non literal strings when appended to a slice")
	t.Run("issue-292.cx", runner.CxCompilationError, "Panic when calling a function from another package where the package name alias a local variable name")
	t.Run("issue-293.cx", runner.CxSuccess, "Garbage memory when passing the address of slice element to a function")
	t.Run("issue-294.cx", runner.CxSuccess, "Type deduction of struct field fails")
	t.Run("issue-295.cx", runner.CxCompilationError, "No compilation error when assigning a i32 value to a []i32 variable")
	t.Run("issue-296.cx", runner.CxCompilationError, "No compilation error when comparing value of different types")
	t.Run("-stack-size 30 issue-297a.cx", runner.CxRuntimeStackOverflowError, "No stack overflow error")
	t.Run("-heap-initial 100 -heap-max 110 issue-297b.cx", runner.CxRuntimeHeapExhaustedError, "No heap exhausted error")
	t.Run("issue-298.cx", runner.CxSuccess, "Argument type deduction failed when passing address of an i32 struct field to a function accepting *i32 argument.")
	t.Run("issue-299.cx", runner.CxCompilationError, "Type checking is not working with receiving variables of unexpected types")
	t.Run("issue-300.cx", runner.CxSuccess, "Crash when using a constant expression in a slice literal expression")
	t.Run("issue-301-a.cx", runner.CxCompilationError, "Can redeclare variables if they are inline initialized")
	t.Run("issue-301-b.cx", runner.CxCompilationError, "Can redeclare variables if they are inline initialized")
	t.Run("issue-302.cx", runner.CxSuccess, "Trying to determine the length of a slice of struct instances throws an error.")
	t.RunEx("issue-68.cx", runner.CxSuccess, "Wrong sprintf behaviour when passing increment expression as argument", runner.TestIssue, 0)
	t.RunEx("issue-67.cx", runner.CxCompilationError, "Panic when using void return value of a function in a for loop expression", runner.TestIssue, 0)
	t.RunEx("issue-66.cx", runner.CxSuccess, "Wrong sprintf behaviour when printing boolean values with %v", runner.TestIssue, 0)
	t.RunEx("issue-65.cx", runner.CxSuccess, "for true {} loop scope is not executed", runner.TestIssue, 0)
	t.RunEx("issue-64.cx", runner.CxSuccess, "for loop using boolean value is not compiling", runner.TestIssue, 0)
	t.Run("issue-303.cx", runner.CxSuccess, "Concatenation of str variables with + operator doesn't work")
	t.Run("issue-304.cx", runner.CxSuccess, "Short declaration doesn't compile with opcode return value")
	t.Run("issue-305.cx", runner.CxSuccess, "Compilation error when struct field is named 'input' or 'output'")
	t.RunEx("issue-63.cx", runner.CxCompilationError, "No compilation error when using empty argument list after function call", runner.TestIssue, 0)
	t.Run("issue-306.cx", runner.CxCompilationError, "No compilation error when using float value in place of boolean expression")
	t.Run("issue-308.cx", runner.CxCompilationError, "Panic when package contains duplicate function signature")
	t.RunEx("issue-62.cx", runner.CxCompilationError, "Left hand side of , is not compiled", runner.TestIssue, 0)
	t.RunEx("issue-61-a.cx", runner.CxSuccess, "Compilation error when using return value of a member method in a expression", runner.TestIssue, 0)
	t.RunEx("issue-61-b.cx", runner.CxSuccess, "Compilation error when using return value of a member method in a expression", runner.TestIssue, 0)
	t.Run("issue-309.cx", runner.CxSuccess, "Compilation error when left hand side of an assignment expression is a struct field")
	t.RunEx("issue-60.cx", runner.CxSuccess, "Crash when INIT_HEAP_SIZE limit is reached ", runner.TestIssue, 0)
	t.RunEx("issue-59-a.cx", runner.CxSuccess, "Crash in garbage collector when heap is resized", runner.TestIssue, 0)
	t.RunEx("issue-59-b.cx", runner.CxSuccess, "Crash in garbage collector when heap is resized", runner.TestIssue, 0)
	t.RunEx("issue-53-a.cx", runner.CxSuccess, "Issues with slice of type T where sizeof T is different than 4 ", runner.TestStable, 0)
	t.RunEx("issue-53-b.cx", runner.CxSuccess, "Issues with slice of type T where sizeof T is different than 4 ", runner.TestStable, 0)
	t.RunEx("issue-53-c.cx", runner.CxSuccess, "Issues with slice of type T where sizeof T is different than 4 ", runner.TestStable, 0)
	t.RunEx("issue-51.cx", runner.CxCompilationError, "No compilation error when global variable is redeclared at local scope", runner.TestIssue, 0)
	t.RunEx("issue-49.cx", runner.CxCompilationError, "No compilation error when assigning an literal which overflow the receiving type", runner.TestIssue, 0)
	t.RunEx("issue-48.cx", runner.CxSuccess, "Cx is not supporting short-circuit evaluation", runner.TestIssue, 0)
	t.RunEx("issue-39.cx", runner.CxSuccess, "func defined with no arguments, called WITH arguments causes PANIC", runner.TestIssue, 0)
	t.Run("issue-2.cx", runner.CxSuccess, "multi-dimensional arrays are not working")
	t.Run("issue-1.cx", runner.CxSuccess, "multi-dimensional slices are not working")
	t.RunEx("issue-120-a.cx", runner.CxCompilationError, "Invalid implicit cast when assigning the result of a math operator to a variable.", runner.TestIssue, 0)
	t.RunEx("issue-120-b.cx", runner.CxCompilationError, "Invalid implicit cast when assigning the result of a math operator to a variable.", runner.TestIssue, 0)
	t.RunEx("issue-120-c.cx", runner.CxCompilationError, "Invalid implicit cast when assigning the result of a math operator to a variable.", runner.TestIssue, 0)
	t.RunEx("issue-121.cx", runner.CxSuccess, "Compilation error when using unary negative operator on a function call", runner.TestIssue, 0)
	t.RunEx("issue-131.cx", runner.CxSuccess, "Panic when using arithmetic operations.", runner.TestIssue, 0)
	t.RunEx("test-ar-1.cx", runner.CxSuccess, "Panic when using string to pointer array.", runner.TestIssue, 0)
	t.RunEx("issue-157.cx", runner.CxSuccess, "expected either 'i32' or 'i64', got 'ident'", runner.TestIssue, 0)
	t.Run("issue-650.cx", runner.CxSuccess, "test should not fail")

	// We need to fix serialization and deserialization as user-callable functions
	// t.RunEx("issue-309.cx", CxSuccess, "Serialization is not taking into account non-default stack sizes.", TestIssue, 0)
	// t.RunEx("issue-310.cx", CxSuccess, "Splitting a serialized program into its blockchain and transaction parts.", TestIssue, 0)
	// t.RunEx("issue-311.cx", CxSuccess, "`CurrentFunction` and `CurrentStruct` are causing errors in programs with more than 1 package.", TestIssue, 0)
	// t.RunEx("issue-312.cx", CxSuccess, "Deserialization is not setting correctly the sizes for the CallStack, HeapStartsAt and StackSize fields of the CXProgram structure.", TestIssue, 0)

	t.Run("issue-28-f32.cx", runner.CxSuccess, "f32")
	t.Run("issue-28-f64.cx", runner.CxSuccess, "f64")
	t.Run("issue-28-i8.cx", runner.CxSuccess, "i8")
	t.Run("issue-28-i16.cx", runner.CxSuccess, "i16")
	t.Run("issue-28-i32.cx", runner.CxSuccess, "i32")
	t.Run("issue-28-i64.cx", runner.CxSuccess, "i64")
	t.Run("issue-28-ui8.cx", runner.CxSuccess, "ui8")
	t.Run("issue-28-ui16.cx", runner.CxSuccess, "ui16")
	t.Run("issue-28-ui32.cx", runner.CxSuccess, "ui32")
	t.Run("issue-28-ui64.cx", runner.CxSuccess, "ui64")
}
