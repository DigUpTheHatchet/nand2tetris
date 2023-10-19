package translator

import (
	"log"
	"os"
	"path"
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
	parsers    []Parser
	codeWriter CodeWriter
}

func NewTranslator(input string) *Translator {

	vmFiles := []string{}
	outputFileName := "testingfiles/" + strings.TrimSuffix(input, ".vm") + ".asm"

	// Single ".vm" file was passed as input
	if strings.HasSuffix(input, ".vm") {
		vmFiles = append(vmFiles, input)
	} else {
		// Directory was passed
		files, err := os.ReadDir("testingfiles/" + input)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".vm") {
				vmFiles = append(vmFiles, path.Join(input, file.Name()))
			}
		}
	}

	cw := NewCodeWriter(outputFileName)
	parsers := []Parser{}
	for _, vmFile := range vmFiles {
		parsers = append(parsers, *NewParser(vmFile))
	}

	t := &Translator{parsers: parsers, codeWriter: *cw}
	return t

}

func (t *Translator) Run() {
	for _, p := range t.parsers {
		t.codeWriter.Filename = p.filename
		for p.hasMoreLines() {
			p.advance()
			if p.currentCommand == "" || strings.HasPrefix(p.currentCommand, "//") {
				continue
			}
			cmdType := p.commandType()

			if cmdType == C_PUSH || cmdType == C_POP {
				t.codeWriter.writePushPop(p.currentCommand, cmdType, p.arg1(), p.arg2())
			} else if cmdType == C_LABEL {
				t.codeWriter.writeLabel(p.arg1())
			} else if cmdType == C_GOTO {
				t.codeWriter.writeGoto(p.arg1())
			} else if cmdType == C_IF {
				t.codeWriter.writeIf(p.arg1())
			} else if cmdType == C_FUNCTION {
				t.codeWriter.writeFunction(p.arg1(), p.arg2())
			} else if cmdType == C_CALL {
				t.codeWriter.writeCall(p.arg1(), p.arg2())
			} else if cmdType == C_RETURN {
				t.codeWriter.writeReturn()
			} else {
				t.codeWriter.writeArithmetic(p.arg1())
			}
		}

		t.codeWriter.writeInfiniteLoop()
	}

	t.codeWriter.Close()
}
