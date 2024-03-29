package cipher

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"
)

func TestJwtAuth_ValidateJwtToken(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ja := NewJwtAuth(key)
	timeCurrent := time.Now()
	tokenOneHourLater, _ := ja.NewJwtTOken(timeCurrent.Add(time.Hour * 1))
	tokenOneHourBefore, _ := ja.NewJwtTOken(timeCurrent.Add(time.Hour * -1))

	type args struct {
		tokenStr string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "OK",
			args:    args{tokenStr: tokenOneHourLater},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Exipred",
			args:    args{tokenStr: tokenOneHourBefore},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ja.ValidateJwtToken(tt.args.tokenStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("JwtAuth.ValidateJwtToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JwtAuth.ValidateJwtToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
