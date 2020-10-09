package stdout

import (
	"github.com/fatih/color"
)

type Printer interface {
	Verbose()
	Print(atr color.Attribute, args ...interface{})
	Println(atr color.Attribute, args ...interface{})
	Printf(atr color.Attribute, template string, args ...interface{})
	VerbosePrint(atr color.Attribute, args ...interface{})
	VerbosePrintln(atr color.Attribute, args ...interface{})
	VerbosePrintf(atr color.Attribute, template string, args ...interface{})
}

type printer struct {
	verbose bool
}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Verbose() {
	p.verbose = true
}

func (p *printer) Print(atr color.Attribute, args ...interface{}) {
	_, _ = color.New(atr).Print(args...)
}

func (p *printer) Println(atr color.Attribute, args ...interface{}) {
	_, _ = color.New(atr).Println(args...)
}

func (p *printer) Printf(atr color.Attribute, template string, args ...interface{}) {
	_, _ = color.New(atr).Printf(template, args...)
}

func (p *printer) VerbosePrint(atr color.Attribute, args ...interface{}) {
	if p.verbose {
		p.Print(atr, args...)
	}
}

func (p *printer) VerbosePrintln(atr color.Attribute, args ...interface{}) {
	if p.verbose {
		p.Println(atr, args...)
	}
}

func (p *printer) VerbosePrintf(atr color.Attribute, template string, args ...interface{}) {
	if p.verbose {
		p.Printf(atr, template, args...)
	}
}
