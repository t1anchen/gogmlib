package sm4

// Context 加密上下文
type Context struct {
	buffer []byte
}

// ToString 上下文:字符串表示
func (c *Context) ToString() string {
	return "Hello"
}
