package category

import (
	"context"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/collection"
	"github.com/e2u/gp-scraper/play/price"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
)

type Options struct {
	Language  string
	Country   string
	PriceType price.Type
}

func Pages(ctx context.Context, query string, opt *Options, fn func(result []*collection.Result) bool) error {

	if opt.Language == "" {
		opt.Language = play.DefaultLanguage
	}

	if opt.Country == "" {
		opt.Country = play.DefaultCountry
	}

	reqUrl := strings.Join([]string{play.BasicURL, "store", "search"}, "/")
	v := url.Values{}
	v.Set("gl", opt.Country)
	v.Set("hl", opt.Language)
	v.Set("q", query)
	v.Set("c", "apps")
	v.Set("price", strconv.Itoa(int(opt.PriceType)))

	reqUrl += "?" + v.Encode()

	logrus.Infof("search url %v", reqUrl)

	u, err := url.Parse(reqUrl)
	if err != nil {
		logrus.Errorf("parse collection url %v error=%v", reqUrl, err)
		return err
	}

	if err := collection.Pages(ctx, u.String(), &collection.Options{Language: opt.Language, Country: opt.Country}, func(result []*collection.Result) bool {
		if next := fn(result); !next {
			return false
		}
		return true
	}); err != nil {
		logrus.Errorf("load collection error=%v", err)
		return err
	}

	return nil
}
