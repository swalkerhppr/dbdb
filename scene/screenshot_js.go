package scene

import (
	"bytes"
	"image/png"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
)


func takeScreenshot(filename string, screen *ebiten.Image) {
		buf := &bytes.Buffer{}
		png.Encode(buf, screen)

		global := js.Global()
		jsData := global.Get("Uint8Array").New(buf.Len())
		js.CopyBytesToJS(jsData, buf.Bytes())

		a := global.Get("document").Call("createElement", "a")
		blob := global.Get("Blob").New(
			[]any{jsData},
			map[string]any{"type": "image/png"},
		)
		a.Set("href", global.Get("URL").Call("createObjectURL", blob))
		a.Set("download", filename)
		a.Call("click")
}
