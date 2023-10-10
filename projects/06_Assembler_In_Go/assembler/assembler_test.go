package assembler

import (
	"fmt"
	"testing"
)

func TestReadASMInputFile(t *testing.T) {
	expected := []string{"@2", "D=A", "@3", "D=D+A", "@0", "M=D"}
	actual := readASMInputFile("Test1.asm")

	equal, ineqIndex := compareSlices(expected, actual)
	if !equal {
		t.Logf("Actual: %s != Expected: %s", actual[ineqIndex], expected[ineqIndex])
		t.Error("Actual != Expected, inequal at index", ineqIndex)
	}
}

func compareSlices(str1, str2 []string) (bool, int) {
	for i, str := range str1 {
		if str != str2[i] {
			return false, i
		}
	}
	return true, -1 // return true if the slices are equal
}

func TestAssembler_Run(t *testing.T) {
	assembler := NewAssembler()
	assembler.Run("Rect.asm")

	t.Log(assembler.SymbolMap)

}

func TestGetDestBits(t *testing.T) {
	actual := getDestBits("MD=M-1")
	expected := "011"

	if actual != expected {
		fmt.Println(actual)
		t.Fail()
	}
}

func TestEncodeCInstruction(t *testing.T) {
	instr := "M=-1"
	actual := "111" + getCompBits(instr) + getDestBits(instr) + getJumpBits(instr)
	expected := "1110111010001000"

	if actual != expected {
		fmt.Println(actual)
		t.Fail()
	}
}
