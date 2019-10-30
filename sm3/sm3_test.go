package sm3

import (
	"testing"
)

var (
	p []uint = []uint{0x8542D698, 0x4C044F18, 0xE8B92435, 0xBF6FF7DE, 0x45728191}
)

func TestContextSize(t *testing.T) {
	c := New()
	actual := c.Size()
	expected := 32
	if actual != expected {
		t.Errorf(`sm3:TestContextSize 失败
期望值=%d
实际值=%d`, expected, actual)
	}
}

func TestContextBlockSize(t *testing.T) {
	c := New()
	actual := c.BlockSize()
	expected := 64
	if actual != expected {
		t.Errorf(`sm3:TestContextSize 失败
期望值=%d
实际值=%d`, expected, actual)
	}
}
