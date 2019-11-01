package sm3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/t1anchen/gogmlib/utils"
)

var (
	p                []uint = []uint{0x8542D698, 0x4C044F18, 0xE8B92435, 0xBF6FF7DE, 0x45728191}
	example1MsgInput        = []byte("abc")
	example2MsgInput        = utils.WordsToBytes([]uint32{
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364})
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

// TestPaddingExample1 A.1 填充后的信息
func TestPaddingExample1(t *testing.T) {
	input := example1MsgInput
	t.Logf("TestPaddingExample1:input = %x", input)

	expected := utils.WordsToBytes([]uint32{
		0x61626380, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000018})
	t.Logf("TestPaddingExample1:expected = %x", expected)

	actual := Padding(input)
	t.Logf("TestPaddingExample1:actual = %x", actual)

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf(`TestPaddingExample1失败
期望值=%x
实际值=%x`, expected, actual)
	}
}

func TestContextMessageExpansionExample1(t *testing.T) {
	var (
		w  [68]uint32
		wp [64]uint32
	)
	ctx := Init()
	input := example1MsgInput
	padded := bytes.NewBuffer(Padding(input))
	blockBuf := bytes.NewBuffer(padded.Next(64))
	for blockBuf.Len() > 0 {
		binary.Read(blockBuf, binary.BigEndian, &ctx.buffer)
		t.Logf("TestContextMessageExpansionExample1:ctx.buffer = %x", ctx.buffer)
		ctx.MessageExpansion(&w, &wp)
		blockBuf = bytes.NewBuffer(padded.Next(64))
	}
	t.Logf("TestContextMessageExpansionExample1:actual_w = %x, actual_wp = %x", w, wp)
	actualWStr := fmt.Sprintf("%x", w)
	actualWpStr := fmt.Sprintf("%x", wp)

	expectedW := []uint32{
		0x61626380, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000018,
		0x9092e200, 0x00000000, 0x000c0606, 0x719c70ed, 0x00000000, 0x8001801f, 0x939f7da9, 0x00000000,
		0x2c6fa1f9, 0xadaaef14, 0x00000000, 0x0001801e, 0x9a965f89, 0x49710048, 0x23ce86a1, 0xb2d12f1b,
		0xe1dae338, 0xf8061807, 0x055d68be, 0x86cfd481, 0x1f447d83, 0xd9023dbf, 0x185898e0, 0xe0061807,
		0x050df55c, 0xcde0104c, 0xa5b9c955, 0xa7df0184, 0x6e46cd08, 0xe3babdf8, 0x70caa422, 0x0353af50,
		0xa92dbca1, 0x5f33cfd2, 0xe16f6e89, 0xf70fe941, 0xca5462dc, 0x85a90152, 0x76af6296, 0xc922bdb2,
		0x68378cf5, 0x97585344, 0x09008723, 0x86faee74, 0x2ab908b0, 0x4a64bc50, 0x864e6e08, 0xf07e6590,
		0x325c8f78, 0xaccb8011, 0xe11db9dd, 0xb99c0545}
	expectedWp := []uint32{
		0x61626380, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000018, 0x9092e200, 0x00000000, 0x000c0606, 0x719c70f5,
		0x9092e200, 0x8001801f, 0x93937baf, 0x719c70ed, 0x2c6fa1f9, 0x2dab6f0b, 0x939f7da9, 0x0001801e,
		0xb6f9fe70, 0xe4dbef5c, 0x23ce86a1, 0xb2d0af05, 0x7b4cbcb1, 0xb177184f, 0x2693ee1f, 0x341efb9a,
		0xfe9e9ebb, 0x210425b8, 0x1d05f05e, 0x66c9cc86, 0x1a4988df, 0x14e22df3, 0xbde151b5, 0x47d91983,
		0x6b4b3854, 0x2e5aadb4, 0xd5736d77, 0xa48caed4, 0xc76b71a9, 0xbc89722a, 0x91a5caab, 0xf45c4611,
		0x6379de7d, 0xda9ace80, 0x97c00c1f, 0x3e2d54f3, 0xa263ee29, 0x12f15216, 0x7fafe5b5, 0x4fd853c6,
		0x428e8445, 0xdd3cef14, 0x8f4ee92b, 0x76848be4, 0x18e587c8, 0xe6af3c41, 0x6753d7d5, 0x49e260d5}
	t.Logf("TestContextMessageExpansionExample1:expected_w = %x, expected_wp = %x", expectedW, expectedWp)
	expectedWStr := fmt.Sprintf("%x", expectedW)
	expectedWpStr := fmt.Sprintf("%x", expectedWp)

	if actualWStr != expectedWStr {
		t.Errorf(`TestContextMessageExpansionExample1失败
expectedWStr=%s
actualWStr=%s`, expectedWStr, actualWStr)
	}

	if actualWpStr != expectedWpStr {
		t.Errorf(`TestContextMessageExpansionExample1失败
expectedWpStr=%s
actualWpStr=%s`, expectedWpStr, actualWpStr)
	}

}

