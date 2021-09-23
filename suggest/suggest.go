package suggest

import (
	"context"
	"encoding/json"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/sirupsen/logrus"
	"net/url"
)

type Options struct {
	Country  string
	Language string
}

const suggestURL = "https://market.android.com/suggest/SuggRequest"

func Query(ctx context.Context, term string, opt *Options) ([]string, error) {
	v := url.Values{}
	v.Set("json", "1")
	v.Set("query", term)
	v.Set("gl", opt.Country)
	v.Set("hl", opt.Language)
	v.Set("c", "3") // c=3 only search apps and games

	_, resp, err := util.HttpGet(ctx, suggestURL+"?="+v.Encode(), nil)
	if err != nil {
		logrus.Errorf("request suggest error=%v", err)
		return nil, err
	}
	var suggests []struct {
		S string `json:"s"`
	}
	if err := json.Unmarshal(resp, &suggests); err != nil {
		logrus.Errorf("parse suggest response error=%v", err)
		return nil, err
	}

	var rs []string
	for _, s := range suggests {
		rs = append(rs, s.S)
	}

	return rs, nil
}
