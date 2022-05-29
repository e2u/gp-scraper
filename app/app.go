package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/e2u/gp-scraper/internal/util"
	"github.com/e2u/gp-scraper/play"
	"github.com/e2u/gp-scraper/play/price"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	detailURL = "https://play.google.com/store/apps/details?id="
)

type Options struct {
	Country  string
	Language string
}

type Feature struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type App struct {
	AdSupported              bool                `json:"ad_supported"`
	AndroidVersion           string              `json:"android_version,omitempty"`
	AppId                    string              `json:"app_id,omitempty"`
	Available                bool                `json:"available"`
	ContentRating            string              `json:"content_rating,omitempty"`
	ContentRatingDescription string              `json:"content_rating_description,omitempty"`
	Country                  string              `json:"country,omitempty"`
	Description              string              `json:"description,omitempty"`
	DescriptionHTML          string              `json:"description_html,omitempty"`
	Developer                string              `json:"developer,omitempty"`
	DeveloperAddress         string              `json:"developer_address,omitempty"`
	DeveloperEmail           string              `json:"developer_email,omitempty"`
	DeveloperID              string              `json:"developer_id,omitempty"`
	DeveloperURL             string              `json:"developer_url,omitempty"`
	DeveloperWebsite         string              `json:"developer_website,omitempty"`
	EditorsChoice            bool                `json:"editors_choice"`
	FamilyGenre              string              `json:"family_genre,omitempty"`
	FamilyGenreID            string              `json:"family_genre_id,omitempty"`
	Features                 []*Feature          `json:"features"`
	Free                     bool                `json:"free"`
	Genre                    string              `json:"genre,omitempty"`
	GenreId                  string              `json:"genre_id,omitempty"`
	HeaderImage              string              `json:"header_image,omitempty"`
	IAPOffers                bool                `json:"iap_offers"`
	IAPRange                 string              `json:"iap_range,omitempty"`
	Icon                     string              `json:"icon,omitempty"`
	Installs                 string              `json:"installs,omitempty"`
	InstallsMax              int64               `json:"installs_max"`
	InstallsMin              int64               `json:"installs_min"`
	Language                 string              `json:"language,omitempty"`
	Permissions              map[string][]string `json:"permissions"`
	Price                    *price.Price        `json:"price"`
	PriceText                string              `json:"price_text"`
	PrivacyPolicy            string              `json:"privacy_policy,omitempty"`
	Ratings                  int64               `json:"ratings,omitempty"`
	RatingsHistogram         map[int]int64       `json:"ratings_histogram,omitempty"`
	RecentChanges            string              `json:"recent_changes,omitempty"`
	RecentChangesHTML        string              `json:"recent_changes_html,omitempty"`
	Released                 string              `json:"released,omitempty"`
	ReviewsTotalCount        int64               `json:"reviews_total_count"`
	Score                    float64             `json:"score"`
	ScoreText                string              `json:"score_text"`
	Screenshots              []string            `json:"screenshots,omitempty"`
	SimilarURL               string              `json:"similar_url,omitempty"`
	Size                     string              `json:"size,omitempty"`
	Summary                  string              `json:"summary,omitempty"`
	Title                    string              `json:"title,omitempty"`
	URL                      string              `json:"url,omitempty"`
	Updated                  time.Time           `json:"updated"`
	Version                  string              `json:"version,omitempty"`
	Video                    string              `json:"video,omitempty"`
	VideoImage               string              `json:"video_image,omitempty"`
}

func Detail(ctx context.Context, appId string, opt *Options) (*App, error) {
	if opt == nil {
		opt = &Options{
			Country:  play.DefaultCountry,
			Language: play.DefaultLanguage,
		}
	}
	app := &App{
		AppId:    appId,
		URL:      detailURL + appId,
		Language: opt.Language,
		Country:  opt.Country,
	}
	if app.Language == "" {
		app.Language = play.DefaultLanguage
	}
	if app.Country == "" {
		app.Country = play.DefaultCountry
	}
	if err := app.loadDetails(ctx); err != nil {
		logrus.Errorf("load %v detials error=%v", appId, err)
		return nil, err
	}
	return app, nil
}

