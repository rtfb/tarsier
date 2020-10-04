package parser

import (
	"github.com/rtfb/tarsier/ast"
	"github.com/rtfb/tarsier/lexer"
	"github.com/rtfb/tarsier/token"
)

// Parser is the object that takes a lexer, reads all tokens from it and
// constructs a corresponding AST.
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

// New creates a Parser.
func New(l *lexer.Lexer) *Parser {
	p := Parser{l: l}
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	default:
		return nil
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
	// TODO: we're skipping the expression until we encounter a semicolon:
	for !p.curTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return &stmt
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
	return false
}
