package category

import (
	"context"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play/age"
	"github.com/e2u/gp-scraper/play/category"
	"github.com/e2u/gp-scraper/play/collection"
	"github.com/e2u/gp-scraper/play/price"
	"github.com/e2u/gp-scraper/play/sort"
	"testing"
)

func TestList(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts *collection.Options
		fn   func([]*collection.Result) bool
	}
	count := 0
	fn := func(results []*collection.Result) bool {
		// fmt.Println(util.IdentJSONString(results))
		count += len(results)
		fmt.Println(count, util.JSONString(results))
		return true
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "top-ration",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort: sort.Rating,
				},
				fn: fn,
			},
		},

		{
			name: "top-ration-limit",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort: sort.Rating,
				},
				fn: fn,
			},
		},

		{
			name: "top-ration-limit-min",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort:     sort.Rating,
					ScoreMin: 4.8,
				},
				fn: fn,
			},
		},

		{
			name: "top-new",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort: sort.Newest,
				},
				fn: fn,
			},
		},
		{
			name: "t01",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort:     sort.Rating,
					Category: category.Business,
				},
				fn: fn,
			},
		},
		{
			name: "t02",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort:     sort.Newest,
					Category: category.Game,
					Age:      age.FiveUnder,
				},
				fn: fn,
			},
		},
		{
			name: "thailand-free",
			args: args{
				ctx: context.TODO(),
				opts: &collection.Options{
					Sort:       sort.Rating,
					Category:   category.Finance,
					PriceType:  price.Free,
					Country:    "th",
					Language:   "th",
					PageNumber: 500,
				},
				fn: fn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Pages(tt.args.ctx, tt.args.opts, tt.args.fn)
		})
	}
}
