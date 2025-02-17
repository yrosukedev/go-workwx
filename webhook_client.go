package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

// WebhookClient 群机器人客户端
type WebhookClient struct {
	opts options

	key string
}

// NewWebhookClient 构造一个群机器人客户端对象，需要提供 webhook 的 key。
func NewWebhookClient(key string, opts ...CtorOption) *WebhookClient {
	optionsObj := defaultOptions()

	for _, o := range opts {
		o.applyTo(&optionsObj)
	}

	return &WebhookClient{
		opts: optionsObj,

		key: key,
	}
}

// Key 返回该群机器人客户端所配置的 webhook key。
func (c *WebhookClient) Key() string {
	return c.key
}

func (c *WebhookClient) composeQyapiURLWithKey(path string, req interface{}) *url.URL {
	values := url.Values{}
	if valuer, ok := req.(urlValuer); ok {
		values = valuer.intoURLValues()
	}

	// add webhook key
	values.Set("key", c.key)

	// TODO: refactor
	base, err := url.Parse(c.opts.QYAPIHost)
	if err != nil {
		// TODO: error_chain
		panic(fmt.Sprintf("qyapiHost invalid: host=%s err=%+v", c.opts.QYAPIHost, err))
	}

	base.Path = path
	base.RawQuery = values.Encode()

	return base
}

func (c *WebhookClient) executeQyapiJSONPost(path string, req interface{}, respObj interface{}) error {
	url := c.composeQyapiURLWithKey(path, req)
	urlStr := url.String()

	body, err := json.Marshal(req)
	if err != nil {
		// TODO: error_chain
		return err
	}

	resp, err := c.opts.HTTP.Post(urlStr, "application/json", bytes.NewReader(body))
	if err != nil {
		// TODO: error_chain
		return err
	}
	defer resp.Body.Close()

	if respObj != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(respObj)
		if err != nil {
			// TODO: error_chain
			return err
		}
	}

	return nil
}
