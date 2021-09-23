package developer

import (
	"context"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play/collection"
	"testing"
)

func TestPages(t *testing.T) {
	type args struct {
		ctx   context.Context
		devId string
		opt   *Options
		fn    func(result []*collection.Result) bool
	}
	count := 0
	fn := func(results []*collection.Result) bool {
		count += len(results)
		t.Log(count, util.JSONString(results))
		return true
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t01-id",
			args: args{
				ctx:   context.TODO(),
				devId: "5700313618786177705",
				opt: &Options{
					Country:  "us",
					Language: "en-US",
				},
				fn: fn,
			},
			wantErr: false,
		},
		{
			name: "t01-developer",
			args: args{
				ctx:   context.TODO(),
				devId: "Google LLC",
				opt: &Options{
					Country:  "us",
					Language: "en-US",
				},
				fn: fn,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Pages(tt.args.ctx, tt.args.devId, tt.args.opt, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Pages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
