// +build cxfx

package cxfx

import (
	"fmt"
	. "github.com/SkycoinProject/cx/cx"
	"sync"
	"time"
)

const (
	APP_NONE EventType = iota
	APP_START
	APP_STOP
	APP_KEY
	APP_MOUSE_BUTTON
	APP_CURSOR_POSITION
	APP_TOUCH_PRESS
	APP_TOUCH_MOVE
	APP_TOUCH_RELEASE
	APP_FOCUS_ON
	APP_FOCUS_OFF
	APP_FRAMEBUFFER_SIZE
	APP_WINDOW_SIZE
	APP_WINDOW_POSITION
	APP_PAINT
)

type EventType uint32
type Event struct {
	etype    EventType
	time     uint64
	width    int32
	height   int32
	index    int32
	ix       int32
	iy       int32
	dx       float64
	dy       float64
	key      int32
	scancode int32
	action   int32
	mods     int32
}

type CXCallback struct {
	prgrm           *CXProgram
	expr            *CXExpression
	fp              int
	windowNameBytes []byte
	windowName      string
	packageName     string
	functionName    string
}

func (cb *CXCallback) init(prgrm *CXProgram, expr *CXExpression, fp int, packageName string) {
	cb.prgrm = prgrm
	cb.expr = expr
	cb.fp = fp
	cb.windowName = ReadStr(fp, expr.Inputs[0])
	cb.windowNameBytes = FromI32(int32(NewWriteObj(FromStr(cb.windowName))))
	cb.functionName = ReadStr(fp, expr.Inputs[1])
	cb.packageName = packageName
}

func (cb *CXCallback) Init(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	cb.init(prgrm, expr, fp, expr.Package.Name)
}

func (cb *CXCallback) InitEx(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	cb.init(prgrm, expr, fp, ReadStr(fp, expr.Inputs[2]))
}

func (cb *CXCallback) Call(inputs [][]byte) {
	fmt.Printf("CALLBACK %s, %s\n", cb.packageName, cb.functionName)
	cb.prgrm.Callback(cb.expr, cb.functionName, cb.packageName, inputs)
}

var appKeyCallback CXCallback
var appCursorPositionCallback CXCallback
var appMouseButtonCallback CXCallback
var appTouchCallback CXCallback
var appFramebufferSizeCallback CXCallback
var appWindowSizeCallback CXCallback
var appWindowPosCallback CXCallback
var appStartCallback CXCallback
var appStopCallback CXCallback

var eventCount int

var events []Event
var polled []Event
var polledCount int

var mutex sync.Mutex

