package translator

import (
	"testing"
)

func TestTranslator_Run(t *testing.T) {
	translator := NewTranslator("StaticTest.vm")
	translator.Run()
}
