package assembler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction int

// Declare related constants for each direction starting with index 1
const (
	A Instruction = iota + 1 // EnumIndex = 1
	C                        // EnumIndex = 2
	L                        // EnumIndex = 3
)

type Assembler struct {
	SymbolMap  map[string]int
	nextVarAdd int
}

func NewAssembler() *Assembler {
	assembler := &Assembler{}
	assembler.initializeSymbolMap()
	assembler.nextVarAdd = 16
	return assembler
}

func (a *Assembler) initializeSymbolMap() {
	symbolMap := map[string]int{
		"SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4, "SCREEN": 16384, "KBD": 24576,
	}

	// R0 ... R15
	for i := 0; i <= 15; i++ {
		key := "R" + strconv.FormatInt(int64(i), 10)
		symbolMap[key] = i
	}
	fmt.Printf("Initial Symbol Map: \n%v \n\n", symbolMap)
	a.SymbolMap = symbolMap
}

func (a *Assembler) Run(inputFile string) error {
	asmLines := readASMInputFile(inputFile)
	a.populateSymbolsMap(asmLines)
	outputLines := a.parseAndEncodeLines(asmLines)
	writeHackOutputFile("xxx", outputLines)

	return nil
}

func (a *Assembler) populateSymbolsMap(asmLines []string) error {
	lineNumber := 0
	for _, line := range asmLines {
		fmt.Printf("\n%v - %v", lineNumber, line)
		instrType := getInstructionType(line)

		if instrType != L {
			lineNumber += 1
			continue
		}
		label := strings.TrimRight(strings.TrimLeft(line, "("), ")")
		fmt.Printf("\nAdding new symbol to map... %s=%v\n", label, lineNumber)
		a.SymbolMap[label] = lineNumber
	}

	return nil
}

func (a *Assembler) parseAndEncodeLines(asmLines []string) []string {
	encodedLines := []string{}

	for _, asmLine := range asmLines {
		// A, C or L(abel) instruction
		switch instrType := getInstructionType(asmLine); instrType {
		case A:
			encodedLines = append(encodedLines, a.encodeAInstruction(asmLine))
		case C:
			encodedLines = append(encodedLines, a.encodeCInstruction(asmLine))
		case L:
			fmt.Println("Skipping over label declaration")
		}
	}
	return encodedLines
}

// e.g. @12345 -> 0011000000111001
func (a *Assembler) encodeAInstruction(instr string) string {
	label := instr[1:]
	var address int

	constant, err := strconv.ParseInt(label, 0, 16)

	if err != nil {
		// A-Instruction referenced a label/var, e.g. @i or @LOOP
		val, exists := a.SymbolMap[label]
		address = val

		if !exists {
			fmt.Printf("\nAdding variable [%s] to Symbol Map, with address = %v\n", label, a.nextVarAdd)
			a.SymbolMap[label] = a.nextVarAdd
			address = a.nextVarAdd
			a.nextVarAdd += 1
		}

	} else {
		// A-Instruction referenced a constant, e.g. @123
		address = int(constant)
	}

	bin := strconv.FormatInt(int64(address), 2)
	// Return as a 16-bit string representation
	// Where first bit is 0 to signify A-Instruction
	numPadBits := 15 - len(bin)
	return strings.Repeat("0", 1+numPadBits) + fmt.Sprintf(bin)
}

func (a *Assembler) encodeCInstruction(instr string) string {
	return "111" + getCompBits(instr) + getDestBits(instr) + getJumpBits(instr)
}

func getJumpBits(instr string) string {
	if !strings.Contains(instr, ";") {
		return "000"
	}
	jumpInstr := strings.TrimSpace(strings.Split(instr, ";")[1])

	jumpBitsMap := map[string]string{
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}

	jumpBits, ok := jumpBitsMap[jumpInstr]
	if !ok {
		log.Fatal("Error coding jumpInstr")
	}
	return jumpBits
}

func getDestBits(instr string) string {
	if !strings.Contains(instr, "=") {
		return "000"
	}
	destInstr := strings.TrimSpace(strings.Split(instr, "=")[0])
	destBitsMap := map[string]string{
		"M":   "001",
		"D":   "010",
		"DM":  "011",
		"MD":  "011",
		"A":   "100",
		"AM":  "101",
		"MM":  "101",
		"AD":  "110",
		"DA":  "110",
		"ADM": "111",
		"AMD": "111",
		"DMA": "111",
		"DAM": "111",
		"MAD": "111",
		"MDA": "111",
	}

	destBits, ok := destBitsMap[destInstr]
	if !ok {
		log.Fatal("Error coding destInstr")
	}
	return destBits
}

func getCompBits(instr string) string {
	destCompInstr := strings.TrimSpace(strings.Split(instr, ";")[0])
	compInstr := destCompInstr
	if strings.Contains(destCompInstr, "=") {
		compInstr = strings.Split(destCompInstr, "=")[1]
	}

	compBitsMap := map[string]string{
		"0":   "0101010",
		"1":   "0111111",
		"-1":  "0111010",
		"D":   "0001100",
		"A":   "0110000",
		"!D":  "0001101",
		"!A":  "0110001",
		"-D":  "0001111",
		"-A":  "0110011",
		"D+1": "0011111",
		"A+1": "0110111",
		"D-1": "0001110",
		"A-1": "0110010",
		"D+A": "0000010",
		"D-A": "0010011",
		"A-D": "0000111",
		"D&A": "0000000",
		"D|A": "0010101",
		"M":   "1110000",
		"!M":  "1110001",
		"-M":  "1110011",
		"M+1": "1110111",
		"M-1": "1110010",
		"D+M": "1000010",
		"D-M": "1010011",
		"M-D": "1000111",
		"D&M": "1000000",
		"D|M": "1010101",
	}

	compBits, ok := compBitsMap[compInstr]
	if !ok {
		log.Fatal("Error coding compInstr")
	}
	return compBits
}

func getInstructionType(line string) Instruction {
	if line[0] == byte('@') {
		return A
	} else if line[0] == byte('(') {
		return L
	}

	return C
}

func readASMInputFile(inputFile string) []string {
	if !strings.HasSuffix(inputFile, ".asm") {
		log.Fatal(errors.New("InputFile: (%s) must have .asm extension.."))
	}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			lines = append(lines, trimmed)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func writeHackOutputFile(filename string, outputLines []string) {
	file, err := os.OpenFile(fmt.Sprintf("%s.hack", filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer file.Close()

	if err != nil {
		log.Fatalf("Failed when creating Hack output file: %s", err)
	}

	datawriter := bufio.NewWriter(file)
	defer datawriter.Flush()

	for _, line := range outputLines {
		_, _ = datawriter.WriteString(line + "\n")
	}
}
