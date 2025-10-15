package main

import "fmt"

type Shape interface {
	Area() float64
	Perimter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

func (r *Rectangle) Area() float64 {
	return r.width * r.height
}

func (r *Rectangle) Perimter() float64 {
	return 2 * (r.width + r.height)
}

func (c *Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (c *Circle) Perimter() float64 {
	return 2 * 3.14 * c.radius
}

type Person struct {
	Name string
	Age  uint8
}

type Employee struct {
	Person     Person
	EmployeeID int
}

func (e *Employee) PrintInfo() {
	fmt.Printf("EmployeeID:%d  Name:%s  age:%d\n", e.EmployeeID, e.Person.Name, e.Person.Age)
}

func main() {
	fmt.Println("---------------------------面积和周长---------------------------")
	r := Rectangle{width: 8, height: 6}
	c := Circle{radius: 7}
	fmt.Printf("矩形的面积是：%.4f\n", r.Area())
	fmt.Printf("矩形的周长是：%.4f\n", r.Perimter())
	fmt.Printf("圆形的面积是：%.4f\n", c.Area())
	fmt.Printf("圆形的周长是：%.4f\n", c.Perimter())

	fmt.Println("---------------------------组合结构体---------------------------")
	employee := Employee{EmployeeID: 1001, Person: Person{Name: "张三", Age: 32}}
	employee.PrintInfo()
}
