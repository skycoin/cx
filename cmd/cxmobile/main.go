// +build cxfx,android

package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"

	"github.com/skycoin/cx/cx/packages/cxfx"
)

func startCXVM() {
	log.Printf("-----------------------------------------------------------------------------")
	log.Printf("STARTING CXDROID N")
	log.Printf("-----------------------------------------------------------------------------")
	CXLogFile(true)
	if CopyAssetsToFilesDir() {
		filesDir := asset.GetFilesDir()
		CXSetWorkingDir(fmt.Sprintf("%s/", filesDir))
		args := parseRunCli(filesDir)
		log.Printf("ARGS %v", args)
		log.Printf("All assets were copied successfuly")
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
				cxfx.PushFramebufferSizeEvent(float64(e.WidthPx), float64(e.HeightPx))
				cxfx.PushWindowSizeEvent(float64(e.WidthPx), float64(e.HeightPx))
			case paint.Event:
				cxfx.PushEvent(cxfx.APP_PAINT)
				//if glctx == nil || e.External {
				//	continue
				//}
			case touch.Event:
				cxfxAction := cxfx.ACTION_RELEASE
				switch e.Type {
				case touch.TypeBegin:
					cxfxAction = cxfx.ACTION_PRESS
				case touch.TypeMove:
					cxfxAction = cxfx.ACTION_MOVE
				case touch.TypeEnd:
					cxfxAction = cxfx.ACTION_RELEASE
				}
				cxfx.PushMouseEvent(cxfx.APP_MOUSE, cxfxAction, int32(e.Sequence), int64(e.Sequence), 0, float64(e.X), float64(e.Y))
			}
		}
	}()

	startCXVM()
}

func parseRunCli(filesDir string) []string {
	cli, _ := CXReadFile("run.cli")
	cliString := string(cli)
	log.Printf("RUN_CLI_STRING %s\n", cliString)
	args := make([]string, 0, 0)
	args = append(args, "--cxpath=.")
	args = append(args, strings.Split(strings.Replace(cliString, "\r\n", "\n", -1), "\n")...)
	argCount := len(args)
	validArgs := make([]string, 0, 0)
	for i := 0; i < argCount; i++ {
		if len(args[i]) > 0 {
			validArgs = append(validArgs, args[i])
		}
	}
	return validArgs
}

func main() {
	app.Main(func(a app.App) {
		eventLoop(a)
	})
}
