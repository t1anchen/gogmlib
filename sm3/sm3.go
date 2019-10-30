package sm3

var (
	// GB/T 32907 4.1 初始值
	iv = [8]uint{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600, 0xa96f30bc,
		0x163138aa, 0xe38dee4d, 0xb0fb0e4e}
)

type Context struct {
	buffer []byte
}
