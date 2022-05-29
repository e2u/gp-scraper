package app

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		appId string
		opt   *Options
	}
	tests := []struct {
		name string
		args args
		want *App
	}{
		{
			name: "com.google.android.apps.translate",
			args: args{
				appId: "com.google.android.apps.translate",
				opt: &Options{
					Country:  "us",
					Language: "en-US",
				},
			},
			want: nil,
		},
		{
			name: "com.gamedevltd.wwh-free-in-app",
			args: args{
				appId: "com.gamedevltd.wwh",
				opt: &Options{
					Country:  "us",
					Language: "us-EN",
				},
			},
			want: nil,
		},
		{
			name: "org.prowl.torque-paid",
			args: args{
				appId: "org.prowl.torque",
				opt: &Options{
					Country:  "us",
					Language: "us-EN",
				},
			},
			want: nil,
		},
		{
			name: "com.google.android.apps.nbu.paisa.user",
			args: args{
				appId: "com.google.android.apps.nbu.paisa.user",
				opt: &Options{
					Country:  "us",
					Language: "us-EN",
				},
			},
			want: nil,
		},
		{
			name: "com.fast.rupee-free-app",
			args: args{
				appId: "com.fast.rupee",
				opt: &Options{
					Country:  "us",
					Language: "us-EN",
				},
			},
			want: nil,
		},
		{
			name: "pdf.tap.scanner-paid-app",
			args: args{
				appId: "pdf.tap.scanner",
				opt: &Options{
					Country:  "us",
					Language: "us-EN",
				},
			},
			want: nil,
		},
	}

	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := Detail(ctx, tt.args.appId, tt.args.opt)
			if err != nil {
				t.Error(err)
			}
			b, _ := json.MarshalIndent(app, "", "\t")
			fmt.Println(string(b))
		})
	}
}
