package sm4

import (
	"testing"
)

var (
	p []uint = []uint{0x8542D698, 0x4C044F18, 0xE8B92435, 0xBF6FF7DE, 0x45728191}
)

func TestContextToString(t *testing.T) {
	c := Context{}
	actual := c.ToString()
	expected := "Hello"
	if actual != expected {
		t.Errorf(`ToString 测试失败
期望值=%s
实际值=%s`, expected, actual)
	}
}
