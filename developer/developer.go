package developer

import (
	"context"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/collection"
	"github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strings"
)

type Options struct {
	Country  string
	Language string
}

func Pages(ctx context.Context, devId string, opt *Options, fn func(result []*collection.Result) bool) error {
	if opt.Language == "" {
		opt.Language = play.DefaultLanguage
	}

	if opt.Country == "" {
		opt.Country = play.DefaultCountry
	}

	rex := regexp.MustCompile("^[0-9]+")

	ps := []string{play.BasicURL, "store", "apps"}
	if rex.MatchString(devId) {
		ps = append(ps, "dev")
	} else {
		ps = append(ps, "developer")
	}
	reqUrl := strings.Join(ps, "/")

	v := url.Values{}
	v.Set("gl", opt.Country)
	v.Set("hl", opt.Language)
	v.Set("id", devId)

	reqUrl += "?" + v.Encode()

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
