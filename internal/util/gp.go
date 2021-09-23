package util

import (
	"bytes"
	"context"
	"errors"
	"github.com/e2u/gp-scraper/internal/vars"
	"regexp"

	"fmt"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
)

const (
	batchExecuteUrl = "https://play.google.com/_/PlayStoreUi/data/batchexecute"
)

func BatchExecute(ctx context.Context, payload string, country, language string) (*gjson.Result, error) {

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
		"User-Agent":   vars.DefaultUserAgent,
	}
	body := url.Values{}
	body.Set("authuser", "")
	body.Set("bl", "boq_playuiserver_20190903.08_p0")
	body.Set("gl", country)
	body.Set("hl", language)
	body.Set("soc-app", "121")
	body.Set("soc-platform", "1")
	body.Set("soc-device", "1")
	body.Set("rpcids", "qnKhOb")
	body.Set("_reqid", "1065213")

	reqUrl := batchExecuteUrl + "?" + body.Encode()

	_, resp, err := HttpPost(ctx, reqUrl, headers, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	jBytes := bytes.TrimLeft(resp, ")]}'")

	if !gjson.ValidBytes(jBytes) {
		return nil, errors.New("invalid json")
	}

	gs := gjson.ParseBytes(jBytes)
	switch {
	case gs.Get("0.0").String() == "er" && gs.Get("0.5").String() == "400":
		return nil, fmt.Errorf("bad request")
	case gs.Get("0.2").String() == "" && gs.Get("0.5.2.0.0").Exists():
		return nil, fmt.Errorf(gs.Get("0.5.2.0.0").String())
	}
	return &gs, nil
}

func ExtractEmbedData(html []byte) map[string]*gjson.Result {
	scriptRegex := regexp.MustCompile(`>AF_initDataCallback[\s\S]*?<\/script`)
	keyRegex := regexp.MustCompile(`(ds:\d*?)'`)
	valueRegex := regexp.MustCompile(`data:([\s\S]*?), sideChannel: {}}\);<\/`)

	data := make(map[string]*gjson.Result)
	scripts := scriptRegex.FindAll(html, -1)
	for _, script := range scripts {
		key := keyRegex.FindSubmatch(script)
		value := valueRegex.FindSubmatch(script)
		if len(key) > 1 && len(value) > 1 {
			gs := gjson.ParseBytes(value[1])
			data[string(key[1])] = &gs
		}
	}
	return data
}

