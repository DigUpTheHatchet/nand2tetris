package translator

import (
	"strings"
)

type CommandType int

const (
	C_ARITHMETIC CommandType = iota + 1 // EnumIndex = 1
	C_PUSH                              // EnumIndex = 2
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type Translator struct {
	parser     Parser
	codeWriter CodeWriter
}

func NewTranslator(inputFile string) *Translator {
	filePrefix := strings.TrimRight(inputFile, ".vm")
	p := NewParser(inputFile)
	cw := NewCodeWriter(filePrefix)
	t := &Translator{parser: *p, codeWriter: *cw}
	return t
}

func (t *Translator) Run() {
	for t.parser.hasMoreLines() {
		t.parser.advance()
		if t.parser.currentCommand == "" || strings.HasPrefix(t.parser.currentCommand, "//") {
			continue
		}
		cmdType := t.parser.commandType()

		if cmdType == C_PUSH || cmdType == C_POP {
			t.codeWriter.writePushPop(t.parser.currentCommand, cmdType, t.parser.arg1(), t.parser.arg2())
		} else {
			t.codeWriter.writeArithmetic(t.parser.arg1())
		}
	}

	t.codeWriter.writeInfiniteLoop()
}
