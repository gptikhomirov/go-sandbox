package main

import "fmt"

type Device struct {
	Name   string
	Type   string
	Active bool
}

func (d Device) Info() string {
	return fmt.Sprintf("name: %s type:%s active:%s", d.Name, d.Type, d.Active)
}

func (d *Device) Activate() {
	d.Active = true
}

func (d *Device) Deactivate() {
	d.Active = false
}

func main() {
	d1 := Device{Name: "macbook", Type: "laptop", Active: false}
	fmt.Println(d1)

	d1.Activate()
	fmt.Println(d1)
	d1.Deactivate()
	fmt.Println(d1)
}
