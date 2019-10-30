package sm2

// Context 加解密上下文
type Context struct {
	buffer []byte
}

// ToString 字符串表示
func (c *Context) ToString() string {
	return "Hello, World"
}
