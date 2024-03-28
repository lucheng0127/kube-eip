package cipher

import (
	"testing"
	"time"

	"bou.ke/monkey"
)

func TestTOTPAuth_Validate(t *testing.T) {
	sec := "ORUW2ZJAORXSA33GMYQGI5LUPE======"
	timeNow := time.Now()
	after5Sec := timeNow.Add(time.Second * 5)
	after30Sec := timeNow.Add(time.Second * 30)
	before5Sec := timeNow.Add(time.Second * -5)
	before30Sec := timeNow.Add(time.Second * -30)

	monkey.Patch(time.Now, func() time.Time {
		return timeNow
	})
	codeNow, _ := GetTOTPCode(sec)

	type fields struct {
		user   string
		secret string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		pTime   time.Time
	}{
		{
			name:    "A 5",
			fields:  fields{user: "", secret: sec},
			args:    args{data: codeNow},
			want:    true,
			wantErr: false,
			pTime:   after5Sec,
		},
		{
			name:    "B 5",
			fields:  fields{user: "", secret: sec},
			args:    args{data: codeNow},
			want:    true,
			wantErr: false,
			pTime:   before5Sec,
		},
		{
			name:    "A 30",
			fields:  fields{user: "", secret: sec},
			args:    args{data: codeNow},
			want:    false,
			wantErr: true,
			pTime:   after30Sec,
		},
		{
			name:    "B 30",
			fields:  fields{user: "", secret: sec},
			args:    args{data: codeNow},
			want:    false,
			wantErr: true,
			pTime:   before30Sec,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta, _ := NewTOTPAuth(tt.fields.user, tt.fields.secret)

			monkey.Patch(time.Now, func() time.Time {
				return tt.pTime
			})

			got, err := ta.Validate(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TOTPAuth.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TOTPAuth.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
