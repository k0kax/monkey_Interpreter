package main

import (
	"fmt"
	"monkey_Interpreter/repl"
	"os"
	user2 "os/user"
)

func main() {
	//在 main 函数内部，首先使用 os/user 包中的 user2.Current() 函数获取当前用户的信息
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programmimg language!\n", user.Username)
	fmt.Printf("Feel freee to type in commamds\n")
	repl.Start(os.Stdin, os.Stdout)
}