func PushEvent(etype EventType) {
	mutex.Lock()
	defer mutex.Unlock()

	event := Event{
		etype,
		uint64(time.Now().UnixNano()),
		0,
		0,
		0,
		0,
		0.0,
		0.0,
		0,
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
	fmt.Printf("PUUUUUUUUUUUUUUUUUUUUUUSH %d\n", eventCount)
}

func PushWindowSizeEvent(width int32, height int32) {
	fmt.Printf("PushWindowSizeEvent %d, %d\n", width, height)
	index := eventCount
	PushEvent(APP_WINDOW_SIZE)
	events[index].width = width
	events[index].height = height
	fmt.Printf("PUSHED %d, %d\n", events[index].width, events[index].height)
}

func PushWindowPositionEvent(x int32, y int32) {
	index := eventCount
	PushEvent(APP_WINDOW_POSITION)
	events[index].ix = x
	events[index].iy = y
}

func PushFramebufferSizeEvent(width int32, height int32) {
	fmt.Printf("PushFramebufferSizeEvent %d, %d\n", width, height)
	index := eventCount
	PushEvent(APP_FRAMEBUFFER_SIZE)
	events[index].width = width
	events[index].height = height
	fmt.Printf("PUSHED %d, %d\n", events[index].width, events[index].height)
}

func PushKeyEvent(key int32, scancode int32, action int32, mods int32) {
	index := eventCount
	PushEvent(APP_KEY)
	events[index].key = key
	events[index].scancode = scancode
	events[index].action = action
	events[index].mods = mods
}

func PushCursorPositionEvent(x float64, y float64) {
	index := eventCount
	PushEvent(APP_CURSOR_POSITION)
	events[index].dx = x
	events[index].dy = y
}

func PushMouseButtonEvent(key int32, action int32, mods int32) {
	index := eventCount
	PushEvent(APP_MOUSE_BUTTON)
	events[index].key = key
	events[index].action = action
	events[index].mods = mods
}

func PushTouchEvent(touchType int32, x int32, y int32, index int32) {
	index := eventCount
	etype := APP_NONE
	switch touchType {
	case 0:
		etype = APP_TOUCH_PRESS
	case 1:
		etype = APP_TOUCH_MOVE
	case 2:
		etype == APP_TOUCH_RELEASE
	}

	PushEvent(etype)
	events[index].touch = index
	events[index].ix = x
	events[index].iy = y
}

func purgeEvents() {
	mutex.Lock()
	defer mutex.Unlock()

	if eventCount > 0 {
		/*fmt.Printf("PURGE %d-- SUM %d\n", eventCount, len(events))
		for i := 0; i < eventCount; i++ {
			e := events[i]
			switch e.etype {
			case APP_WINDOW_SIZE:
				fmt.Printf("PURGE APP_WINDOW_SIZE  %d, %d\n", e.width, e.height)
			case APP_FRAMEBUFFER_SIZE:
				fmt.Printf("PURGE APP_FRAMEBUFFER_SIZE %d, %d\n", e.width, e.height)
			default:
				fmt.Printf("PURGE ELSE %d\n", e.etype)
			}
		}*/
		polled = append(polled[0:0], events[0:eventCount]...)
		eventCount = 0
	}
}

func PollEvents() {
	fmt.Printf("---------------> POLL_EVENTS\n")
	purgeEvents()
	eventCount := len(polled)
	for i := 0; i < eventCount; i++ {
		e := polled[i]
		switch e.etype {
		case APP_START:
			var inputs [][]byte = make([][]byte, 1)
			inputs[0] = appStartCallback.windowNameBytes
			appStartCallback.Call(inputs)
		case APP_STOP:
			var inputs [][]byte = make([][]byte, 1)
			inputs[0] = appStopCallback.windowNameBytes
			appStopCallback.Call(inputs)
		case APP_KEY:
			var inputs [][]byte = make([][]byte, 5)
			inputs[0] = appKeyCallback.windowNameBytes
			inputs[1] = FromI32(e.key)
			inputs[2] = FromI32(e.scancode)
			inputs[3] = FromI32(e.action)
			inputs[4] = FromI32(e.mods)
			appKeyCallback.Call(inputs)
		case APP_CURSOR_POSITION:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appCursorPositionCallback.windowNameBytes
			inputs[1] = FromF64(e.dx)
			inputs[2] = FromF64(e.dy)
			appCursorPositionCallback.Call(inputs)
		case APP_MOUSE_BUTTON:
			var inputs [][]byte = make([][]byte, 4)
			inputs[0] = appMouseButtonCallback.windowNameBytes
			inputs[1] = FromI32(e.key)
			inputs[2] = FromI32(e.action)
			inputs[3] = FromI32(e.mods)
			appMouseButtonCallback.Call(inputs)
		case APP_TOUCH:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appTouchCallback.windowNameBytes
			inputs[1] = FromI32(e.ix)
			inputs[2] = FromI32(e.iy)
			appTouchCallback.Call(inputs)
		case APP_FRAMEBUFFER_SIZE:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appFramebufferSizeCallback.windowNameBytes
			inputs[1] = FromI32(e.width)
			inputs[2] = FromI32(e.height)
			fmt.Printf("APP_FRAMEBUFFER_SIZE %d, %d\n", e.width, e.height)
			appFramebufferSizeCallback.Call(inputs)
		case APP_WINDOW_SIZE:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appWindowSizeCallback.windowNameBytes
			inputs[1] = FromI32(e.width)
			inputs[2] = FromI32(e.height)
			fmt.Printf("APP_WINDOW_SIZE %d, %d\n", e.width, e.height)
			appWindowSizeCallback.Call(inputs)
		case APP_WINDOW_POSITION:
			var inputs [][]byte = make([][]byte, 3)
			inputs[0] = appWindowPosCallback.windowNameBytes
			inputs[1] = FromI32(e.ix)
			inputs[2] = FromI32(e.iy)
			appWindowPosCallback.Call(inputs)
		}
	}
	polled = polled[0:0]
}
