package util

import (
	"context"
	"github.com/easonchen147/foundation/cfg"
	"github.com/go-resty/resty/v2"
	"time"
)

var httpClient *resty.Client

func init() {
	httpClient = resty.New()
	httpClient.SetTimeout(time.Second * time.Duration(cfg.AppConf.HttpTimeout))
	httpClient.SetDebug(cfg.AppConf.IsDevEnv())
}

// GetHttpClient 获取http client 实例
func GetHttpClient(ctx context.Context) *resty.Client {
	return httpClient
}

// Post 快速发起post请求
func Post(ctx context.Context, url string, req interface{}, result interface{}, respErr interface{}) error {
	_, err := httpClient.R().SetBody(req).SetResult(&result).SetError(&respErr).Post(url)
	return err
}

// Get  快速发起get请求
func Get(ctx context.Context, url string, queries map[string]string, result interface{}, respErr interface{}) error {
	_, err := httpClient.R().SetQueryParams(queries).SetResult(&result).SetError(&respErr).Get(url)
	return err
}

// PostWithHeader 快速发起post请求，带headers
func PostWithHeader(ctx context.Context, headers map[string]string, url string, req interface{}, result interface{}, respErr interface{}) error {
	_, err := httpClient.R().SetHeaders(headers).SetBody(req).SetResult(&result).SetError(&respErr).Post(url)
	return err
}

// GetWithHeader  快速发起get请求，带headers
func GetWithHeader(ctx context.Context, headers map[string]string, queries map[string]string, url string, result interface{}, respErr interface{}) error {
	_, err := httpClient.R().SetHeaders(headers).SetQueryParams(queries).SetResult(&result).SetError(&respErr).Get(url)
	return err
}