func (a *App) loadPermissions(ctx context.Context) error {
	payload := url.Values{}
	payload.Set("f.req", fmt.Sprintf(`[[["xdSrCf","[[null,[\"%s\",7],[]]]",null,"1"]]]`, a.AppId))
	data, err := util.BatchExecute(ctx, payload.Encode(), a.Country, a.Language)
	if err != nil {
		return err
	}
	a.Permissions = make(map[string][]string)
	for _, result := range data.Get("0.2").Array() {
		gjson.Parse(result.String()).Get("0").ForEach(func(key, value gjson.Result) bool {
			permKey := value.Get("0").String()
			value.Get("2").ForEach(func(key, value gjson.Result) bool {
				a.Permissions[permKey] = append(a.Permissions[permKey], value.Get("1").String())
				return true
			})
			return true
		})
	}
	return nil
}

func (a *App) loadDetails(ctx context.Context) error {
	v := url.Values{}
	v.Set("gl", a.Country)
	v.Set("hl", a.Language)
	_, resp, err := util.HttpGet(ctx, a.URL+"&"+v.Encode(), nil)
	if err != nil {
		return err
	}

	data := util.ExtractEmbedData(resp)
	if data == nil {
		return errors.New("extract embed data nil")
	}

	// tb, _ := json.Marshal(data)
	// fmt.Println(string(tb))

	normPath := func(path string) string {
	again:
		path = strings.ReplaceAll(path, " ", "")
		path = strings.ReplaceAll(path, ",", ".")
		if strings.Contains(path, " ") {
			goto again
		}
		return path
	}

	getVal := func(key, path string) (*gjson.Result, bool) {
		path = normPath(path)
		if gv, ok := data[key]; ok && gv.Get(path).Exists() {
			val := gv.Get(path)
			return &val, true
		}
		return nil, false
	}
	getStringVal := func(key, path string) string {
		if gv, ok := getVal(key, path); ok {
			return gv.String()
		}
		return ""
	}

	getIntVal := func(key, path string) int64 {
		if gv, ok := getVal(key, path); ok {
			return gv.Int()
		}
		return 0
	}

	getFloatVal := func(key, path string) float64 {
		if gv, ok := getVal(key, path); ok {
			return gv.Float()
		}
		return 0
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.loadPermissions(ctx); err != nil {
			logrus.Errorf("load %v permissions error=%v", a.AppId, err)
		}
	}()
	// https://github.com/facundoolano/google-play-scraper/blob/dev/lib/mappers/details.js
	a.AdSupported = getStringVal("ds:4", "1, 2, 48") != ""
	a.AndroidVersion = getStringVal("ds:4", "1, 2, 140, 1, 1, 0, 0, 1")
	// a.Available = getIntVal("ds:6", "0.12.11.0") != 0 // XX
	a.ContentRating = getStringVal("ds:4", "1, 2, 9, 0")
	a.ContentRatingDescription = getStringVal("ds:4", "1, 2, 9, 2, 1")
	a.DescriptionHTML = getStringVal("ds:4", "1, 2, 72, 0, 1")
	a.Description = util.HTMLToText(a.DescriptionHTML)
	a.Developer = getStringVal("ds:4", "1, 2, 68, 0")
	a.DeveloperAddress = getStringVal("ds:4", "1, 2, 69, 2, 0")
	a.DeveloperEmail = getStringVal("ds:4", "1, 2, 69, 1, 0")
	//developerUrl := getStringVal("ds:6", "0.12.5.5.4.2") // XX
	a.DeveloperID = getStringVal("ds:4", "1, 2, 68, 1, 4, 2")
	// a.DeveloperURL = play.BasicURL + developerUrl
	a.DeveloperWebsite = getStringVal("ds:6", "1, 2, 69, 0, 5, 2")
	a.EditorsChoice = getStringVal("ds:4", "0, 12, 15, 0") != ""
	// a.FamilyGenre = getStringVal("ds:4", "0.12.13.1.0")   // XXX
	// a.FamilyGenreID = getStringVal("ds:6", "0.12.13.1.2") // XXX
	a.Features = func() []*Feature {
		var rs []*Feature
		vv, ok := getVal("ds:6", "0.12.16.2")
		if !ok {
			return rs
		}

		vv.ForEach(func(k gjson.Result, v gjson.Result) bool {
			rs = append(rs, &Feature{
				Title:       v.Get("0").String(),
				Description: v.Get("1.0.0.1").String(),
			})
			return true
		})
		return rs
	}()
	a.Free = getIntVal("ds:4", "1, 2, 57, 0, 0, 0, 0, 1, 0, 0") == 0
	a.Genre = getStringVal("ds:4", "1, 2, 79, 0, 0, 0")
	a.GenreId = getStringVal("ds:4", "1, 2, 79, 0, 0, 2")
	a.HeaderImage = getStringVal("ds:4", "1, 2, 96, 0, 3, 2")
	a.IAPRange = getStringVal("ds:6", "0.12.12.0")
	a.IAPOffers = a.IAPRange != ""
	a.Icon = getStringVal("ds:4", "1, 2, 95, 0, 3, 2")
	a.Installs = getStringVal("ds:4", "1, 2, 13, 0")
	a.InstallsMin = getIntVal("ds:4", "1, 2, 13, 1")
	a.InstallsMax = getIntVal("ds:4", "1, 2, 13, 2")
	a.Price = &price.Price{
		Currency: getStringVal("ds:4", "1, 2, 57, 0, 0, 0, 0, 1, 0, 1"),
		Value:    getFloatVal("ds:4", "1, 2, 57, 0, 0, 0, 0, 1, 0, 0"), // / 1000000.00,
	}
	a.PriceText = getStringVal("ds:4", "1, 2, 19, 0")
	a.PrivacyPolicy = getStringVal("ds:4", "1, 2, 99, 0, 5, 2")
	a.Ratings = getIntVal("ds:4", "1, 2, 51, 2, 1")
	a.RatingsHistogram = func() map[int]int64 {
		rm := make(map[int]int64)
		mv, ok := getVal("ds:4", "1, 2, 51, 1")
		if !ok {
			return rm
		}
		rm[1] = mv.Get("1.1").Int()
		rm[2] = mv.Get("2.1").Int()
		rm[3] = mv.Get("3.1").Int()
		rm[4] = mv.Get("4.1").Int()
		rm[5] = mv.Get("5.1").Int()
		return rm
	}()
	a.RecentChangesHTML = getStringVal("ds:4", "1, 2, 144, 1, 1")
	a.RecentChanges = util.HTMLToText(a.RecentChangesHTML)
	a.Released = getStringVal("ds:4", "1, 2, 10, 0")
	a.ReviewsTotalCount = getIntVal("ds:4", "1, 2, 51, 3, 1")
	a.Score = getFloatVal("ds:4", "1, 2, 51, 0, 1")
	a.ScoreText = getStringVal("ds:4", "1, 2, 51, 0, 0")
	a.Screenshots = func() []string {
		var rs []string
		vv, ok := getVal("ds:4", "1, 2, 78, 0")
		if !ok {
			return rs
		}
		vv.ForEach(func(k gjson.Result, v gjson.Result) bool {
			rs = append(rs, v.Get("3.2").String())
			return true
		})
		return rs
	}()
	// a.Size = getStringVal("ds:8", "0")
	a.Summary = getStringVal("ds:4", "1, 2, 73, 0, 1")
	a.Title = getStringVal("ds:4", "1, 2, 0, 0")
	a.Updated = time.Unix(getIntVal("ds:4", "1, 2, 145, 0, 1, 0"), 0)
	a.Version = getStringVal("ds:4", "1, 2, 140, 0, 0, 0") // Varies with device
	a.Video = getStringVal("ds:4", "1, 2, 100, 0, 0, 3, 2")
	a.VideoImage = getStringVal("ds:4", "1, 2, 100, 1, 0, 3, 2")

	//if similarURL := getStringVal("ds:8", "1.1.0.0.3.4.2"); similarURL != "" {
	//	a.SimilarURL = play.BasicURL + similarURL
	//}

	return nil
}
