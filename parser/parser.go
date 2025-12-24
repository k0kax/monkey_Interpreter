package parser

//语法分析器
import (
	"fmt"
	"monkey_Interpreter/ast"
	"monkey_Interpreter/lexer"
	"monkey_Interpreter/token"
	"strconv"
)

/*
普拉特语法分析器的主要思想是将解析函数（普拉特称为语义代码）与词法单元类型
相关联。每当遇到某个词法单元类型时，都会调用相关联的解析函数来解析对应的表达
式，最后返回生成的AST节点。每个词法单元类型最多可以关联两个解析函数，这取决于
词法单元的位置，是位于前缀位置还是中缀位置。
*/

// 前缀解析函数与中缀解析函数
type (
	prefixParseFn func() ast.Expression               //前缀解析函数，左侧为空
	infixParseFn  func(ast.Expression) ast.Expression //中缀解析函数，接受的参数时中缀表达式左边的内容
)

// 语法分析器结构
type Parser struct {
	l *lexer.Lexer //指向词法分析器实例的指针

	curToken  token.Token //当前词法单元
	peekToken token.Token //当前词法单元的下一位

	errors []string //错误集合

	//添加两个解析函数的映射
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// 解析优先级
const (
	_ int = iota
	LOWEST
	AND         //&&与
	EQUALS      //==
	LESSGREATER //> or <
	SUM         //+
	PRODUCT     //*
	PREFIX      //-X or !x
	CALL        //函数func
)

// 优先级表
var precedences = map[token.TokenType]int{
	token.AND:      AND, //逻辑与
	token.OR:       AND, //逻辑或
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER, //<
	token.GT:       LESSGREATER, //>
	token.PLUS:     SUM,         //+
	token.MINUS:    SUM,         //-
	token.SLASH:    PRODUCT,     //除/
	token.ASTERISK: PRODUCT,     //*
	token.LPAREN:   CALL,
}

// 实例化语法分析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors: []string{},
	} //语法分析器实例

	//读取两个词法单元，以设置curToken和peekToken
	p.nextToken()
	p.nextToken()

	//关联解析函数
	//前缀解析函数
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) //初始化映射
	p.registerPrefix(token.IDENT, p.parseIdentifier)           //注册ident标识符相关的解析函数（parseIdentifier）
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         //注册integer整形相关的解析函数（parseIntegerLiteral）
	p.registerPrefix(token.BANG, p.parsePrefixExpression)      //注册！非的解析函数
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     //注册-负号的解析函数

	//中缀解析函数
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)     //注册加号+的解析函数
	p.registerInfix(token.MINUS, p.parseInfixExpression)    //注册-负号的解析函数
	p.registerInfix(token.SLASH, p.parseInfixExpression)    //注册/的解析函数
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) //注册*的解析函数
	p.registerInfix(token.EQ, p.parseInfixExpression)       //注册==的解析函数
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)   //注册!=的解析函数
	p.registerInfix(token.LT, p.parseInfixExpression)       //注册<的解析函数
	p.registerInfix(token.GT, p.parseInfixExpression)       //注册>的解析函数

	//布尔型字面量 前缀表达式
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	//分组表达式() 前缀表达式
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	//解析if语句 前缀表达式
	p.registerPrefix(token.IF, p.parseIfExpression)

	//解析函数字面量 前缀表达式
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	//解析调用函数表达式
	p.registerInfix(token.LPAREN, p.parseCallExpression) //以左括号为中心，注册一个中缀表达式

	//解析字符串表达式
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	//解析逻辑运算符
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	return p
}

// 获取下一个词法单元 前移curToken和peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken() //递归
}

// 入口点
// 解析程序AST的根节点
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}              //声明构造ast根节点program
	program.Statements = []ast.Statement{} //语句接口切片集

	for p.curToken.Type != token.EOF { //碰到词法法单元Token EOF文件结尾 表示已将遍历完终止
		stmt := p.parseStatement() //解析具体句子
		if stmt != nil {
			program.Statements = append(program.Statements, stmt) //不断解析语句，并且存到statements切片中
		}
		p.nextToken() //下移
	}

	return program
}

