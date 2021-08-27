// +build cxfx

package cxfx

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cx/types"
	"sync"
	"time"
)

const (
	ACTION_RELEASE ActionType = 0 // TODO : break from glfw.Action
	ACTION_PRESS              = 1
	ACTION_REPEAT             = 2
	ACTION_MOVE               = 3
)

const (
	APP_NONE EventType = iota
	APP_START
	APP_STOP
	APP_KEYBOARD
	APP_MOUSE
	APP_FOCUS_ON
	APP_FOCUS_OFF
	APP_FRAMEBUFFER_SIZE
	APP_WINDOW_SIZE
	APP_WINDOW_POSITION
	APP_PAINT
	APP_CURSOR_POS   // TODO : to deprecate
	APP_MOUSE_BUTTON // TODO : to deprecate
)

type ActionType uint32
type EventType uint32
type Event struct {
	eventType EventType
	action    ActionType
	time      uint64
	x         float64
	y         float64
	key       int32
	scancode  int64
	mods      int32
}

type CXCallback struct {
	windowNameBytes []byte
	windowName      string
	packageName     string
	functionName    string
}

func (cb *CXCallback) init(inputs []ast.CXValue, outputs []ast.CXValue, packageName string) {
	cb.windowName = inputs[0].Get_str()
	var windowHeapPtr = types.AllocWrite_obj_data(ast.PROGRAM.Memory, []byte(cb.windowName))
	var windowName [types.POINTER_SIZE]byte
	types.Write_ptr(windowName[:], 0, windowHeapPtr)
	cb.windowNameBytes = windowName[:]
	cb.functionName = inputs[1].Get_str()
	cb.packageName = packageName
}

func (cb *CXCallback) Init(inputs []ast.CXValue, outputs []ast.CXValue) {
	cb.init(inputs, outputs, inputs[0].Expr.Package.Name)
}

func (cb *CXCallback) InitEx(inputs []ast.CXValue, outputs []ast.CXValue) {
	cb.init(inputs, outputs, inputs[2].Get_str())
}

func (cb *CXCallback) Call(inputs [][]byte) {
	if fn, err := ast.PROGRAM.GetFunction(cb.functionName, cb.packageName); err == nil {
		execute.Callback(ast.PROGRAM, fn, inputs)
	}
}

var appKeyboardCallback CXCallback
var appMouseCallback CXCallback
var appFramebufferSizeCallback CXCallback
var appWindowSizeCallback CXCallback
var appWindowPosCallback CXCallback
var appStartCallback CXCallback
var appStopCallback CXCallback

var appCursorPositionCallback CXCallback // TODO : to deprecate
var appMouseButtonCallback CXCallback    // TODO : to deprecate

var eventCount int

var events []Event
var polled []Event
var polledCount int

var mutex sync.Mutex

func PushEvent(eventType EventType) {
	mutex.Lock()
	defer mutex.Unlock()

	event := Event{
		eventType,
		ACTION_RELEASE, // TODO : ACTION_NONE
		uint64(time.Now().UnixNano()),
		0.0,
		0.0,
		0,
		0,
		0,
	}
	if eventCount < len(events) {
		events[eventCount] = event
	} else {
		events = append(events, event)
	}
	eventCount++
}

func PushWindowSizeEvent(width float64, height float64) {
	index := eventCount
	PushEvent(APP_WINDOW_SIZE)
	events[index].x = width
	events[index].y = height
}

func PushWindowPositionEvent(x float64, y float64) {
	index := eventCount
	PushEvent(APP_WINDOW_POSITION)
	events[index].x = x
	events[index].y = y
}

func PushFramebufferSizeEvent(width float64, height float64) {
	index := eventCount
	PushEvent(APP_FRAMEBUFFER_SIZE)
	events[index].x = width
	events[index].y = height
}

func PushKeyboardEvent(action ActionType, key int32, scancode int64, mods int32) {
	index := eventCount
	PushEvent(APP_KEYBOARD)
	events[index].action = action
	events[index].key = key
	events[index].scancode = scancode
	events[index].mods = mods
}

