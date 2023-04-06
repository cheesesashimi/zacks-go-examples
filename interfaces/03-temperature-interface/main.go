package main

import "fmt"

// An interface in Go is essentially a named collection of method signatures.
// So with that in mind, lets first create a Unit interface:
type Unit interface {
	Unit() string
}

// And just like with structs, you can embed one interface inside another. Here
// is our Temperature interface:
type Temperature interface {
	Kelvin() Kelvin
	Fahrenheit() Fahrenheit
	Celsius() Celsius

	Unit
}

type Fahrenheit float64
type Celsius float64
type Kelvin float64

// But there's a problem: Nothing actually *implements* the Temperature
// interface yet. So let's go ahead and do that! What each of these function
// declarations is doing is attaching a method onto its associated type. In
// other words, any basic type in Go can have helpers attached to it; not just
// structs.
func (f Fahrenheit) Fahrenheit() Fahrenheit {
	// Because this represents our underlying type, we just return it as-is.
	return f
}

func (f Fahrenheit) Celsius() Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func (f Fahrenheit) Kelvin() Kelvin {
	// The formula for converting Fahrenheit to Kelvin requires an intermediate
	// conversion to Celsius. The Celsius type knows how to convert itself to
	// Kelvin. So first, we call our Celsius method which performs the
	// intermediate conversion, then we call the Kelvin method to do the final
	// conversion.
	return f.Celsius().Kelvin()
}

func (f Fahrenheit) Unit() string {
	return "Fahrenheit"
}

// More on what this function does later. For now, notice that I did not
// include it in the Temperature interface definition. This was intentional.
func (f Fahrenheit) String() string {
	return fmt.Sprintf("%.2f °F", f.Fahrenheit())
}

// This method is only found on the Fahrenheit type intentionally. More on this
// later.
func (f Fahrenheit) IsThisCold() bool {
	return f <= 70.0
}

func (c Celsius) Fahrenheit() Fahrenheit {
	return Fahrenheit((c * 9 / 5) + 32)
}

func (c Celsius) Celsius() Celsius {
	return c
}

func (c Celsius) Kelvin() Kelvin {
	return Kelvin(c + 273.15)
}

func (c Celsius) Unit() string {
	return "Celsius"
}

// More on what this function does later. For now, notice that I did not
// include it in the Temperature interface definition. This was intentional.
func (c Celsius) String() string {
	return fmt.Sprintf("%.2f °C", c)
}

func (k Kelvin) Fahrenheit() Fahrenheit {
	// Similar to the above case, we must convert Kelvin to Celsius first before
	// we can convert to Fahrenheit. Since Celsius knows how to convert itself to
	// Fahrenheit, we can make use of its Fahrenheit method.
	return k.Celsius().Fahrenheit()
}

func (k Kelvin) Celsius() Celsius {
	return Celsius(k - 273.15)
}

func (k Kelvin) Kelvin() Kelvin {
	return k
}

func (k Kelvin) Unit() string {
	return "Kelvin"
}

// More on what this function does later. For now, notice that I did not
// include it in the Temperature interface definition. This was intentional.
func (k Kelvin) String() string {
	return fmt.Sprintf("%.2f °K", k)
}

// Now, we can verify each of our temperature types implements the Temperature
// interface. What this does is attempt to assign an instance of your type to a
// variable that expects the Temperature interface. Since the underscore is
// used for the variable name, the variable will be assigned then immediately
// discarded; taking up no actual memory at runtime. So if your type fails to
// implement the interface, this will fail to compile.
var _ Temperature = Kelvin(0.0)
var _ Temperature = Fahrenheit(0.0)
var _ Temperature = Celsius(0.0)

// Here are some units of distance which implement the Unit interface, but do
// not implement the Temperature interface because they are not temperatures.
// We'll reuse code which can work with units later. For brevity, I did not
// include conversion functions nor did I include an interface for them other
// than the Unit interface.
type Mile float64

func (m Mile) String() string {
	return fmt.Sprintf("%.2f miles", m)
}

func (m Mile) Unit() string {
	return "Mile"
}

type Kilometer float64

func (k Kilometer) String() string {
	return fmt.Sprintf("%.2f kilometers", k)
}

func (k Kilometer) Unit() string {
	return "Kilometer"
}

// https://en.wikipedia.org/wiki/Smoot
type Smoot float64

func (s Smoot) String() string {
	return fmt.Sprintf("%.2f smoots", s)
}

func (s Smoot) Unit() string {
	return "Smoot"
}

var _ Unit = Mile(0.0)
var _ Unit = Kilometer(0.0)
var _ Unit = Smoot(0.0)

// But why is this useful? I know this looks like way more code than my
// previous examples. And to some degree, it is! However, something you'll
// notice is that each temperature type knows what unit it is.
func getPrettyTemperature(temp Temperature) string {
	unit := ""
	switch temp.Unit() {
	case "Kelvin":
		unit = "K"
	case "Celsius":
		unit = "C"
	case "Fahrenheit":
		unit = "F"
	}

	return fmt.Sprintf("%.2f °%s", temp, unit)
}

