package suggest

import (
	"context"
	"testing"
)

func TestQuery(t *testing.T) {
	type args struct {
		ctx  context.Context
		term string
		opt  *Options
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "t01",
			args: args{
				ctx:  context.TODO(),
				term: "word",
				opt: &Options{
					Country:  "us",
					Language: "en-US",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.ctx, tt.args.term, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got: %#v", got)
		})
	}
}
