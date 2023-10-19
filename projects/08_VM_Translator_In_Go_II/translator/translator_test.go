package translator

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

// Project 07 Tests
func TestTranslator_BasicTest(t *testing.T) {
	translator := NewTranslator("BasicTest.vm", false)
	translator.Run()

	expected := readASMFileContents("BasicTestExpected.asm")
	actual := readASMFileContents("BasicTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_StackTest(t *testing.T) {
	translator := NewTranslator("StackTest.vm", false)
	translator.Run()

	expected := readASMFileContents("StackTestExpected.asm")
	actual := readASMFileContents("StackTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_SimpleAdd(t *testing.T) {
	translator := NewTranslator("SimpleAdd.vm", false)
	translator.Run()

	expected := readASMFileContents("SimpleAddExpected.asm")
	actual := readASMFileContents("SimpleAdd.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_PointerTest(t *testing.T) {
	translator := NewTranslator("PointerTest.vm", false)
	translator.Run()

	expected := readASMFileContents("PointerTestExpected.asm")
	actual := readASMFileContents("PointerTest.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_StaticTest(t *testing.T) {
	translator := NewTranslator("StaticTest.vm", false)
	translator.Run()

	expected := readASMFileContents("StaticTestExpected.asm")
	actual := readASMFileContents("StaticTest.asm")

	assertSlicesEqual(t, expected, actual)
}

// Project 08 Tests
func TestTranslator_BasicLoop(t *testing.T) {
	translator := NewTranslator("BasicLoop.vm", false)
	translator.Run()

	expected := readASMFileContents("BasicLoopExpected.asm")
	actual := readASMFileContents("BasicLoop.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_FibonacciSeries(t *testing.T) {
	translator := NewTranslator("FibonacciSeries.vm", false)
	translator.Run()

	expected := readASMFileContents("FibonacciSeriesExpected.asm")
	actual := readASMFileContents("FibonacciSeries.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_SimpleFunction(t *testing.T) {
	translator := NewTranslator("SimpleFunction.vm", false)
	translator.Run()

	expected := readASMFileContents("SimpleFunctionExpected.asm")
	actual := readASMFileContents("SimpleFunction.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_FibonacciElement(t *testing.T) {
	translator := NewTranslator("FibonacciElement", true)
	translator.Run()

	expected := readASMFileContents("FibonacciElement/FibonacciElementExpected.asm")
	actual := readASMFileContents("FibonacciElement/FibonacciElement.asm")

	assertSlicesEqual(t, expected, actual)
}

func TestTranslator_StaticsTest(t *testing.T) {
	translator := NewTranslator("StaticsTest", true)
	translator.Run()

	expected := readASMFileContents("StaticsTest/StaticsTestExpected.asm")
	actual := readASMFileContents("StaticsTest/StaticsTest.asm")

	assertSlicesEqual(t, expected, actual)
}

// Helper Functions
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
