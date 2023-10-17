package translator

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

func TestTranslator_BasicTest(t *testing.T) {
	translator := NewTranslator("BasicTest")
	translator.Run()

	expected := readASMFileContents("BasicTestExpected.asm")
	actual := readASMFileContents("BasicTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_StackTest(t *testing.T) {
	translator := NewTranslator("StackTest")
	translator.Run()

	expected := readASMFileContents("StackTestExpected.asm")
	actual := readASMFileContents("StackTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_SimpleAdd(t *testing.T) {
	translator := NewTranslator("SimpleAdd")
	translator.Run()

	expected := readASMFileContents("SimpleAddExpected.asm")
	actual := readASMFileContents("SimpleAdd.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_PointerTest(t *testing.T) {
	translator := NewTranslator("PointerTest")
	translator.Run()

	expected := readASMFileContents("PointerTestExpected.asm")
	actual := readASMFileContents("PointerTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_StaticTest(t *testing.T) {
	translator := NewTranslator("StaticTest")
	translator.Run()

	expected := readASMFileContents("StaticTestExpected.asm")
	actual := readASMFileContents("StaticTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func readASMFileContents(filename string) []string {
	if !strings.HasSuffix(filename, ".asm") {
		log.Fatal(errors.New("Filename: (%s) must have .asm extension.."))
	}

	file, err := os.Open("testfiles/" + filename)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	var asmLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		asmLines = append(asmLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return asmLines
}

func assertSlicesEqual(t *testing.T, expected, actual []string) {
	if len(expected) != len(actual) {
		t.Fatalf("Slice are not equal, lengths differ..\n")
	}
	for i, element := range expected {
		if element != actual[i] {
			t.Logf("Expected: %v != Actual: %v", element, actual[i])
			t.Fatalf("Slice equality check failed on index: %v", i)
		}
	}
	t.Logf("Slices with length %v were equal!\n", len(expected))
}
