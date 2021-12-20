package reviews

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play/sort"
	"testing"
)

func TestPages(t *testing.T) {
	type args struct {
		ctx   context.Context
		appId string
		opt   *Options
		fn    func([]*Review) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t01",
			args: args{
				ctx:   context.TODO(),
				appId: "com.gamedevltd.wwh",
				opt: &Options{
					Country:    "us",
					Language:   "en-US",
					PageNumber: 200,
					Sorting:    sort.Rating,
				},
				fn: func(reviews []*Review) bool {
					b, _ := json.MarshalIndent(reviews, "", "\t")
					fmt.Println(">", string(b))
					return false
				},
			},
			wantErr: false,
		},
		{
			name: "t02",
			args: args{
				ctx:   context.TODO(),
				appId: "com.memo.cash",
				opt: &Options{
					Country:    "id",
					Language:   "id",
					PageNumber: 10,
					Sorting:    sort.Rating,
				},
				fn: func(reviews []*Review) bool {
					// b, _ := json.MarshalIndent(reviews, "", "\t")
					// fmt.Println(">",string(b))
					fmt.Println(util.IdentJSONString(reviews))
					return false
				},
			},
			wantErr: false,
		},
		{
			name: "t03-replay",
			args: args{
				ctx:   context.TODO(),
				appId: "com.funcamerastudio.videomaker",
				opt: &Options{
					Country:    "tw",
					Language:   "zh-TW",
					PageNumber: 200,
					Sorting:    sort.Helpfulness,
				},
				fn: func(reviews []*Review) bool {
					b, _ := json.MarshalIndent(reviews, "", "\t")
					fmt.Println(">", string(b))
					return true
				},
			},
			wantErr: false,
		},
		{
			name: "t04",
			args: args{
				ctx:   context.TODO(),
				appId: "com.fastrupee.fastrupeepeefa",
				opt: &Options{
					Country:    "in",
					Language:   "en-IN",
					PageNumber: 10,
					Sorting:    sort.Rating,
				},
				fn: func(reviews []*Review) bool {
					// b, _ := json.MarshalIndent(reviews, "", "\t")
					// fmt.Println(">",string(b))
					fmt.Println(util.IdentJSONString(reviews))
					return false
				},
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
