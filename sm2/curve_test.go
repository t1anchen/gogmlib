package sm2

import (
	"testing"
)

func TestGetSm2P256(t *testing.T) {
	curve := GetSm2P256()
	t.Logf("P:%s\n", curve.Params().P.Text(16))
	t.Logf("B:%s\n", curve.Params().B.Text(16))
	t.Logf("N:%s\n", curve.Params().N.Text(16))
	t.Logf("Gx:%s\n", curve.Params().Gx.Text(16))
	t.Logf("Gy:%s\n", curve.Params().Gy.Text(16))
}
