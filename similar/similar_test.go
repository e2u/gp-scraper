package similar

import (
	"context"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play/collection"
	"testing"
)

func TestPages(t *testing.T) {
	type args struct {
		ctx   context.Context
		appId string
		opt   *Options
		fn    func(result []*collection.Result) bool
	}
	count := 0

	fn := func(results []*collection.Result) bool {
		count += len(results)
		fmt.Println(count, util.JSONString(results))
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
				appId: "com.google.android.apps.nbu.paisa.user",
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
			if err := Pages(tt.args.ctx, tt.args.appId, tt.args.opt, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Pages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
