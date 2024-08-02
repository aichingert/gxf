package main

import (
	"encoding/json"
	"github.com/aichingert/dxf"
	color "github.com/aichingert/dxf/pkg/colors"
	"syscall/js"
)

func Parse() js.Func {
	jsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return "Invalid number of arguments passed"
		}

		buffer := make([]byte, args[1].Get("byteLength").Int())
		js.CopyBytesToGo(buffer, args[1])

		drawing, err := dxf.Parse(args[0].String(), buffer)

		if err != nil {
			return "Error parsing dxf file"
		}

		var value map[string]interface{}
		buffer, err = json.Marshal(drawing)

		if err != nil {
			return "Error marshaling drawing"
		}

		err = json.Unmarshal(buffer, &value)
		if err != nil {
			return "Error unmarshaling drawing"
		}

		return value
	})

	return jsFunc
}

func DxfColorToRgb() js.Func {
	jsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid number of arguments passed"
		}

		colorCode := args[0].Int()
		result := js.Global().Get("Uint8Array").New(3)
		js.CopyBytesToJS(result, color.DxfColorToRGB[colorCode])
		return result
	})

	return jsFunc
}

func main() {
	js.Global().Set("parseDxf", Parse())
	js.Global().Set("dxfColorToRgb", DxfColorToRgb())
	<-make(chan struct{})
}
