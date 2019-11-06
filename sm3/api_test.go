package sm3

import (
	"fmt"
	"testing"

	"github.com/t1anchen/gogmlib/utils"
)

func TestApiReset(t *testing.T) {
	c := New()
	t.Logf("api_test:TestNew:c = %x\n", c)
}

func TestSm3ApiExample1(t *testing.T) {
	input := example1MsgInput
	h := New()
	h.Write(input)
	actual := h.Sum(nil)
	actualStr := fmt.Sprintf("%x", actual)

	expected := utils.WordsToBytes([]uint32{
		0x66c7f0f4, 0x62eeedd9, 0xd1f2d46b, 0xdc10e4e2,
		0x4167c487, 0x5cf2f7a2, 0x297da02b, 0x8f4ba8e0})
	expectedStr := fmt.Sprintf("%x", expected)

	if actualStr != expectedStr {
		t.Errorf(`TestSm3ApiExample1失败
期望值=%s
实际值=%s`, expectedStr, actualStr)
	}

}

func TestSm3ApiExample2(t *testing.T) {
	input := example2MsgInput
	h := New()
	h.Write(input)
	actual := h.Sum(nil)
	actualStr := fmt.Sprintf("%x", actual)

	expected := utils.WordsToBytes([]uint32{
		0xdebe9ff9, 0x2275b8a1, 0x38604889, 0xc18e5a4d,
		0x6fdb70e5, 0x387e5765, 0x293dcba3, 0x9c0c5732})
	expectedStr := fmt.Sprintf("%x", expected)

	if actualStr != expectedStr {
		t.Errorf(`TestSm3ApiExample2失败
期望值=%s
实际值=%s`, expectedStr, actualStr)
	}
}
