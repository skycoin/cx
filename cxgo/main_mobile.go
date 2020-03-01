// +build cxfx,android

package main

import (
	"fmt"
	. "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxfx"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
	"log"
)

func startCXVM() {
	filesDir := asset.GetFilesDir()
	args := parseManifest(filesDir)
	log.Printf("-----------------------------------------------------------------------------")
	log.Printf("STARTING CXDROID N")
	log.Printf("ARGS %v", args)
	log.Printf("-----------------------------------------------------------------------------")
	CXLogFile(true)
	if CopyAssetsToFilesDir() {
		log.Printf("All assets were copied successfuly")
		CXSetWorkingDir(fmt.Sprintf("%s/", filesDir))
		Run(args)
	}
	log.Printf("-----------------------------------------------------------------------------")
	log.Printf("STOPPING CXDROID N")
	log.Printf("-----------------------------------------------------------------------------")
}

func eventLoop(a app.App) {
	go func() {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageFocused) {
				case lifecycle.CrossOn:
					var glctx gl.Context
					glctx, _ = e.DrawContext.(gl.Context)
					cxfx.SetGLContext(glctx)
					cxfx.SetGOApp(a)
					cxfx.PushEvent(cxfx.APP_START)
				case lifecycle.CrossOff:
					cxfx.PushEvent(cxfx.APP_STOP)
				}
			case size.Event:
				cxfx.PushFramebufferSizeEvent(int32(e.WidthPx), int32(e.HeightPx))
				cxfx.PushWindowSizeEvent(int32(e.WidthPx), int32(e.HeightPx))
			case paint.Event:
				cxfx.PushEvent(cxfx.APP_PAINT)
				//if glctx == nil || e.External {
				//	continue
				//}
			case touch.Event:
				cxfx.PushTouchEvent(int32(e.Type), int32(e.X), int32(e.Y))
			}
		}
	}()

	startCXVM()
}

func parseManifest(filesDir string) []string {
	// TODO : extract args from manifest
	args := []string{
		"--stack-size=128M",
		"--heap-initial=800M",
		"-heap-max=800M",
		//"--debug-profile=100",
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
		"cxfx/src/phx/physics.cx",
		//"cxfx/tutorials/0_colored_quad.cx",
		//"cxfx/tutorials/1_textured_quad.cx",
		//"cxfx/tutorials/2_text.cx",
		//"cxfx/tutorials/3_perspective.cx",
		//"cxfx/tutorials/4_camera.cx",
		//"cxfx/tutorials/5_batch.cx",
		"cxfx/tutorials/6_model.cx",
		//"cxfx/tutorials/7_menu.cx",
		//"cxfx/tutorials/8_sound.cx",
		// "cxfx/tutorials/9_button.cx",
		// "cxfx/tutorials/10_dialog.cx",
		//"cxfx/games/skylight/src/skylight.cx",
		"++data=/cxfx/resources/",
		"++hints=resizable",
		"++glVersion=gles31",
		//"++fps=30"
	}

	return args
}

func main() {
	app.Main(func(a app.App) {
		eventLoop(a)
	})
}
