package sm2

// Sign 用私钥生成基于 ASN.1 DER-encoded 的签名
func (sk *PrivKey) Sign(userID []byte, data []byte) ([]byte, error) {
	return sk.SignToASN1DER(userID, data)
}
