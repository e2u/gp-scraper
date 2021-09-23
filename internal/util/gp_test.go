package util

import (
	"context"
	"fmt"
	"testing"
)

func TestBatchExecute(t *testing.T) {
	type args struct {
		payload  string
		country  string
		language string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t01",
			args: args{
				payload:  "f.req=%5B%5B%5B%22UsvDTd%22%2C%22%5Bnull%2Cnull%2C%5B2%2C1%2C%5B50%2Cnull%2Cnull%5D%2Cnull%2C%5B%5D%5D%2C%5B%5C%22com.google.android.apps.translateFFF%5C%22%2C7%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D",
				country:  "us",
				language: "en-US",
			},
			want:    "",
			wantErr: false,
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BatchExecute(ctx, tt.args.payload, tt.args.country, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchExecute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestExtractEmbedData(t *testing.T) {
	ctx := context.TODO()
	_, html, err := HttpGet(ctx, "https://play.google.com/store/apps/details?id=com.google.android.apps.translate", nil)
	if err != nil {
		panic(err)
	}
	type args struct {
		html []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "t02",
			args: args{
				html: html,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractEmbedData(tt.args.html)
			for k, v := range result {
				fmt.Println(k, ">>>>", v)
			}
		})
	}
}