func PushMouseEvent(eventType EventType, action ActionType, button int32, touch int64, mods int32, x float64, y float64) {
	index := eventCount
	PushEvent(eventType)
	events[index].action = action
	events[index].key = button
	events[index].scancode = touch
	events[index].mods = mods
	events[index].x = x
	events[index].y = y
}

func purgeEvents() {
	mutex.Lock()
	defer mutex.Unlock()

	if eventCount > 0 {
		polled = append(polled[0:0], events[0:eventCount]...)
		eventCount = 0
	}
}

func fromI32(value int32) []byte {
	var memory [4]byte
	types.Write_i32(memory[:], 0, value)
	return memory[:]
}

func fromI64(value int64) []byte {
	var memory [8]byte
	types.Write_i64(memory[:], 0, value)
	return memory[:]
}

func fromF64(value float64) []byte {
	var memory [8]byte
	types.Write_f64(memory[:], 0, value)
	return memory[:]
}

func PollEvents() {
	purgeEvents()
	eventCount := len(polled)
	for i := 0; i < eventCount; i++ {
		e := polled[i]
		switch e.eventType {
		case APP_START:
			var inputs [][]byte = make([][]byte, 1)
			inputs[0] = appStartCallback.windowNameBytes
			appStartCallback.Call(inputs)
		case APP_STOP:
			var inputs [][]byte = make([][]byte, 1)
			inputs[0] = appStopCallback.windowNameBytes
			appStopCallback.Call(inputs)
		case APP_KEYBOARD:
			var inputs [][]byte = make([][]byte, 5)
			inputs[0] = appKeyboardCallback.windowNameBytes
			inputs[1] = fromI32(e.key)
			inputs[2] = fromI32(int32(e.scancode))
			inputs[3] = fromI32(int32(e.action))
			inputs[4] = fromI32(e.mods)
			appKeyboardCallback.Call(inputs)
		case APP_MOUSE:
			var inputs [][]byte = make([][]byte, 7)
			inputs[0] = appMouseCallback.windowNameBytes
			inputs[1] = fromI32(e.key)
			inputs[2] = fromI64(e.scancode)
			inputs[3] = fromI32(int32(e.action))
			inputs[4] = fromI32(e.mods)
			inputs[5] = fromF64(e.x)
			inputs[6] = fromF64(e.y)
			appMouseCallback.Call(inputs)
		case APP_FRAMEBUFFER_SIZE:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appFramebufferSizeCallback.windowNameBytes
			inputs[1] = fromI32(int32(e.x)) // TODO : use float64 (deprecate int32)
			inputs[2] = fromI32(int32(e.y)) // TODO : use float64 (deprecate int32)
			appFramebufferSizeCallback.Call(inputs)
		case APP_WINDOW_SIZE:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appWindowSizeCallback.windowNameBytes
			inputs[1] = fromI32(int32(e.x)) // TODO : use float64 (deprecate int32)
			inputs[2] = fromI32(int32(e.y)) // TODO : use float64 (deprecate int32)
			appWindowSizeCallback.Call(inputs)
		case APP_WINDOW_POSITION:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appWindowPosCallback.windowNameBytes
			inputs[1] = fromI32(int32(e.x)) // TODO : use float64 (deprecate int32)
			inputs[2] = fromI32(int32(e.y)) // TODO : use float64 (deprecate int32)
			appWindowPosCallback.Call(inputs)
		case APP_CURSOR_POS: // TODO : to deprecate
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appCursorPositionCallback.windowNameBytes
			inputs[1] = fromF64(e.x)
			inputs[2] = fromF64(e.y)
			appCursorPositionCallback.Call(inputs)
		case APP_MOUSE_BUTTON: // TODO to deprecate
			var inputs [][]byte = make([][]byte, 4)
			inputs[0] = appMouseButtonCallback.windowNameBytes
			inputs[1] = fromI32(e.key)
			inputs[2] = fromI32(int32(e.action))
			inputs[3] = fromI32(e.mods)
			appMouseButtonCallback.Call(inputs)
		}
	}
	polled = polled[0:0]
}
