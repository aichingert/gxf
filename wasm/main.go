package main

import (
    "fmt"
    "syscall/js"
    "encoding/json"

    "github.com/aichingert/gxf"
)

func Parse() js.Func {
    jsFunc := js.FuncOf(func(this js.Value, args[]js.Value) any {
        if len(args) != 1 {
            return fmt.Sprintf("Invalid amount of arguments passed: expected 1 got %d", len(args))
        }

        buffer := make([]byte, args[0].Get("byteLength").Int())
        js.CopyBytesToGo(buffer, args[0])

        plan, err := gxf.Parse(buffer)

        if err != nil {
            return "Error: parsing dxf file"
        }

        var value map[string]interface{}
        buffer, err = json.Marshal(plan)

        if err != nil {
            return "Error: marshaling plan"
        }

        err = json.Unmarshal(buffer, &value)
        if err != nil {
            return "Error unmarshaling plan"
        }

        return value
    })

    return jsFunc
}

func main() {
    js.Global().Set("parse", Parse())

    // NOTE: used for go wasm runtime to not stop
    <-make(chan struct{})
}
