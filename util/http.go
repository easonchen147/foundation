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
	httpClient.SetTimeout(time.Second * 5)
	httpClient.SetDebug(cfg.AppConf.IsDevEnv())
}

// GetHttpClient 获取http client 实例
func GetHttpClient(ctx context.Context) *resty.Client {
	return httpClient
}

// Post 快速发起post请求
func Post(ctx context.Context, url string, data interface{}, result interface{}) error {
	_, err := httpClient.R().SetBody(data).SetResult(&result).Post(url)
	return err
}

// PostWithHeader 快速发起post请求，带headers
func PostWithHeader(ctx context.Context, headers map[string]string, url string, data interface{}, result interface{}) error {
	_, err := httpClient.R().SetHeaders(headers).SetBody(data).SetResult(&result).Post(url)
	return err
}

// GetWithHeader  快速发起get请求，带headers
func GetWithHeader(ctx context.Context, headers map[string]string, queries map[string]string, url string, result interface{}) error {
	_, err := httpClient.R().SetHeaders(headers).SetQueryParams(queries).SetResult(&result).Get(url)
	return err
}
