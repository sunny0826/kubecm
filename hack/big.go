package main

import "github.com/pterm/pterm"

func main() {
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("guoxudong.io")).Srender()
	pterm.Println(s) // Print BigLetters with the default CenterPrinter
}