// Not only does each temperature know what unit it is, it also knows how to
// convert to other units. With this in mind, a function that accepts a
// Temperature doesn't have to care what the underlying unit is. It can easily
// convert it to the unit that is required!
func printingTemps(temps []Temperature) {
	fmt.Println("using getPrettyTemperature():")
	for _, temp := range temps {
		fmt.Println("\tOriginal:", getPrettyTemperature(temp))
		fmt.Println("\tFahrenheit:", getPrettyTemperature(temp.Fahrenheit()))
		fmt.Println("\tCelsius:", getPrettyTemperature(temp.Celsius()))
		fmt.Println("\tKelvin:", getPrettyTemperature(temp.Kelvin()))
		fmt.Println("")
	}

	// But wait, you promised you'd go over what the String() method does! So,
	// this method holds special significance in Go. Why is that? In short,
	// String() string satisfies a built-in interface in Go called Stringer. What
	// does this mean?
	//
	// It means that any type that implements the String() interface can control
	// what its own string representation looks like. What does this mean in
	// practice? Remember our getPrettyTemperature function above? Technically,
	// we don't need it!
	fmt.Println("using each Temperature's stringer interface:")
	for _, temp := range temps {
		// fmt.Println() will call an objects' String() method and use its return
		// value if it satisfies the Stringer interface.
		fmt.Println("\tOriginal:", temp)
		fmt.Println("\tFahrenheit:", temp.Fahrenheit())
		fmt.Println("\tCelsius:", temp.Celsius())
		fmt.Println("\tKelvin:", temp.Kelvin())
		fmt.Println("")
	}

	// This also means we can do something like this:
	fmt.Println("using each Temperature's stringer interface in a formatted string:")
	for _, temp := range temps {
		fmt.Printf("\t")
		printTemperatureLine(temp)
	}

	// In short, a concrete type like Fahrenheit, Celsius, or Kelvin, can satisfy
	// more than one interface at a time.
}

// There is a downside to using interfaces like this. Sometimes, an underlying
// type might have a method on it that we want to call that's not specified in
// the interface. If that method signature is not in the interface, we can't
// directly use it if all we know about is that it satisfies a given interface.
//
// If you'll recall, only the Fahrenheit type has a IsThisCold() method
// attached to it. So given the interface, how can we call that method? The
// short answer: Type assertions!
func typeAssertions(temps []Temperature) {
	// A type assertion provides a way to access an interfaces underlying
	// concrete type. In this case, we accept a slice of Temperatures but only
	// a few of them is a Fahrenheit.
	for _, temp := range temps {
		// I should mention that type assertions occur during runtime, which means
		// that it is possible that an unchecked bad type assertion can cause a
		// panic.
		//
		// This syntax allows us to conditionally execute based upon the type
		// assertion. Not doing this conditionally will cause the program to panic
		// because the input does not match the underlying type.
		//
		// In this particular situation, if we don't match the underlying type, we
		// can convert to it quite easily. However, that won't always be the case.
		f, ok := temp.(Fahrenheit)
		if ok {
			fmt.Printf("Found a native Fahrenheit: ")
		} else {
			f = temp.Fahrenheit()
			fmt.Printf("Converted a non-native Fahrenheit from %s %s: ", temp.Unit(), temp)
		}

		printIsThisCold(f)
	}

	var fahrCount int
	var celCount int
	var kelCount int

	for _, temp := range temps {
		// You can also switch based upon type!
		switch tObj := temp.(type) {
		case Fahrenheit:
			printUnit(tObj)
			fahrCount += 1
		case Celsius:
			printUnit(tObj)
			celCount += 1
		case Kelvin:
			printUnit(tObj)
			kelCount += 1
		default:
			fmt.Println("Unknown temperature type!")
		}
	}

	fmt.Printf("Found %d Fahrenheits, %d Celsiuses, %d Kelvins\n", fahrCount, celCount, kelCount)
}

// Just like how you can create anonymous structs, you can also create
// anonymous interfaces! Remember how Fahrenheit has the IsThisCold() method
// attached to it?
func printIsThisCold(in interface{ IsThisCold() bool }) {
	fmt.Printf("IsThisCold (%s)? %v\n", in, in.IsThisCold())
}

// This function will only print a Unit.
func printUnit(in Unit) {
	fmt.Printf("Got a %s: %s\n", in.Unit(), in)
}

func printTemperatureLine(temp Temperature) {
	fmt.Printf("Original: %s \tCelsius: %s \tFahrenheit: %s \tKelvin: %s\n", temp, temp.Celsius(), temp.Fahrenheit(), temp.Kelvin())
}

func printUnits(units []Unit) {
	for _, unit := range units {
		printUnit(unit)
	}
}

func main() {
	// Now, let's declare multiple temperature types in a single slice! We can do
	// this becuase all of our temperature types implement the Temperature
	// interface.
	temps := []Temperature{
		Kelvin(0.0),
		Celsius(0.0),
		Fahrenheit(0.0),
		Fahrenheit(212.0),
		Fahrenheit(80.0),
		Celsius(70.0),
		Kelvin(300.0),
	}

	// Let's print our temps!
	fmt.Println("Printing temps:")
	printingTemps(temps)

	// Let's assert their types!
	fmt.Println("Assert their types:")
	typeAssertions(temps)

	fmt.Println("Units:")

	// Now lets do units!
	units := []Unit{
		Smoot(5.5),
		Mile(5.5),
		Kilometer(5.5),
	}

	// We cannot directly concatenate a list of Temperatures and a list of Units
	// since they are fundamentally different types; even if they have the Unit
	// interface in common!
	for _, temp := range temps {
		units = append(units, temp)
	}

	printUnits(units)
}
