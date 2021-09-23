package collection

import (
	"context"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/age"
	"github.com/e2u/gp-scraper/play/category"
	"github.com/e2u/gp-scraper/play/price"
	"github.com/e2u/gp-scraper/play/sort"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math"
	"net/url"
)

const (
	DefaultPageNumber = 200
)

type Collection string

const (
	TopFree     Collection = "topselling_free"
	TopGrossing Collection = "topgrossing"
	TopNewFree  Collection = "topselling_new_free"
	TopNewPaid  Collection = "topselling_new_paid"
	TopPaid     Collection = "topselling_paid"
	TopTrending Collection = "movers_shakers"
)

type Result struct {
	Country      string       `json:"country,omitempty"`
	Developer    string       `json:"developer,omitempty"`
	DeveloperId  string       `json:"developer_id,omitempty"`
	DeveloperURL string       `json:"developer_url,omitempty"`
	Free         bool         `json:"free"`
	Icon         string       `json:"icon,omitempty"`
	Language     string       `json:"language,omitempty"`
	Price        *price.Price `json:"price,omitempty"`
	PriceText    string       `json:"price_text,omitempty"`
	Score        float64      `json:"score"`
	ScoreText    string       `json:"score_text,omitempty"`
	Summary      string       `json:"summary,omitempty"`
	Title        string       `json:"title,omitempty"`
	URL          string       `json:"url,omitempty"`
}

func parseResultFromGJson(gs gjson.Result, opt *Options) *Result {
	r := &Result{
		Language: opt.Language,
		Country:  opt.Country,
	}
	r.Developer = gs.Get("4.0.0.0").String()
	r.DeveloperURL = play.BasicURL + gs.Get("4.0.0.1.4.2").String()
	r.DeveloperId = func() string {
		u, err := url.Parse(r.DeveloperURL)
		if err != nil {
			logrus.Errorf("get developer id error=%v", err)
			return ""
		}
		return u.Query().Get("id")
	}()
	r.Icon = gs.Get("1.1.0.3.2").String()
	r.PriceText = gs.Get("7.0.3.2.1.0.2").String()
	r.Price = &price.Price{
		Currency: gs.Get("7.0.3.2.1.0.1").String(),
		Value:    gs.Get("7.0.3.2.1.0.0").Float(),
	}
	r.Free = r.Price.Value <= 0
	r.Score = gs.Get("6.0.2.1.1").Float()
	r.ScoreText = gs.Get("6.0.2.1.0").String()
	r.Title = gs.Get("2").String()
	r.Summary = gs.Get("4.1.1.1.1").String()
	r.URL = play.BasicURL + gs.Get("9.4.2").String() + "&gl=" + r.Country + "&hl=" + r.Language

	if opt.PriceMin > 0 && opt.PriceMax > 0 {
		bv := r.Price.Value / float64(1000000)
		if bv >= opt.PriceMin && bv <= opt.PriceMax {
			return r
		}
		return nil
	}

	if opt.ScoreMin > 0 && opt.ScoreMax > 0 {
		if r.Score >= opt.ScoreMin && r.Score <= opt.ScoreMax {
			return r
		}
		return nil
	}

	return r
}

type Options struct {
	sort.Sort
	category.Category
	age.Age
	PriceType  price.Type
	PageNumber int64
	Country    string
	Language   string
	PriceMin   float64
	PriceMax   float64
	ScoreMin   float64
	ScoreMax   float64
}

func Pages(ctx context.Context, reqUrl string, opt *Options, fn func(results []*Result) bool) error {

	if opt.ScoreMax <= 0 {
		opt.ScoreMax = math.MaxFloat64
	}
	if opt.PriceMax <= 0 {
		opt.PriceMax = math.MaxFloat64
	}
	if opt.PageNumber <= 0 || opt.PageNumber > DefaultPageNumber {
		opt.PageNumber = DefaultPageNumber
	}

	logrus.Infof("collection send http get request url=%v", reqUrl)
	_, resp, err := util.HttpGet(ctx, reqUrl, nil)
	if err != nil {
		logrus.Errorf("send request %v error=%v", reqUrl, err)
		return err
	}

	// first page
	var results []*Result
	data := util.ExtractEmbedData(resp)

	data["ds:3"].Get("0.1.0.0.0").ForEach(func(key, value gjson.Result) bool {
		if result := parseResultFromGJson(value, opt); result != nil {
			results = append(results, result)
		}
		return true
	})

	if next := fn(results); !next {
		return nil
	}

	nextToken := data["ds:3"].Get("0.1.0.0.7.1").String()

	if nextToken == "" {
		return nil
	}

nextPage:
	results, nextToken, err = loadMoreCollection(ctx, opt, nextToken)
	if err != nil {
		logrus.Errorf("load collction error=%v", err)
		return err
	}

	if next := fn(results); !next {
		return nil
	}

	if nextToken != "" {
		goto nextPage
	}

	return nil
}

func loadMoreCollection(ctx context.Context, opt *Options, nextToken string) ([]*Result, string, error) {
	payload := url.Values{}
	payload.Set("f.req", fmt.Sprintf(`[[["qnKhOb","[[null,[[10,[10,%d]],true,null,[96,27,4,8,57,30,110,79,11,16,49,1,3,9,12,104,55,56,51,10,34,77]],null,\"%s\"]]",null,"generic"]]]`, opt.PageNumber, nextToken))

	result, err := util.BatchExecute(ctx, payload.Encode(), opt.Country, opt.Language)
	if err != nil {
		logrus.Errorf("load rw batch execute error=%v", err)
		return nil, "", err
	}
	data := gjson.Parse(result.Get("0.2").String())
	nextToken = data.Get("0.0.7.1").String()

	var results []*Result
	data.Get("0.0.0").ForEach(func(key, value gjson.Result) bool {
		if result := parseResultFromGJson(value, opt); result != nil {
			results = append(results, result)
		}
		return true
	})
	return results, nextToken, nil
}
