package main

import "fmt"

// Suppose you have a function that is supposed to pretty-print a temperature
// value with the degree (°) symbol.
func PrintTemperature(temp float64) {
	fmt.Printf("%f °???\n", temp)
}

// Already, we've bumped into a problem: Depending on what country you live in
// or even what occupation you have, you might use a different temperature
// unit. The three most common units of measurement for temperature are
// Fahrenheit, Celsius, and Kelvin. So let's write functions that can deal with
// each of those:
func PrintFahrenheit(temp float64) {
	fmt.Printf("%f °F\n", temp)
}

func PrintCelsius(temp float64) {
	fmt.Printf("%f °C\n", temp)
}

func PrintKelvin(temp float64) {
	fmt.Printf("%f °K\n", temp)
}

// We've run into another problem. Each temperature unit has a different scale
// of measurement. For example, water boils at 100 °C, 212 °F, and 373.15 °K.
// Similarly, water freezes at 0 °C, 32 °F, and 273.15 °K.
//
// What's to stop us from passing a Fahrenheit temperature into the function
// that prints Kelvin temperatures?
//
// The short answer: Nothing! This problem gets even worse when we have to
// start considering conversions from one temperature unit to another.
func main() {
	temp := 100.0

	PrintTemperature(temp)
	PrintFahrenheit(temp)
	PrintCelsius(temp)
	PrintKelvin(temp)
}