// TestChecksumExample1 A.1 杂凑值
// func TestChecksumExample1(t *testing.T) {
// 	buf := new(bytes.Buffer)
// 	binary.Write(buf, binary.BigEndian, []byte("abc"))
// 	actual := Checksum(buf.Bytes())

// 	expectedBuffer := new(bytes.Buffer)
// 	binary.Write(expectedBuffer, binary.BigEndian, []uint32{
// 		0x66c7f0f4, 0x62eeedd9, 0xd1f2d46b, 0xdc10e4e2, 0x4167c487, 0x5cf2f7a2, 0x297da02b, 0x8f4ba8e0})
// 	expected := expectedBuffer.Bytes()

// 	if bytes.Compare(expected, actual) != 0 {
// 		t.Errorf(`sm3:TestChecksum失败
// 期望值=%x
// 实际值=%x`, expected, actual)
// 	}

// }

// TestPaddingExample2 A.2 填充后的信息
func TestPaddingExample2(t *testing.T) {
	input := example2MsgInput
	t.Logf("TestPaddingExample2:input = %x", input)

	expected := utils.WordsToBytes([]uint32{
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x80000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000200})
	t.Logf("TestPaddingExample2:expected = %x", expected)

	actual := Padding(input)
	t.Logf("TestPaddingExample2:actual = %x", actual)

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf(`TestPaddingExample2失败
期望值=%x
实际值=%x`, expected, actual)
	}
}

