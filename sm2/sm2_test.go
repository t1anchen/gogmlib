package sm2

// func TestSignAndVerify(t *testing.T) {
// 	msg := []byte("abc")
// 	t.Logf("msg = %x", msg)
// 	priv, err := GenKey(rand.Reader)
// 	if err != nil {
// 		panic("GenerateKey failed")
// 	}
// 	t.Logf("priv = %+v", priv)

// 	hfunc := sm3.New()
// 	hfunc.Write(msg)
// 	hash := hfunc.Sum(nil)
// 	t.Logf("hash(msg) = %x", hash)

// 	sigBytes, err := priv.Sign(rand.Reader, hash)
// 	var sig Signature
// 	asn1.Unmarshal(sigBytes, &sig)
// 	t.Logf("sig(msg) = %+v", sig)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	pk := priv.PK()
// 	ret := pk.Verify(msg, sigBytes)
// 	fmt.Println(ret)
// }
