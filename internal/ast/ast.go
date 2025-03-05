package ast

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Node struct {
	Value    float64
	Operator string
	Left     *Node
	Right    *Node
}

type Parser struct {
	src string
	pos int
}

func Parse(expression string) (*Node, error) {
	expr := strings.ReplaceAll(expression, " ", "")
	if expr == "" {
		return nil, fmt.Errorf("error in expression")
	}
	p := &Parser{expr, 0}
	return p.parseExpr()
}

func (p *Parser) next() rune {
	if p.pos < len(p.src) {
		ch := rune(p.src[p.pos])
		p.pos++
		return ch
	}
	return 0
}

func (p *Parser) peek() rune {
	if p.pos < len(p.src) {
		return rune(p.src[p.pos])
	}
	return 0
}

func (p *Parser) parseExpr() (*Node, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '+' || ch == '-' {
			p.next()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			node = &Node{Operator: string(ch), Left: node, Right: right}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseTerm() (*Node, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '*' || ch == '/' {
			p.next()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			node = &Node{Operator: string(ch), Left: node, Right: right}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseFactor() (*Node, error) {
	if p.peek() == '(' {
		p.next()
		node, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.peek() != ')' {
			return nil, fmt.Errorf("error in expression")
		}
		p.next()
		return node, nil
	}
	start := p.pos
	for unicode.IsDigit(p.peek()) || p.peek() == '.' {
		p.next()
	}
	number := p.src[start:p.pos]
	if number == "" {
		return nil, fmt.Errorf("error in expression")
	}
	value, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return nil, fmt.Errorf("error in expression")
	}
	return &Node{Value: value}, nil
}
