package parser

import (
	"fmt"
	"strconv"

	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/token"
)

// Operator precedence constants.
const (
	_ int = iota
	Lowest
	Equals      // ==
	LessGreater // < or >
	Sum         // +
	Product     // *
	Prefix      // -x or !x
	Call        // someFunc(x)
	Index       // array[index]
)

var precedences = map[token.Type]int{
	token.Equals:    Equals,
	token.NotEquals: Equals,
	token.LT:        LessGreater,
	token.GT:        LessGreater,
	token.Plus:      Sum,
	token.Minus:     Sum,
	token.Slash:     Product,
	token.Asterisk:  Product,
	token.LParen:    Call,
	token.LBracket:  Index,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser is the object that takes a lexer, consumes all tokens from it and
// constructs a corresponding AST.
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New creates a Parser.
func New(l *lexer.Lexer) *Parser {
	p := Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[token.Type]prefixParseFn),
		infixParseFns:  make(map[token.Type]infixParseFn),
	}
	// prefix parse funcs:
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Num, p.parseIntegerLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Function, p.parseFunctionLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.LBracket, p.parseArrayLiteral)
	p.registerPrefix(token.LBrace, p.parseHashLiteral)
	p.registerPrefix(token.Macro, p.parseMacroLiteral)
	// infix parse funcs:
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Equals, p.parseInfixExpression)
	p.registerInfix(token.NotEquals, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)
	p.registerInfix(token.LBracket, p.parseIndexExpression)
	// read two tokens so that curToken and peekToken are both set:
	p.nextToken()
	p.nextToken()
	return &p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram is the main entry point of the parser.
func (p *Parser) ParseProgram() *ast.Program {
	program := ast.Program{
		Statements: []ast.Statement{},
	}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return &program
}

// Errors returns a list of errors encountered during parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := ast.LetStatement{
		Token: p.curToken,
	}
	if !p.expectPeek(token.Ident) {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	if !p.expectPeek(token.Assign) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(Lowest)
	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return &stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(Lowest)
	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return &stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(Lowest),
	}
	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return &stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := ast.IntegerLiteral{
		Token: p.curToken,
	}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return &lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := ast.ArrayLiteral{
		Token: p.curToken,
	}
	array.Elements = p.parseExpressionList(token.RBracket)
	return &array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}
	p.nextToken()
	exp.Index = p.parseExpression(Lowest)
	if !p.expectPeek(token.RBracket) {
		return nil
	}
	return &exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := ast.HashLiteral{
		Token: p.curToken,
	}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	for !p.peekTokenIs(token.RBrace) {
		p.nextToken()
		key := p.parseExpression(Lowest)
		if !p.expectPeek(token.Colon) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(Lowest)
		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBrace) && !p.expectPeek(token.Comma) {
			return nil
		}
	}
	if !p.expectPeek(token.RBrace) {
		return nil
	}
	return &hash
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	macro := ast.MacroLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LParen) {
		return nil
	}
	macro.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBrace) {
		return nil
	}
	macro.Body = p.parseBlockStatement()
	return &macro
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.True),
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(Prefix)
	return &expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(Lowest)
	if !p.expectPeek(token.RParen) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := ast.IfExpression{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LParen) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(Lowest)
	if !p.expectPeek(token.RParen) {
		return nil
	}
	if !p.expectPeek(token.LBrace) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.Else) {
		p.nextToken()
		if !p.expectPeek(token.LBrace) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return &expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LParen) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBrace) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return &lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, &ident)
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident := ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, &ident)
	}
	if !p.expectPeek(token.RParen) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBrace) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return &block
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return &expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseExpressionList(token.RParen)
	return &exp
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(Lowest))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(Lowest))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return Lowest
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return Lowest
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
