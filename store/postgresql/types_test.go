package postgresql

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestTypes_GetAndReturn(t *testing.T) {

	t.Run("ok", func(t *testing.T) {

		// 0x840c8122433a5aa7ad60c1bcdc36ab9dccf761a5
		a := []byte{132, 12, 129, 34, 67, 58, 90, 167, 173, 96, 193, 188, 220, 54, 171, 157, 204, 247, 97, 165}
		ad := common.BytesToAddress(a)
		t.Log(ad.Hex())
		t.Log(ad.String())
		t.Log(ad.Hash().Big().String())
		h := ad.Hash().String()
		t.Log(h)

		bytesA := ad.Bytes()
		p := new(big.Int)
		p.SetBytes(bytesA)
		t.Log(p)

		pad := common.Address{}
		pad.SetBytes(p.Bytes())
		t.Log(pad)
		t.Log(pad.Hex())

	})
}
