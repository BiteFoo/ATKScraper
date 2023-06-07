package http

import (
	"fmt"
	"time"

	resty "github.com/go-resty/resty/v2"
)

var (
	client *resty.Client
)

func init() {
	// json := jsoniter.ConfigCompatibleWithStandardLibrary

	client = resty.New().
		SetRetryCount(3).
		SetTimeout(60 * time.Second).
		SetRetryMaxWaitTime(60 * time.Second).
		SetRetryAfter(func(c *resty.Client, r *resty.Response) (time.Duration, error) {
			return 0, fmt.Errorf("网络连接失败，尝试超过%v次无法链接", c.RetryCount)
		})

}

func Get(url string) ([]byte, error) {

	resp, err := client.R().EnableTrace().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("网络请求失败，statusCode: %v", resp.StatusCode())
	}

	return resp.Body(), nil
}
