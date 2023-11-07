package main

import "fmt"

// Car represents the complex object being built
type Car struct {
	color         string
	engineType    string
	hasSunroof    bool
	hasNavigation bool
}

// CarBuilder provides an interface for constructing the parts of the car
type CarBuilder interface {
	SetColor(color string) CarBuilder
	SetEngineType(engineType string) CarBuilder
	SetSunroof(hasSunroof bool) CarBuilder
	SetNavigation(hasNavigation bool) CarBuilder
	Build() *Car
}

// carBuilder implements the CarBuilder interface
type carBuilder struct {
	car *Car
}

// NewCarBuilder creates a new CarBuilder
func NewCarBuilder() CarBuilder {
	return &carBuilder{
		car: &Car{}, // Initialize the car attribute
	}
}

func (cb *carBuilder) SetColor(color string) CarBuilder {
	cb.car.color = color
	return cb
}

func (cb *carBuilder) SetEngineType(engineType string) CarBuilder {
	cb.car.engineType = engineType
	return cb
}

func (cb *carBuilder) SetSunroof(hasSunroof bool) CarBuilder {
	cb.car.hasSunroof = hasSunroof
	return cb
}

func (cb *carBuilder) SetNavigation(hasNavigation bool) CarBuilder {
	cb.car.hasNavigation = hasNavigation
	return cb
}

func (cb *carBuilder) Build() *Car {
	return cb.car
}

// Director provides an interface to build cars
type Director struct {
	builder CarBuilder
}

func (d *Director) ConstructCar(color, engineType string, hasSunroof, hasNavigation bool) *Car {
	d.builder.SetColor(color).
		SetEngineType(engineType).
		SetSunroof(hasSunroof).
		SetNavigation(hasNavigation)

	return d.builder.Build()
}

func main() {
	// Create a new car builder.
	builder := NewCarBuilder()

	// Create a car with the director.
	director := &Director{builder: builder}
	myCar := director.ConstructCar("blue", "electric", true, true)

	// Use the car object with the chosen configuration.
	fmt.Println(myCar)
}
