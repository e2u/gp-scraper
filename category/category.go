package category

import (
	"context"
	"errors"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/collection"
	"github.com/e2u/gp-scraper/play/price"
	"github.com/e2u/gp-scraper/play/sort"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
	"strings"
)

const (
	DefaultPageNumber = 199
)

type Options = collection.Options

func Pages(ctx context.Context, opt *Options, fn func(result []*collection.Result) bool) error {

	if opt.Language == "" {
		opt.Language = play.DefaultLanguage
	}

	if opt.Country == "" {
		opt.Country = play.DefaultCountry
	}

	path := []string{play.BasicURL, "store", "apps"}
	switch opt.Sort {
	case sort.Rating:
		path = append(path, "top")
	case sort.Newest:
		path = append(path, "new")
	}
	if opt.Category != "" {
		path = append(path, "category/"+string(opt.Category))
	}
	if opt.PageNumber <= 0 {
		opt.PageNumber = DefaultPageNumber
	}

	reqUrl := strings.Join(path, "/")
	v := url.Values{}
	v.Set("gl", opt.Country)
	v.Set("hl", opt.Language)

	if opt.Age != "" {
		v.Set("age", string(opt.Age))
	}

	reqUrl += "?" + v.Encode()
	_, resp, err := util.HttpGet(ctx, reqUrl, nil)
	if err != nil {
		return err
	}

	logrus.Infof("load collection url=%v", reqUrl)
	data := util.ExtractEmbedData(resp)
	if !data["ds:3"].Get("0.1.0").Exists() && !data["ds:3"].Get("0.1.1").Exists() {
		logrus.Errorf("no result")
		return errors.New("no result")
	}

	var appResult gjson.Result

	switch opt.PriceType {
	case price.Free:
		appResult = data["ds:3"].Get("0.1.0")
	case price.Paid:
		appResult = data["ds:3"].Get("0.1.1")
	default:
		appResult = data["ds:3"].Get("0.1.0")
	}

	collectionUrl := fmt.Sprintf("%s%s&hl=%s&gl=%s", play.BasicURL, appResult.Get("0.3.4.2").String(), opt.Language, opt.Country)
	u, err := url.Parse(collectionUrl)
	if err != nil {
		logrus.Errorf("parse collection url %v error=%v", collectionUrl, err)
		return err
	}

	if err := collection.Pages(ctx, u.String(), opt, func(result []*collection.Result) bool {
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

func extractCollectionKey(b []byte) string {
	r := regexp.MustCompile(`([a-z_]+_[A-Z_]+)`)
	if !r.Match(b) {
		return ""
	}
	return string(r.FindAll(b, 1)[0])
}
