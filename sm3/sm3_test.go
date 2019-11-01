package sm3

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/t1anchen/gogmlib/utils"
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

func TestContextReset(t *testing.T) {
	c := Context{}
	c.Reset()
	expected := [8]uint32{
		0x7380166f,
		0x4914b2b9,
		0x172442d7,
		0xda8a0600,
		0xa96f30bc,
		0x163138aa,
		0xe38dee4d,
		0xb0fb0e4e}
	actual := c.state
	if actual != expected {
		t.Errorf(`sm3:TestContextReset失败
期望值=%v
实际值=%v`, expected, actual)
	}

}

// -----------------------------------------------------------------------------
// GB/T 32908 附录A 运算实例
// -----------------------------------------------------------------------------

// TestContextPaddingExample1 A.1 填充后的信息
func TestContextPaddingExample1(t *testing.T) {
	input := []byte("abc")
	t.Logf("TestContextPaddingExample1:input = %x", input)

	expected := utils.WordsToBytes([]uint32{
		0x61626380, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000018})
	t.Logf("TestContextPaddingExample1:expected = %x", expected)

	actual := Padding(input)
	t.Logf("TestContextPaddingExample1:actual = %x", actual)

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf(`TestContextPaddingExample1失败
期望值=%x
实际值=%x`, expected, actual)
	}
}

// TestChecksumExample1 A.1 杂凑值
func TestChecksumExample1(t *testing.T) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []byte("abc"))
	actual := Checksum(buf.Bytes())

	expectedBuffer := new(bytes.Buffer)
	binary.Write(expectedBuffer, binary.BigEndian, []uint32{
		0x66c7f0f4, 0x62eeedd9, 0xd1f2d46b, 0xdc10e4e2, 0x4167c487, 0x5cf2f7a2, 0x297da02b, 0x8f4ba8e0})
	expected := expectedBuffer.Bytes()

	if bytes.Compare(expected, actual) != 0 {
		t.Errorf(`sm3:TestChecksum失败
期望值=%x
实际值=%x`, expected, actual)
	}

}

// TestContextPaddingExample2 A.2 填充后的信息
func TestContextPaddingExample2(t *testing.T) {
	input := utils.WordsToBytes([]uint32{
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364})
	t.Logf("TestContextPaddingExample2:input = %x", input)

	expected := utils.WordsToBytes([]uint32{
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x80000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000200})
	t.Logf("TestContextPaddingExample2:expected = %x", expected)

	actual := Padding(input)
	t.Logf("TestContextPaddingExample2:actual = %x", actual)

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf(`TestContextPaddingExample2失败
期望值=%x
实际值=%x`, expected, actual)
	}

	// 	inputBuf := new(bytes.Buffer)
	// 	binary.Write(inputBuf, binary.BigEndian,
	// 	input := inputBuf.Bytes()

	// 	c := Context{
	// 		iv,
	// 		uint64(len(input) * 8),
	// 		[16]uint32{},
	// 		//input,
	// 	}

	// 	expectedBuf := new(bytes.Buffer)
	// 	binary.Write(expectedBuf, binary.BigEndian,
	// 	expected := expectedBuf.Bytes()
	// 	expectedLen := len(expected)

	// 	actual, err := c.Padding()
	// 	if err != nil {
	// 		t.Error(err.Error())
	// 	}

	// 	if expectedLen != len(actual) {
	// 		t.Errorf(`sm3:TestContextPadding失败
	// 期望值长度=%d
	// 实际值长度=%d`, expectedLen, len(actual))
	// 	}

	// 	if bytes.Compare(expected, actual) != 0 {
	// 		t.Errorf(`sm3:TestContextPadding失败
	// 期望值=%x
	// 实际值=%x`, expected, actual)
	// 	}
}

// TestChecksumExample1 A.2
func TestChecksumExample2(t *testing.T) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, []uint32{
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364,
		0x61626364})
	actual := Checksum(buf.Bytes())

	expectedBuffer := new(bytes.Buffer)
	binary.Write(expectedBuffer, binary.BigEndian, []uint32{
		0xdebe9ff9,
		0x2275b8a1,
		0x38604889,
		0xc18e5a4d,
		0x6fdb70e5,
		0x387e5765,
		0x293dcba3,
		0x9c0c5732})
	expected := expectedBuffer.Bytes()

	if bytes.Compare(expected, actual) != 0 {
		t.Errorf(`sm3:TestChecksum失败
期望值=%x
实际值=%x`, expected, actual)
	}

}