// 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// 解析let语句 以为例let x=5;
// 此时：curtoken=let peektoken=x
func (p *Parser) parseLetStatement() *ast.LetStatement {
	//1.初始化 LetStatement 节点：把当前 curToken（就是 LET 类型的 "let"）绑定到节点
	stmt := &ast.LetStatement{Token: p.curToken}
	//运行后：stmt.token=let curtoken=let peektoken=x

	//2.检测标识符
	//检测下一个token(也就是peektoken)不是标识符indent，不是则退出（此处检测到为x是标识符，不退）,是则peektoken、curtoken后移一位
	if !p.expectPeek(token.IDENT) { //执行expectPeek(),检测到peektoken.type=IDENT,不执行{}内容，peektoken、curtoken都后移一位
		return nil
	}
	//运行后：stmt.token=let curtoken=x peektoken = "="

	//3.
	//将当前词法单元作为标识符的 Token 字段，并将其字面值literal作为标识符indent的值value赋给 stmt.Name
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	//运行后：stmt.token=let curtoken=x peektoken = "=" stmt.Name=&{Token: IDENT("x"), Value: "x"}

	//4.检测等号=
	//检查下一个token（peektoken）,,是则继续，不是则退出，peektoken、curtoken后移一位
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	//运行后：stmt.token=let curtoken="=" peektoken = "5"  stmt.Name=&{Token: IDENT("x"), Value: "x"}

	//5.TODO: 跳过对表达式的处理parseExpression()
	//运行后：stmt.token=let curtoken="5" peektoken = ";"  stmt.Name=&{Token: IDENT("x"), Value: "x"} tmt.Value = &IntegerLiteral {Token: INT ("5"), Value: 5}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	//6.检测分号（;）     处理语句末尾的分号（;）
	//检测当前token（curtoken）是否是分号（;）
	if !p.curTokenIs(token.SEMICOLON) { //是，则不需要移动
		p.nextToken() //不是，则peektoken、curtoken后移一位，直接解析下一句
	}
	//运行后：stmt.token=let curtoken=";" peektoken = ""  stmt.Name=&{Token: IDENT("x"), Value: "x"} stmt.

	//7.直接返回stmt
	return stmt
	//LetStatement {Token: LET ("let"), Name: Identifier ("x"), Value: IntegerLiteral (5)}

}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	//TODO:跳过对表达式的处理，直接遇到分号
	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 辅助断言函数
// 当前token判断
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 下一个token判断
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 用于判断下一个词法单元的类型是否与给定的类型匹配，并移动到下一个词法单元
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken() //后移一位 curtoken变成peektoken,peektoken变成peektoken的下一位
		return true
	} else {
		p.peekErrors(t)
		return false
	}
}

// 错误检测
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be \"%s\",got=%s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// 辅助函数
// 向前缀、中缀解析函数映射表添加数据
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 根据词法单元类型返回优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 返回当前词法单元的优先级
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 解析表达式 普卡特语法解析器核心 1+2*3
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 第一步：找当前token对应的「前缀解析函数」
	prefix := p.prefixParseFns[p.curToken.Type]
	//当前token：1，注册前缀解析函数prefix也就是parseIntegerLiteral()
	// 第二步：如果没有对应的前缀解析函数，报错并返回nil
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// 第三步：执行前缀解析函数，得到表达式的「左部分」
	leftExp := prefix()
	// 执行parseIntegerLiteral()，leftExp = &ast.IntegerLiteral{Value: 1}
	// 此时 token 状态：curToken=INT(1)，peekToken=PLUS(+)

	// 第四步：循环解析「中缀表达式」（核心循环，处理优先级）左结合
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 循环条件（两个都满足才继续）：
		// 1. !p.peekTokenIs(token.SEMICOLON)：下一个token不是分号（说明表达式还没结束）；
		// 2. precedence < p.peekPrecedence()：当前表达式的优先级 < 下一个token的优先级（保证先解析高优先级的）；

		//下一个token为+，进入循环

		// 找到下一个token对应的中缀解析函数
		infix := p.infixParseFns[p.peekToken.Type]
		// p.peekToken.Type = PLUS(+) → 拿到 parseInfixExpression 函数
		// infix = parseInfixExpression

		// 如果没有中缀解析函数，说明表达式结束，返回当前的左部分
		if infix == nil {
			return leftExp
		}
		// infix不为空跳过

		// 移动token
		p.nextToken()
		// 执行后 token 状态：curToken=PLUS(+)，peekToken=INT(2)

		// 执行中缀解析函数（核心递归），更新左部分为「完整的中缀表达式」
		leftExp = infix(leftExp)
		// 传入 leftExp=1，执行 parseInfixExpression(1)
	}

	// 第五步：返回最终解析好的表达式
	return leftExp //调用该解析函数
}

