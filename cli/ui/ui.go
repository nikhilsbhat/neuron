package ui

import (
        "fmt"
        "github.com/fatih/color"
        "io"
)

type UiWriter struct {
        Writer io.Writer
}

type Color struct {
        colour color.Attribute
}

type NeuronUi struct {
        *UiWriter
}

var (
        Red     = Color{color.FgRed}
        Blue    = Color{color.FgBlue}
        Green   = Color{color.FgGreen}
        Cyan    = Color{color.FgCyan}
        Yellow  = Color{color.FgYellow}
        Magenta = Color{color.FgMagenta}
)

func (c *Color) makeitcolorful(msg string) string {
        return color.New(c.colour).SprintFunc()(msg)
}

func (n NeuronUi) NeuronSaysItsInfo(msg string) {
        n.neuronPrints(Green.makeitcolorful(msg))
}

func (n NeuronUi) NeuronSaysItsError(msg string) {
        n.neuronPrints(Red.makeitcolorful("Error: " + msg))
}

func (n NeuronUi) NeuronSaysItsWarn(msg string) {
        n.neuronPrints(Yellow.makeitcolorful(msg))
}

func Info(msg string) string {
        return Green.makeitcolorful(msg)
}

func Error(msg string) string {
        return Red.makeitcolorful(msg)
}

func Warn(msg string) string {
        return Yellow.makeitcolorful(msg)
}

func Debug(msg string) string {
        return Magenta.makeitcolorful(msg)
}

func (n NeuronUi) neuronPrints(msg string) {
        fmt.Fprint(n.Writer, msg)
        fmt.Fprintf(n.Writer, "\n")
}
