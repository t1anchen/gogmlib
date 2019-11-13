package sm2

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/t1anchen/gogmlib/utils"
)

func TestGenKey(t *testing.T) {
	sk, pk, err := GenKey(rand.Reader)
	t.Logf("pk = %s\nsk = %s\nerr = %s", pk, sk, err)
}

func TestSign(t *testing.T) {
	sk := &PrivKey{
		D:     utils.NewBigIntFromHexString("5DD701828C424B84C5D56770ECF7C4FE882E654CAC53C7CC89A66B1709068B9D"),
		Curve: GetSm2P256()}
	t.Logf("priv.D = %x", sk.D)
	sig, err := sk.Sign(nil, utils.HexStringToBytes("0102030405060708010203040506070801020304050607080102030405060708"))
	if err != nil {
		t.Errorf("TestSign:Sign:err = %s\n", err)
		return
	}
	t.Logf("TestSign:Sign:sig = %x\n", sig)
}

func TestSignToBigInt(t *testing.T) {
	sk := &PrivKey{
		D:     utils.NewBigIntFromHexString("5DD701828C424B84C5D56770ECF7C4FE882E654CAC53C7CC89A66B1709068B9D"),
		Curve: GetSm2P256()}
	msg := utils.HexStringToBytes("0102030405060708010203040506070801020304050607080102030405060708")
	actualR, actualS, err := sk.SignToBigInt(
		sm2SignDefaultUserID,
		msg)
	t.Logf("actualR:%x\nactualS:%x\nerr:%v",
		utils.BigIntToBytes(actualR),
		utils.BigIntToBytes(actualS),
		err)

	pk := &PubKey{
		Curve: GetSm2P256(),
		X:     utils.NewBigIntFromHexString("FF6712D3A7FC0D1B9E01FF471A87EA87525E47C7775039D19304E554DEFE0913"),
		Y:     utils.NewBigIntFromHexString("F632025F692776D4C13470ECA36AC85D560E794E1BCCF53D82C015988E0EB956")}

	actualResult := pk.VerifyFromBigInt(nil, msg, actualR, actualS)
	t.Logf("actualResult:%v", actualResult)
}

func TestSignToASN1DER(t *testing.T) {
	sk := &PrivKey{
		D:     utils.NewBigIntFromHexString("5DD701828C424B84C5D56770ECF7C4FE882E654CAC53C7CC89A66B1709068B9D"),
		Curve: GetSm2P256()}
	msg := utils.HexStringToBytes("0102030405060708010203040506070801020304050607080102030405060708")
	sigBytes, _ := sk.SignToASN1DER(nil, msg)
	t.Logf("actualSigBytes:%x", sigBytes)
	pk := &PubKey{
		Curve: GetSm2P256(),
		X:     utils.NewBigIntFromHexString("FF6712D3A7FC0D1B9E01FF471A87EA87525E47C7775039D19304E554DEFE0913"),
		Y:     utils.NewBigIntFromHexString("F632025F692776D4C13470ECA36AC85D560E794E1BCCF53D82C015988E0EB956")}

	actualResult := pk.VerifyFromASN1DER(nil, msg, sigBytes)
	t.Logf("actualResult:%v", actualResult)
}

func TestEncrypt(t *testing.T) {
	msg := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sk, pk, err := GenKey(rand.Reader)
	if err != nil {
		t.Error(err.Error())
		return
	}

	encrypted, err := pk.Encrypt(msg)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("cipher text:%x\n", encrypted)

	plain, err := sk.Decrypt(encrypted)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("plain text:%x\n", plain)

	if !bytes.Equal(msg, plain) {
		t.Errorf("加解密失败，msg=%x 不等于 plain=%x", msg, plain)
		return
	}

}

func TestEncryptWithASN1DER(t *testing.T) {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sk, pk, err := GenKey(rand.Reader)
	if err != nil {
		t.Error(err.Error())
		return
	}

	cipherText, err := pk.Encrypt(src)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("cipher text:%x\n", cipherText)

	derCipher, err := MarshalEncryptedToASN1DER(cipherText)
	if err != nil {
		t.Error(err.Error())
		return
	}
	//err = ioutil.WriteFile("derCipher.dat", derCipher, 0644)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	cipherText, err = UnmarshalEncryptedFromASN1DER(derCipher)
	if err != nil {
		t.Error(err.Error())
		return
	}

	plainText, err := sk.Decrypt(cipherText)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("plain text:%x\n", plainText)

	if !bytes.Equal(src, plainText) {
		t.Errorf("加解密失败，msg=%x 不等于 plain=%x", src, plainText)
		return
	}
}
