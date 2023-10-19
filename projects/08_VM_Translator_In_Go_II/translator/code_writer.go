package translator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CodeWriter struct {
	writer     bufio.Writer
	Close      func() // Closure containing references to the resources that need to be closed
	labelId    int    // Used to assign a unique identifier to each label
	segmentMap map[string]string
	filename   string
}

func NewCodeWriter(directory string, filename string) *CodeWriter {
	outputFilename := directory + "/" + filename + ".asm"
	file, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

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

	cw := &CodeWriter{writer: writer, labelId: 1, segmentMap: segmentMap, filename: filename}
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

func (cw *CodeWriter) writePop(command string, segment string, index int) {
	cmds := []string{}
	cmds = append(cmds, "// "+command)

	segmentVar := cw.segmentMap[segment]

	if segment == "temp" {
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, fmt.Sprintf("@%s%v", segmentVar, (5+index)))
		cmds = append(cmds, "M=D")
	} else if segment == "pointer" {
		var thisThat string
		if index == 0 {
			thisThat = "THIS"
		} else {
			thisThat = "THAT"
		}
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@"+thisThat)
		cmds = append(cmds, "M=D")
	} else if segment == "static" {
		varLabel := cw.filename + "." + strconv.FormatInt(int64(index), 10)
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "AM=M-1")
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@"+varLabel)
		cmds = append(cmds, "M=D")
	} else {
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
	}
	cw.appendASMCommands(cmds)
}

func (cw *CodeWriter) writePush(command string, segment string, index int) {
	cmds := []string{}
	cmds = append(cmds, "// "+command)

	segmentVar := cw.segmentMap[segment]

	if segment == "constant" {
		// Pushing a constant
		cmds = append(cmds, fmt.Sprintf("@%v", index))
		cmds = append(cmds, "D=A")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
	} else if segment == "temp" {
		cmds = append(cmds, fmt.Sprintf("@%s%v", segmentVar, (5+index)))
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
	} else if segment == "pointer" {
		var thisThat string
		if index == 0 {
			thisThat = "THIS"
		} else {
			thisThat = "THAT"
		}
		cmds = append(cmds, "@"+thisThat)
		cmds = append(cmds, "D=M")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "A=M")
		cmds = append(cmds, "M=D")
		cmds = append(cmds, "@SP")
		cmds = append(cmds, "M=M+1")
	} else if segment == "static" {
		varLabel := cw.filename + "." + strconv.FormatInt(int64(index), 10)
		cmds = append(cmds, "@"+varLabel)
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

func (cw *CodeWriter) writePushPop(command string, commandType CommandType, segment string, index int) {
	if commandType == C_PUSH {
		cw.writePush(command, segment, index)
	} else {
		cw.writePop(command, segment, index)
	}
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

func (cw *CodeWriter) writeLabel(label string) {
	cmds := []string{}
	cmds = append(cmds, "// label "+label)
	cmds = append(cmds, fmt.Sprintf("(%v)", label))
	cw.appendASMCommands(cmds)

}

func (cw *CodeWriter) writeGoto(label string) {
	cmds := []string{}
	cmds = append(cmds, "// goto "+label)
	cmds = append(cmds, fmt.Sprintf("@%v", label))
	cmds = append(cmds, "0;JMP") // Unconditional jump
	cw.appendASMCommands(cmds)
}

func (cw *CodeWriter) writeIf(label string) {
	cmds := []string{}
	cmds = append(cmds, "// if-goto "+label)
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "AM=M-1")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, fmt.Sprintf("@%v", label))
	cmds = append(cmds, "D;JNE") // Jump if D != 0
	cw.appendASMCommands(cmds)

}

func (cw *CodeWriter) writeFunction(functionName string, nVars int) {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("// function %s %v", functionName, nVars))
	// function entry label
	cmds = append(cmds, fmt.Sprintf("(%s)", functionName))
	// Skip initialization if nArgs == 0
	cmds = append(cmds, fmt.Sprintf("@%v", nVars))
	cmds = append(cmds, "D=A")
	cmds = append(cmds, fmt.Sprintf("@%s$INIT_END", functionName))
	cmds = append(cmds, "D;JLE")
	// Set temp counter var = nArgs
	cmds = append(cmds, "@R13")
	cmds = append(cmds, "M=D")
	// Initialization loop
	cmds = append(cmds, fmt.Sprintf("(%s$INIT_LOOP)", functionName))
	// Push zero onto stack
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "A=M")
	cmds = append(cmds, "M=0")
	// Increment SP
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "M=M+1")
	// Decremept temp args counter
	cmds = append(cmds, "@R13")
	cmds = append(cmds, "MD=M-1")
	// Continue loop if temp args counter > 0
	cmds = append(cmds, fmt.Sprintf("@%s$INIT_LOOP", functionName))
	cmds = append(cmds, "D;JGT")
	// We've pushed 0 onto stack for each arg, now end
	cmds = append(cmds, fmt.Sprintf("(%s$INIT_END)", functionName))

	cw.appendASMCommands(cmds)
}

func (cw *CodeWriter) writeCall(functionName string, nVars int) {

}

func (cw *CodeWriter) writeReturn() {
	// Saved Caller Frame Addresses:
	// LCL-5 = Return Address
	// LCL-4 = LCL Address
	// LCL-3 = ARG Address
	// LCL-2 = THIS Address
	// LCL-1 = THAT Address

	cmds := []string{}
	cmds = append(cmds, "// return")
	// Pop the return value of the top of stack, store value in ARG
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "AM=M-1")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@ARG")
	cmds = append(cmds, "A=M")
	cmds = append(cmds, "M=D")
	// Set the SP to the address after ARG
	cmds = append(cmds, "@ARG")
	cmds = append(cmds, "D=M+1")
	cmds = append(cmds, "@SP")
	cmds = append(cmds, "M=D")
	// Store the return address (the value of LCL-5) in R13
	cmds = append(cmds, "@5")
	cmds = append(cmds, "D=A")
	cmds = append(cmds, "@LCL")
	cmds = append(cmds, "D=M-D")
	cmds = append(cmds, "@R13")
	cmds = append(cmds, "M=D")
	// Set the THAT value to be the value of LCL - 1
	cmds = append(cmds, "@LCL")
	cmds = append(cmds, "A=M-1")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@THAT")
	cmds = append(cmds, "M=D")
	// Set the THIS value be the value of LCL - 2
	cmds = append(cmds, "@2")
	cmds = append(cmds, "D=A")
	cmds = append(cmds, "@LCL")
	cmds = append(cmds, "A=M-D")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@THIS")
	cmds = append(cmds, "M=D")
	// Set the ARG value be the value of LCL - 3
	cmds = append(cmds, "@3")
	cmds = append(cmds, "D=A")
	cmds = append(cmds, "@LCL")
	cmds = append(cmds, "A=M-D")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@ARG")
	cmds = append(cmds, "M=D")
	// Set the LCL value to be the value of LCL - 4
	cmds = append(cmds, "@LCL")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@4")
	cmds = append(cmds, "A=D-A")
	cmds = append(cmds, "D=M")
	cmds = append(cmds, "@LCL")
	cmds = append(cmds, "M=D")
	// Jump to the return address
	cmds = append(cmds, "@R13")
	cmds = append(cmds, "A=M")
	cmds = append(cmds, "0;JMP")

	cw.appendASMCommands(cmds)
}