func TestContextMessageExpansionExample2(t *testing.T) {
	var (
		w  [68]uint32
		wp [64]uint32
	)
	ctx := Init()
	input := example2MsgInput
	padded := bytes.NewBuffer(Padding(input))
	blockBuf := bytes.NewBuffer(padded.Next(64))

	// First Round
	binary.Read(blockBuf, binary.BigEndian, &ctx.buffer)
	t.Logf("TestContextMessageExpansionExample2:ctx.buffer = %x", ctx.buffer)
	ctx.MessageExpansion(&w, &wp)
	blockBuf = bytes.NewBuffer(padded.Next(64))

	// t.Logf("TestContextMessageExpansionExample2:actual_w = %x, actual_wp = %x", w, wp)
	actualWStrFirstBlock := fmt.Sprintf("%x", w)
	actualWpStrFirstBlock := fmt.Sprintf("%x", wp)

	expectedWFirstBlock := []uint32{
		0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364, 0x61626364,
		0xa121a024, 0xa121a024, 0xa121a024, 0x6061e0e5, 0x6061e0e5, 0x6061e0e5, 0xa002e345, 0xa002e345,
		0xa002e345, 0x49c969ed, 0x49c969ed, 0x49c969ed, 0x85ae5679, 0xa44ff619, 0xa44ff619, 0x694b6244,
		0xe8c8e0c4, 0xe8c8e0c4, 0x240e103e, 0x346e603e, 0x346e603e, 0x9a517ab5, 0x8a01aa25, 0x8a01aa25,
		0x0607191c, 0x25f8a37a, 0xd528936a, 0x89fbd8ae, 0x00606206, 0x10501256, 0x7cff7ef9, 0x3c78b9f9,
		0xcc2b8a69, 0x9f03f169, 0xdf45be20, 0x9ec5bee1, 0x0a212906, 0x49ff72c0, 0x46717241, 0x67e09a19,
		0x6efaa333, 0x2ebae676, 0x3475c386, 0x201dcff6, 0x2f18fccf, 0x2c5f2b5c, 0xa80b9f38, 0xbc139f34,
		0xc47f18a7, 0xa25ce71d, 0x42743705, 0x51baf619}
	expectedWpFirstBlock := []uint32{
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0xc043c340, 0xc043c340, 0xc043c340, 0x01038381,
		0xc14040c1, 0xc14040c1, 0x01234361, 0xc06303a0, 0xc06303a0, 0x29a88908, 0xe9cb8aa8, 0xe9cb8aa8,
		0x25acb53c, 0xed869ff4, 0xed869ff4, 0x20820ba9, 0x6d66b6bd, 0x4c8716dd, 0x8041e627, 0x5d25027a,
		0xdca680fa, 0x72999a71, 0xae0fba1b, 0xbe6fca1b, 0x32697922, 0xbfa9d9cf, 0x5f29394f, 0x03fa728b,
		0x06677b1a, 0x35a8b12c, 0xa9d7ed93, 0xb5836157, 0xcc4be86f, 0x8f53e33f, 0xa3bac0d9, 0xa2bd0718,
		0xc60aa36f, 0xd6fc83a9, 0x9934cc61, 0xf92524f8, 0x64db8a35, 0x674594b6, 0x7204b1c7, 0x47fd55ef,
		0x41e25ffc, 0x02e5cd2a, 0x9c7e5cbe, 0x9c0e50c2, 0xeb67e468, 0x8e03cc41, 0xea7fa83d, 0xeda9692d}
	// t.Logf("TestContextMessageExpansionExample2:expected_w = %x, expected_wp = %x", expectedW, expectedWp)
	expectedWStrFirstBlock := fmt.Sprintf("%x", expectedWFirstBlock)
	expectedWpStrFirstBlock := fmt.Sprintf("%x", expectedWpFirstBlock)

	if actualWStrFirstBlock != expectedWStrFirstBlock {
		t.Errorf(`TestContextMessageExpansionExample2失败
expectedWStr=%s
actualWStr=%s`, expectedWStrFirstBlock, actualWStrFirstBlock)
	}

	if actualWpStrFirstBlock != expectedWpStrFirstBlock {
		t.Errorf(`TestContextMessageExpansionExample2失败
expectedWpStr=%s
actualWpStr=%s`, expectedWpStrFirstBlock, actualWpStrFirstBlock)
	}

	// Second Round
	binary.Read(blockBuf, binary.BigEndian, &ctx.buffer)
	t.Logf("TestContextMessageExpansionExample2:ctx.buffer = %x", ctx.buffer)
	ctx.MessageExpansion(&w, &wp)
	blockBuf = bytes.NewBuffer(padded.Next(64))

	actualWStrSecondBlock := fmt.Sprintf("%x", w)
	actualWpStrSecondBlock := fmt.Sprintf("%x", wp)

	expectedWSecondBlock := []uint32{
		0x80000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000200,
		0x80404000, 0x00000000, 0x01008080, 0x10005000, 0x00000000, 0x002002a0, 0xac545c04, 0x00000000,
		0x09582a39, 0xa0003000, 0x00000000, 0x00200280, 0xa4515804, 0x20200040, 0x51609838, 0x30005701,
		0xa0002000, 0x008200aa, 0x6ad525d0, 0x0a0e0216, 0xb0f52042, 0xfa7073b0, 0x20000000, 0x008200a8,
		0x7a542590, 0x22a20044, 0xd5d6ebd2, 0x82005771, 0x8a202240, 0xb42826aa, 0xeaf84e59, 0x4898eaf9,
		0x8207283d, 0xee6775fa, 0xa3e0e0a0, 0x8828488a, 0x23b45a5d, 0x628a22c4, 0x8d6d0615, 0x38300a7e,
		0xe96260e5, 0x2b60c020, 0x502ed531, 0x9e878cb9, 0x218c38f8, 0xdcae3cb7, 0x2a3e0e0a, 0xe9e0c461,
		0x8c3e3831, 0x44aaa228, 0xdc60a38b, 0x518300f7}
	expectedWpSecondBlock := []uint32{
		0x80000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000200, 0x80404000, 0x00000000, 0x01008080, 0x10005200,
		0x80404000, 0x002002a0, 0xad54dc84, 0x10005000, 0x09582a39, 0xa02032a0, 0xac545c04, 0x00200280,
		0xad09723d, 0x80203040, 0x51609838, 0x30205581, 0x04517804, 0x20a200ea, 0x3bb5bde8, 0x3a0e5517,
		0x10f50042, 0xfaf2731a, 0x4ad525d0, 0x0a8c02be, 0xcaa105d2, 0xd8d273f4, 0xf5d6ebd2, 0x828257d9,
		0xf07407d0, 0x968a26ee, 0x3f2ea58b, 0xca98bd88, 0x08270a7d, 0x5a4f5350, 0x4918aef9, 0xc0b0a273,
		0xa1b37260, 0x8ced573e, 0x2e8de6b5, 0xb01842f4, 0xcad63ab8, 0x49eae2e4, 0xdd43d324, 0xa6b786c7,
		0xc8ee581d, 0xf7cefc97, 0x7a10db3b, 0x776748d8, 0xadb200c9, 0x98049e9f, 0xf65ead81, 0xb863c496}
	// t.Logf("TestContextMessageExpansionExample2:expected_w = %x, expected_wp = %x", expectedW, expectedWp)
	expectedWStrSecondBlock := fmt.Sprintf("%x", expectedWSecondBlock)
	expectedWpStrSecondBlock := fmt.Sprintf("%x", expectedWpSecondBlock)

	if actualWStrSecondBlock != expectedWStrSecondBlock {
		t.Errorf(`TestContextMessageExpansionExample2失败
expectedWStr=%s
actualWStr=%s`, expectedWStrSecondBlock, actualWStrSecondBlock)
	}

	if actualWpStrSecondBlock != expectedWpStrSecondBlock {
		t.Errorf(`TestContextMessageExpansionExample2失败
expectedWpStr=%s
actualWpStr=%s`, expectedWpStrSecondBlock, actualWpStrSecondBlock)
	}

}

// TestChecksumExample1 A.2
// func TestChecksumExample2(t *testing.T) {
// 	buf := new(bytes.Buffer)
// 	binary.Write(buf, binary.BigEndian, []uint32{
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364,
// 		0x61626364})
// 	actual := Checksum(buf.Bytes())

// 	expectedBuffer := new(bytes.Buffer)
// 	binary.Write(expectedBuffer, binary.BigEndian, []uint32{
// 		0xdebe9ff9,
// 		0x2275b8a1,
// 		0x38604889,
// 		0xc18e5a4d,
// 		0x6fdb70e5,
// 		0x387e5765,
// 		0x293dcba3,
// 		0x9c0c5732})
// 	expected := expectedBuffer.Bytes()

// 	if bytes.Compare(expected, actual) != 0 {
// 		t.Errorf(`sm3:TestChecksum失败
// 期望值=%x
// 实际值=%x`, expected, actual)
// 	}

// }
