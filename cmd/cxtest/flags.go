package main

const (
	TestNone   Bits = 0
	TestStable Bits = 1 << iota
	TestIssue
	TestGui
	TestAll = TestStable | TestIssue | TestGui
)

var TestBits = map[string]Bits{
	"all":    TestAll,
	"stable": TestStable,
	"issue":  TestIssue,
	"gui":    TestGui,
}

const (
	LogNone    Bits = 0
	LogSuccess Bits = 1 << iota
	LogStderr
	LogFail
	LogSkip
	LogTime
	LogAll = LogSuccess | LogStderr | LogFail | LogSkip | LogTime
)

var LogBits = map[string]Bits{
	"all":     LogAll,
	"success": LogSuccess,
	"stderr":  LogStderr,
	"fail":    LogFail,
	"skip":    LogSkip,
	"time":    LogTime,
}

const (
	CxSuccess = iota
	CxCompilationError
	CxPanic // 2
	CxInternalError
	CxAssert
	CxRuntimeError
	CxRuntimeStackOverflowError
	CxRuntimeHeapExhaustedError
	CxRuntimeInvalidArgument
	CxRuntimeSliceIndexOutOfRange
	CxRuntimeNotImplemented
)
