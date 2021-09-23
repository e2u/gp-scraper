package similar

import (
	"context"
	"github.com/e2u/gp-scraper/app"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/collection"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Language string
	Country  string
}

func Pages(ctx context.Context, appId string, opt *Options, fn func(result []*collection.Result) bool) error {
	if opt.Language == "" {
		opt.Language = play.DefaultLanguage
	}

	if opt.Country == "" {
		opt.Country = play.DefaultCountry
	}

	a, err := app.Detail(ctx, appId, &app.Options{
		Country:  opt.Country,
		Language: opt.Language,
	})
	if err != nil {
		logrus.Errorf("load app %v detail error=%v", appId, err)
		return err
	}

	if a.SimilarURL == "" {
		logrus.Infof("app %v no huave simailar apps", appId)
		return nil
	}

	if err := collection.Pages(ctx, a.SimilarURL, &collection.Options{Language: opt.Language, Country: opt.Country}, func(result []*collection.Result) bool {
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
