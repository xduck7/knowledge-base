package factory

import (
	"fmt"
	"log"
)

type Vehicle interface {
	Drive() string
}

type Car struct{}
type Motorcycle struct{}
type Plane struct{}

func (Car) Drive() string        { return "Driving a car" }
func (Motorcycle) Drive() string { return "Riding a motorcycle" }
func (Plane) Drive() string      { return "Flying a plane" }

const (
	TypeCar        = "car"
	TypeMotorcycle = "motorcycle"
	TypePlane      = "plane"
)

func NewVehicle(t string) (Vehicle, error) {
	switch t {
	case TypeCar:
		return Car{}, nil
	case TypeMotorcycle:
		return Motorcycle{}, nil
	case TypePlane:
		return Plane{}, nil
	default:
		return nil, fmt.Errorf("unknown vehicle type: %s", t)
	}
}

func Example() {
	car, err := NewVehicle(TypeCar)
	if err != nil {
		log.Fatal(err)
	}
	bike, _ := NewVehicle(TypeMotorcycle)
	plane, _ := NewVehicle(TypePlane)

	fmt.Println(car.Drive())
	fmt.Println(bike.Drive())
	fmt.Println(plane.Drive())
}
