package main

import "fmt"

func main() {
	a := 5
	addTen(&a)
	fmt.Println(a)

	b := []int{1, 2, 3}
	multiTwo(&b)
	fmt.Println(b)
}

func addTen(x *int) {
	*x += 10
}

func multiTwo(p *[]int) {
	fmt.Println("原切片内容：", *p)
	// for循环 适合需要修改元素或依赖索引的场景
	for i := 0; i < len(*p); i++ {
		(*p)[i] *= 2
	}
	fmt.Println("通过for循环改变元素，切片内容改变为：", *p)

	// for range 适合只需要读取元素的场景，修改元素不会影响原切片
	for _, v := range *p {
		v += 10
	}
	fmt.Println("通过for range改变元素，原切不会变", *p)
}
