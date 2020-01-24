// +build android

package main

import (
	//"encoding/binary"
	"fmt"
	. "github.com/SkycoinProject/cx/cx"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"log"
	//"golang.org/x/mobile/exp/app/debug"
	//"golang.org/x/mobile/exp/f32"
	//"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

/*var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32
)*/

type EventType int

const (
	APP_NONE EventType = iota
	APP_START
	APP_STOP
	APP_RESIZE
	APP_TOUCH
	APP_PAINT
)

type Event struct {
	a   app.App
	t   EventType
	msg string
	w   int
	h   int
	x   float32
	y   float32
}

var events []Event
var eventCount int

func pushEvent(a app.App, eventType EventType, msg string) int {
	e := eventCount
	if eventCount < len(events) {
	} else {
		events = append(events, Event{})
	}
	log.Printf("NEW_EVENT %s, INDEX %d, TYPE %d", msg, e, eventType)
	events[e] = Event{a: a, t: eventType, msg: msg}
	eventCount++
	return e
}

func startCXVM() {
	log.Printf("-----------------------------------------------------------------------------")
	log.Printf("CXDROID N")
	log.Printf("-----------------------------------------------------------------------------")
	//go eventLoop(a)
	args := parseManifest()
	log.Printf("ARGS %v", args)
	if CopyAssetsToFilesDir() {
		log.Printf("All assets were copied successfuly")
		CXSetWorkingDir(fmt.Sprintf("%s/", asset.GetFilesDir()))
		Run(args)
	}
}

func eventLoop(a app.App) {
	for e := range a.Events() {
		switch e := a.Filter(e).(type) {
		case lifecycle.Event:
			switch e.Crosses(lifecycle.StageVisible) {
			case lifecycle.CrossOn:
				_ = pushEvent(a, APP_START, "START")
				glctx, _ := e.DrawContext.(gl.Context)
				SetGLContext(glctx)
				go startCXVM()
			case lifecycle.CrossOff:
				_ = pushEvent(a, APP_STOP, "STOP")
			}
		case size.Event:
			i := pushEvent(a, APP_RESIZE, "RESIZE")
			events[i].w = e.WidthPx
			events[i].h = e.HeightPx
		case paint.Event:
			_ = pushEvent(a, APP_PAINT, "PAINT")
			if /*glctx == nil ||*/ e.External {
				continue
			}

			//onPaint(glctx, size.Event)
			//a.Publish()
			//a.Send(paint.Event{})
		case touch.Event:
			i := pushEvent(a, APP_TOUCH, "TOUCH")
			events[i].x = e.X
			events[i].y = e.Y
		}
	}
}

func parseManifest() []string {
	// TODO : extract args from manifest
	args := []string{
		"--stack-size=128M",
		"--heap-initial=800M",
		"-heap-max=800M",
		"--debug-profile=100",
		"lib/args.cx",
		"lib/json.cx",
		"cxfx/src/mat/math.cx",
		"cxfx/src/mat/v1d.cx",
		"cxfx/src/mat/v1f.cx",
		"cxfx/src/mat/v2f.cx",
		"cxfx/src/mat/v3f.cx",
		"cxfx/src/mat/v4f.cx",
		"cxfx/src/mat/q4f.cx",
		"cxfx/src/mat/m44f.cx",
		"cxfx/src/app/application.cx",
		"cxfx/src/app/event.cx",
		"cxfx/src/fps/profiler.cx",
		"cxfx/src/fps/framerate.cx",
		"cxfx/src/gfx/batch.cx",
		"cxfx/src/gfx/graphics.cx",
		"cxfx/src/gfx/state.cx",
		"cxfx/src/gfx/effect.cx",
		"cxfx/src/gfx/shader.cx",
		"cxfx/src/gfx/program.cx",
		"cxfx/src/gfx/mesh.cx",
		"cxfx/src/gfx/gltf.cx",
		"cxfx/src/gfx/model.cx",
		"cxfx/src/gfx/texture.cx",
		"cxfx/src/gfx/particle.cx",
		"cxfx/src/gfx/text.cx",
		"cxfx/src/gfx/target.cx",
		"cxfx/src/gfx/scissor.cx",
		"cxfx/src/snd/sounds.cx",
		"cxfx/src/snd/audio.cx",
		"cxfx/src/snd/voice.cx",
		"cxfx/src/gui/callback.cx",
		"cxfx/src/gui/layer.cx",
		"cxfx/src/gui/skin.cx",
		"cxfx/src/gui/scope.cx",
		"cxfx/src/gui/font.cx",
		"cxfx/src/gui/animation.cx",
		"cxfx/src/gui/control.cx",
		"cxfx/src/gui/label.cx",
		"cxfx/src/gui/picture.cx",
		"cxfx/src/gui/screen.cx",
		"cxfx/src/gui/interface.cx",
		"cxfx/src/gui/focus.cx",
		"cxfx/src/gui/splitter.cx",
		"cxfx/src/gui/window.cx",
		"cxfx/src/gui/keyboard.cx",
		"cxfx/src/gui/list.cx",
		"cxfx/src/gui/graph.cx",
		"cxfx/src/gui/lifter.cx",
		"cxfx/src/gui/scroller.cx",
		"cxfx/src/gui/binder.cx",
		"cxfx/src/gui/combo.cx",
		"cxfx/src/gam/camera.cx",
		"cxfx/src/gam/game.cx",
		"cxfx/src/phx/physics.cx",
		"cxfx/tutorials/0_colored_quad.cx",
		// "cxfx/tutorials/1_textured_quad.cx",
		// "cxfx/tutorials/2_text.cx",
		// "cxfx/tutorials/3_perspective.cx",
		// "cxfx/tutorials/4_camera.cx",
		// "cxfx/tutorials/5_batch.cx",
		// "cxfx/tutorials/6_model.cx",
		// "cxfx/tutorials/7_menu.cx",
		// "cxfx/tutorials/8_sound.cx",
		// "cxfx/tutorials/9_button.cx",
		// "cxfx/tutorials/10_dialog.cx",
		// "cxfx/game/skylight/src/skylight.cx",
		"++data=resources/",
		"++hints=resizable",
		//"++fps=30"
	}

	return args
}

func main() {
	app.Main(func(a app.App) {
		eventLoop(a)
	})
}

/*func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}
	glctx.Uniform4f(color, 0, green, 0, 1)

	glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	glctx.DisableVertexAttribArray(position)

	fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.4, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;
attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`*/
