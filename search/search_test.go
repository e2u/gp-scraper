package category

import (
	"context"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play/collection"
	"github.com/e2u/gp-scraper/play/price"
	"testing"
)

func TestPages(t *testing.T) {
	var count int

	fn := func(results []*collection.Result) bool {
		count += len(results)
		t.Log(count, util.JSONString(results))
		return true
	}

	type args struct {
		ctx   context.Context
		query string
		opt   *Options
		fn    func(result []*collection.Result) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "google-paid",
			args: args{
				ctx:   context.TODO(),
				query: "google",
				opt: &Options{
					PriceType: price.Paid,
					Country:   "us",
					Language:  "en-US",
				},
				fn: fn,
			},
			wantErr: false,
		},
		{
			name: "google-free",
			args: args{
				ctx:   context.TODO(),
				query: "google",
				opt: &Options{
					PriceType: price.Free,
					Country:   "us",
					Language:  "en-US",
				},
				fn: fn,
			},
			wantErr: false,
		},
		{
			name: "google-all",
			args: args{
				ctx:   context.TODO(),
				query: "google",
				opt: &Options{
					PriceType: price.All,
					Country:   "us",
					Language:  "en-US",
				},
				fn: fn,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Pages(tt.args.ctx, tt.args.query, tt.args.opt, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Pages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
