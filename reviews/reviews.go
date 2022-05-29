package reviews

import (
	"context"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/sort"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/url"
	"sync"
	"time"
)

const (
	DefaultPageNumber = 199
)

var reviewsPool = sync.Pool{
	New: func() interface{} {
		return []*Review{}
	},
}

type Options struct {
	Country    string
	Language   string
	PageNumber int
	Sorting    sort.Sort
}

// Review of app
type Review struct {
	Id        string           `json:"id,omitempty"`
	UserName  string           `json:"user_name"`
	Avatar    string           `json:"avatar,omitempty"`
	Time      time.Time        `json:"time,omitempty"`
	Score     int64            `json:"score"`
	URL       string           `json:"url,omitempty"`
	Text      string           `json:"text,omitempty"`
	ReplyTime time.Time        `json:"reply_time,omitempty"`
	ReplyText string           `json:"reply_text,omitempty"`
	Version   string           `json:"version,omitempty"`
	ThumbsUp  int64            `json:"thumbs_up"`
	Criteria  map[string]int64 `json:"criteria,omitempty"`
}

func Pages(ctx context.Context, appId string, opt *Options, fn func([]*Review) bool) error {
	if opt == nil {
		opt = &Options{
			Country:    play.DefaultCountry,
			Language:   play.DefaultLanguage,
			PageNumber: DefaultPageNumber,
			Sorting:    sort.Helpfulness,
		}
	}

	if opt.Country == "" {
		opt.Country = play.DefaultCountry
	}

	if opt.Language == "" {
		opt.Language = play.DefaultLanguage
	}

	if opt.PageNumber <= 0 || opt.PageNumber > 199 {
		opt.PageNumber = DefaultPageNumber
	}

	var rs []*Review
	var nextToken string
	var err error
nextPage:
	rs, nextToken, err = loadReviews(ctx, appId, opt, nextToken)
	if err != nil {
		logrus.Errorf("load %v reviews error=%v", appId, err)
		return err
	}
	if next := fn(rs); !next {
		return nil
	}
	if nextToken != "" {
		goto nextPage
	}
	return nil
}

func loadReviews(ctx context.Context, appId string, opt *Options, nextToken string) ([]*Review, string, error) {
	payload := url.Values{}

	if nextToken == "" {
		payload.Set("f.req", fmt.Sprintf(`[[["UsvDTd","[null,null,[2,%d,[%d,null,null],null,[]],[\"%s\",7]]",null,"generic"]]]`, opt.Sorting, opt.PageNumber, appId))
	} else {
		payload.Set("f.req", fmt.Sprintf(`[[["UsvDTd","[null,null,[2,%d,[%d,null,\"%s\"],null,[]],[\"%s\",7]]",null,"generic"]]]`, opt.Sorting, opt.PageNumber, nextToken, appId))
	}

	result, err := util.BatchExecute(ctx, payload.Encode(), opt.Country, opt.Language)
	if err != nil {
		logrus.Errorf("load rw batch execute error=%v", err)
		return nil, "", err
	}
	data := gjson.Parse(result.Get("0.2").String())
	nextToken = data.Get("1.1").String()

	rs := reviewsPool.Get().([]*Review)
	defer reviewsPool.Put(rs)

	for _, rw := range data.Get("0").Array() {
		r := &Review{Id: rw.Get("0").String()}
		r.UserName = rw.Get("1.0").String()
		r.Avatar = rw.Get("1.1.3.2").String()
		r.Time = time.Unix(rw.Get("5.0").Int(), rw.Get("5.1").Int())
		r.Score = rw.Get("2").Int()
		r.URL = fmt.Sprintf("%s/store/apps/details?id=%s&reviewId=%s", play.BasicURL, appId, r.Id)
		r.Text = rw.Get("4").String()
		r.ReplyTime = time.Unix(rw.Get("7.2.0").Int(), rw.Get("7.2.1").Int())
		r.ReplyText = rw.Get("7.1").String()
		r.Version = rw.Get("10").String()
		r.ThumbsUp = rw.Get("6").Int()
		r.Criteria = func() map[string]int64 {
			cm := make(map[string]int64)
			rw.Get("12.0").ForEach(func(key, value gjson.Result) bool {
				cm[value.Get("0").String()] = value.Get("1.0").Int()
				return true
			})
			return cm
		}()
		rs = append(rs, r)
	}

	return rs, nextToken, nil
}
