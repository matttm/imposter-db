package main

import "github.com/charmbracelet/huh"

func PromptSelection(title string, items []string) []string {
	selections := []string{}
	options := huh.NewOptions(items...)
	w := huh.NewMultiSelect[string]().
		Options(options...).
		Title(title).
		Limit(1).
		Value(&selections)
	if err := w.Run(); err != nil {
		panic(err)
	}
	return selections
}
