package translator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type CodeWriter struct {
	writer     bufio.Writer
	Close      func()
	labelId    int
	segmentMap map[string]string
}

func NewCodeWriter(outputFile string) *CodeWriter {
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		log.Fatalf("Failed when creating Hack Asm output file: %s", err)
	}

	writer := *bufio.NewWriter(file)

	segmentMap := map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
		"temp":     "R",
	}

	cw := &CodeWriter{writer: writer, labelId: 1, segmentMap: segmentMap}
	cw.Close = func() {
		writer.Flush()
		file.Close()
	}

	return cw
}

func (cw *CodeWriter) writeArithmetic(command string) {
	cmds := []string{}

	if command == "eq" || command == "lt" || command == "gt" {
		cw.writeComparison(command)
	} else if command == "neg" {
		// Neg
		cmds = append(cmds, "// "+command)
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M-1")
		cmds = append(cmds, "M=-M")
		cw.appendASMCommands(cmds)
	} else if command == "not" {
		// Neg
		cmds = append(cmds, "// "+command)
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M-1")
		cmds = append(cmds, "M=!M")
		cw.appendASMCommands(cmds)
	} else {
		var op string
		switch command {
		case "add":
			op = "+"
		case "sub":
			op = "-"
		case "or":
			op = "|"
		case "and":
			op = "&"
		}
		// Add
		cmds = append(cmds, "// "+command)
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "M=M"+op+"D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
		cw.appendASMCommands(cmds)
	}
}

func (cw *CodeWriter) writeComparison(command string) {
	// eq, gt, lt
	op := "J" + strings.ToUpper(command)
	labelPrefix := fmt.Sprintf("%v_%v", strings.ToUpper(command), cw.labelId)
	cw.labelId++

	cmds := []string{}
	cmds = append(cmds, "// "+command)
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "AM=M-1")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "A=M-1")
	cmds = append(cmds, "D=M-D")
	cmds = append(cmds, "@"+labelPrefix+"_TRUE")
	cmds = append(cmds, "D;"+op)
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "A=M-1")
	cmds = append(cmds, "M=0")
	cmds = append(cmds, "@"+labelPrefix+"_END")
	cmds = append(cmds, "0;JMP")
	cmds = append(cmds, "("+labelPrefix+"_TRUE)")
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "A=M-1")
	cmds = append(cmds, "M=-1")
	cmds = append(cmds, "("+labelPrefix+"_END)")

	cw.appendASMCommands(cmds)
}

func (cw *CodeWriter) writePushPop(command string, commandType CommandType, segment string, index int) {
	cmds := []string{}
	cmds = append(cmds, "// "+command)

	segmentVar := cw.segmentMap[segment]

	if commandType == C_POP && segment == "temp" {
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, fmt.Sprintf("@%s%v", segmentVar, (5+index)))
		cmds = append(cmds, "M=D")
	} else if commandType == C_POP {
		cmds = append(cmds, fmt.Sprintf("@%s", segmentVar))
		cmds = append(cmds, "D=M")
		cmds = append(cmds, fmt.Sprintf("@%v", index))
		cmds = append(cmds, "D=D+A")
		cmds = append(cmds, "@R13")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@R13")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
	} else if commandType == C_PUSH && segment == "constant" {
		// Pushing a constant
		cmds = append(cmds, fmt.Sprintf("@%v", index))
		cmds = append(cmds, "D=A")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
	} else if commandType == C_PUSH && segment == "temp" {
		cmds = append(cmds, fmt.Sprintf("@%s%v", segmentVar, (5+index)))
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
	} else {
		cmds = append(cmds, fmt.Sprintf("@%s", segmentVar))
		cmds = append(cmds, "D=M")
		cmds = append(cmds, fmt.Sprintf("@%v", index))
		cmds = append(cmds, "A=D+A")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
	}
	cw.appendASMCommands(cmds)
}

func (cw *CodeWriter) appendASMCommands(asmCommands []string) {
	for _, cmd := range asmCommands {
		_, err := cw.writer.WriteString(cmd + "\n")

		if err != nil {
			log.Fatal(err)
		}
	}

	cw.writer.Flush()
}

func (cw *CodeWriter) writeInfiniteLoop() {
	cmds := []string{}
	cmds = append(cmds, "// Infinite Loop")
	cmds = append(cmds, "(END)")
	cmds = append(cmds, "@END")
	cmds = append(cmds, "0;JMP")
	cw.appendASMCommands(cmds)
}
