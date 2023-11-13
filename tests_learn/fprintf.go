package main

import (
	"fmt"
	"os"
)

func main() {
	name := "Alice"
	age := 25

	// 将格式化的字符串输出到标准输出
	fmt.Fprintf(os.Stdout, "Hello, %s! You are %d years old.\n", name, age)

	// 将格式化的字符串写入文件
	file, _ := os.Create("tests_learn/output.txt")
	defer file.Close()
	fmt.Fprintf(file, "Hello, %s! You are %d years old.\n", name, age)
}
