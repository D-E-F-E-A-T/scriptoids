package cli

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type Display struct {
	NoColor   bool
	NoSymbols bool
}

func (d *Display) getPrefix(symbol string, label string) string {
	if d.NoSymbols {
		return label
	} else {
		return symbol
	}
}

func (d *Display) Info(msg string, a ...interface{}) {
	prefix := d.getPrefix(".", "Info:")
	formattedMsg := fmt.Sprintf(msg, a...)

	if d.NoColor {
		fmt.Println(prefix, formattedMsg)
	} else {
		fmt.Println(Bold(Black(prefix)), Bold(Black(formattedMsg)))
	}
}

func (d *Display) Success(msg string, a ...interface{}) {
	prefix := d.getPrefix("✔", "Success:")
	formattedMsg := fmt.Sprintf(msg, a...)

	if d.NoColor {
		fmt.Println(prefix, formattedMsg)
	} else {
		fmt.Println(Green(prefix), formattedMsg)
	}
}

func (d *Display) Failure(msg string, a ...interface{}) {
	prefix := d.getPrefix("✘", "Error:")
	formattedMsg := fmt.Sprintf(msg, a...)

	if d.NoColor {
		fmt.Println(prefix, formattedMsg)
	} else {
		fmt.Println(Red(prefix), formattedMsg)
	}
}
