// Util to calculate hash

package framework

import (
	"crypto"
	"fmt"
	"testing"
)

func TestCalcHash(t *testing.T) {
	type args struct {
		hashType crypto.Hash
		obj      any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Valid result",
			args{crypto.SHA256, []int{}},
			"4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945", // hardcode value
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcHash(tt.args.hashType, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fmt.Sprintf("%x", got) != tt.want {
				t.Errorf("CalcHash() = %v, want %v", fmt.Sprintf("%x", got), tt.want)
			}
		})
	}
}
