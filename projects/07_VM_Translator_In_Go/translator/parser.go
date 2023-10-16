package translator

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	scanner        bufio.Scanner
	currentCommand string
	fields         []string
	moreLines      bool
}

func NewParser(inputFile string) *Parser {
	if !strings.HasSuffix(inputFile, ".vm") {
		log.Fatal(errors.New("InputFile: (%s) must have .vm extension.."))
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	p := &Parser{moreLines: true, scanner: *bufio.NewScanner(file)}
	return p
}

func (p *Parser) hasMoreLines() bool {
	return p.moreLines
}

func (p *Parser) advance() {
	p.moreLines = p.scanner.Scan()
	trimmed := strings.TrimSpace(p.scanner.Text())
	p.currentCommand = trimmed
	p.fields = strings.Fields(trimmed)

	if err := p.scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (p *Parser) commandType() CommandType {

	var cmdType CommandType

	switch p.fields[0] {
	case "push":
		cmdType = C_PUSH
	case "pop":
		cmdType = C_POP
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		cmdType = C_ARITHMETIC
	default:
		// TODO: Raise an error?
		cmdType = C_ARITHMETIC
	}

	return cmdType
}

func (p *Parser) arg1() string {
	// add -> arg1 = 'add'
	// lt -> arg1 = 'lt'

	if p.commandType() == C_ARITHMETIC {
		return p.fields[0]
	}

	// push local 2 -> arg1 = 'local'
	return p.fields[1]
}

func (p *Parser) arg2() int {

	// push local 2 -> arg2 = '2'
	i64, _ := strconv.ParseInt(p.fields[2], 10, 0)
	return int(i64)
}