// 解析表达式语句expression
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken} //构建AST
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { //检查分号（有没有分号都可以运行）
		p.nextToken()
	}
	return stmt
}

// 解析函数：解析标识符indent
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 解析函数：解析整形Int
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// 将格式化错误信息添加到语法分析器的errors字段
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 中缀解析函数
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// 1. 初始化中缀表达式节点（+）
	expression := &ast.InfixExpression{
		Token:    p.curToken,         //+
		Operator: p.curToken.Literal, //“+”
		Left:     left,               //1
	}

	precedence := p.curPrecedence()                  //记录当前词法单元+的优先级 SUM
	p.nextToken()                                    //curToken=INT(2)，peekToken=ASTERISK(*)
	expression.Right = p.parseExpression(precedence) //递归调用 parseExpression(SUM) 解析右边的 2*3（重点！）

	return expression
}

// 布尔值解析函数
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// 解析分组表达式()
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) { //遇到)就结束
		return nil
	}

	return exp
}

// 解析if语句
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	//if部分 首选
	//识别条件
	if !p.expectPeek(token.LPAREN) { //下一个token是左括号(，移动token，继续，进入条件识别；不是则直接退出
		return nil
	}

	p.nextToken()                                    //token移动开始识别条件
	expression.Condition = p.parseExpression(LOWEST) //低权限识别条件，并写入条件中

	if !p.expectPeek(token.RPAREN) { //下一个token是右括号)，移动token，继续，表明条件结束
		return nil
	}

	//识别结果
	if !p.expectPeek(token.LBRACE) { //下一个token是左大括号{，移动token，继续，开始进行结果识别；不是则直接退出
		return nil
	}

	//写入首选结果
	expression.Consequence = p.parseBlockStatement()

	//中间选项
	for p.peekTokenIs(token.ELIF) {
		p.nextToken()
		elif := p.parseElIfExpression()
		expression.Alternatives = append(expression.Alternatives, elif)
	}

	//else部分 最后选项
	if p.peekTokenIs(token.ELSE) { //识别下一个token是else,进入else的判断；不是则直接退出
		p.nextToken() //后移token 匹配到真正的

		if !p.expectPeek(token.LBRACE) { //识别下一个token是左大括号{
			return nil
		}

		expression.LastAlternative = p.parseBlockStatement() //写入可替换结果
	}

	return expression
}

// 解析ELIF 和解析IF差不多
func (p *Parser) parseElIfExpression() *ast.ElIfExpression {
	expression := &ast.ElIfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) { //下一个token是左括号(，移动token，继续，进入条件识别；不是则直接退出
		return nil
	}

	p.nextToken()                                    //token移动开始识别条件
	expression.Condition = p.parseExpression(LOWEST) //低权限识别条件，并写入条件中

	if !p.expectPeek(token.RPAREN) { //下一个token是右括号)，移动token，继续，表明条件结束
		return nil
	}

	//识别结果
	if !p.expectPeek(token.LBRACE) { //下一个token是左大括号{，移动token，继续，开始进行结果识别；不是则直接退出
		return nil
	}

	//写入首选结果
	expression.Consequence = p.parseBlockStatement()

	return expression
}

// 解析块
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// 解析函数字面量
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	//开始解析参数
	if !p.expectPeek(token.LPAREN) { //识别左括号(
		return nil
	}
	//解析参数
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) { //识别右括号)
		return nil
	}
	//解析参数结束

	//解析函数体
	lit.Body = p.parseBlockStatement()

	return lit
}

// 解析函数字面量里的参数
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) { //如果识别到下一个token是右括号则，解析参数结束。
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	//识别第一个参数
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) { //遇到的下一个token是逗号,说明有多个参数，开始识别多个参数
		p.nextToken()
		p.nextToken() //跳过逗号，将新参数写入参数切片
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	//识别所有参数结束

	if !p.expectPeek(token.RPAREN) { //下一个token是右括号)，则下移token，返回所有参数
		return nil
	}

	return identifiers
}

// 解析函数：调用函数表达式
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

// 解析函数调用的各个参数
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) { //下一个token是右括号，说明解析完所有参数
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) { //下一个token是右括号，说明解析完所有参数，并后移一位
		return nil
	}

	return args
}

// 解析函数 字符串
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}
