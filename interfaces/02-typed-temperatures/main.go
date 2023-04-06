package main

import (
	"fmt"
)

// To avoid this mistake, we can assign special meaning to the number by
// creating custom types:
type Fahrenheit float64
type Celsius float64
type Kelvin float64

// This lets you lean on the Go compiler to ensure that a function which
// expects a Fahrenheit value cannot inadvertantly accept a Celsius value.
// Here are our print functions:
func PrintFahrenheit(f Fahrenheit) {
	fmt.Printf("%f °F\n", f)
}

func PrintCelsius(c Celsius) {
	fmt.Printf("%f °C\n", c)
}

func PrintKelvin(k Kelvin) {
	fmt.Printf("%f °K\n", k)
}

// Now, attempting to write code like:
//
// var f Fahrenheit = 212.0
// PrintCelsius(f)
//
// Will result in a compile error! So far so good, right? Yes!
//
// Now let's say that we want to convert between the various units. This is
// where the type system can really help us it makes it more difficult to
// accidentally pass a Celsius temperature into a function that expects
// Fahrenheit.
func FahrenheitToCelsius(f Fahrenheit) Celsius {
	return Celsius((float64(f) - 32) * 5 / 9)
}

func CelsiusToFahrenheit(c Celsius) Fahrenheit {
	return Fahrenheit((float64(c) * 9 / 5) + 32)
}

func CelsiusToKelvin(c Celsius) Kelvin {
	return Kelvin(float64(c) + 273.15)
}

func KelvinToCelsius(k Kelvin) Celsius {
	return Celsius(float64(k) - 273.15)
}

func FahrenheitToKelvin(f Fahrenheit) Kelvin {
	// Note: The commonly-accepted formula uses Celsius as an intermediate
	// conversion before converting to Kelvin. So we just reuse those conversion
	// functions here.
	return CelsiusToKelvin(FahrenheitToCelsius(f))
}

func KelvinToFahrenheit(k Kelvin) Fahrenheit {
	// Note: The commonly-accepted formula uses Celsius as an intermediate
	// conversion before converting to Kelvin. So we just reuse those conversion
	// functions here.
	return CelsiusToFahrenheit(KelvinToCelsius(k))
}

func main() {
	c := Celsius(100.0)

	// Uncomment me and try to build. You should get a compile error because
	// you're passing this into a function which accepts a different type.
	// PrintFahrenheit(c)

	f := Fahrenheit(212.0)
	k := Kelvin(373.15)

	PrintCelsius(c)
	PrintFahrenheit(f)
	PrintKelvin(k)

	PrintCelsius(FahrenheitToCelsius(f))
	PrintFahrenheit(CelsiusToFahrenheit(c))
	PrintKelvin(FahrenheitToKelvin(f))

	// However, using multiple types to represent different temperature units
	// still has one shortcoming: We can still cast from one type to another
	// without actually converting the underlying meaning:
	PrintCelsius(Celsius(f))

	// Another problem that we'll run into is adding support for an additional
	// temperature unit. For the sake of brevity (and a lack of creativity), I
	// won't actually introduce a new temperature unit. But let's imagine that I
	// did. We'll need, at a minimum:
	//
	// 1. A new type.
	// 2. Three conversion functions that will convert from our new unit to Kelvin,
	// Celsius, and Fahrenheit.
	// 3. Three conversion functions that will convert from Kelvin, Celsius, and
	// Fahrenheit to our new unit.
	//
	// Note: I deliberately chose temperatures for this example since there are
	// only three commonly-accepted scales. However, there are multiple measures of
	// distance such as inches, centimeters, kilometers, miles, yards, etc. This
	// quickly turns into a combinatoric nightmare.
	//
	// While we can't completely eliminate the combinatoric nightmare, we do have
	// ways of ensuring that code which uses these units does not have to care
	// about what the underlying unit is; just that it can convert it to the needed
	// unit.
	//
	// This is where interfaces come into play.
}
