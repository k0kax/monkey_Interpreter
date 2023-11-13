// repl/repl.go
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey_Interpreter/lexer"
	"monkey_Interpreter/token"
)

const PROMPT = ">>" //表示 REPL 提示符

// 函数接受输入流 in 和输出流 out 作为参数
func Start(in io.Reader, out io.Writer) {

	//在函数内部，创建了一个 bufio.Scanner 对象 scanner，用于从输入流中读取用户输入
	scanner := bufio.NewScanner(in)

	for {
		//使用 fmt.Fprintf 函数将提示符输出到输出流 out
		fmt.Fprintf(out, PROMPT)

		//调用 scanner.Scan() 方法来等待用户输入，并返回一个布尔值表示是否成功读取到输入
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		//如果成功读取到输入，将用户输入的文本保存在变量 line 中
		line := scanner.Text()

		//创建了一个 lexer.Lexer 对象 l，并使用用户输入的文本作为输入来初始化该对象
		l := lexer.New(line)

		//进入一个循环，每次迭代中调用 l.NextToken() 方法来获取下一个词法单元，并将其输出到输出流 out
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
