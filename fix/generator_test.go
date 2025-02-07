package fix

import (
	"bytes"
	"testing"
)

func TestCalcCheckSum(t *testing.T) {
	cases := [][]byte{
		[]byte("8=FIX.4.2|9="),
		[]byte("8=FIX.4.2|9=0"),
		[]byte("8=FIX.4.2|9=0|"),
		[]byte(""),
		nil,
	}
	for _, c := range cases {
		t.Run(string(c), func(t *testing.T) {
			if !bytes.Equal(CalcCheckSumOptimized(c), CalcCheckSum(c)) {
				t.Fatalf("CalcCheckSumOptimized %s != CalcCheckSum %s", CalcCheckSumOptimized(c), CalcCheckSum(c))
			}
		})
	}
}
