package downloader

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SplashConn struct {
	host            string
	user            string
	password        string
	timeout         int
	resourceTimeout int
	wait            int
	lua             string
}

type SplashResponse struct {
	Num1 struct {
		HeadersSize int    `json:"headersSize"`
		URL         string `json:"url"`
		Ok          bool   `json:"ok"`
		Content     struct {
			Encoding string `json:"encoding"`
			Text     string `json:"text"`
			MimeType string `json:"mimeType"`
			Size     int    `json:"size"`
		} `json:"content"`
		Headers []struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"headers"`
		Cookies     []interface{} `json:"cookies"`
		HTTPVersion string        `json:"httpVersion"`
		Status      int           `json:"status"`
		StatusText  string        `json:"statusText"`
		RedirectURL string        `json:"redirectURL"`
	} `json:"1"`
}

//NewSplashConn opens new connection to Splash Server
func NewSplashConn(host, user, password string, timeout, resourceTimeout, wait int, lua string) SplashConn {

	return SplashConn{
		//	config:     cnf,
		host:            host,
		user:            user,
		password:        password,
		timeout:         timeout,
		resourceTimeout: resourceTimeout,
		wait:            wait,
		lua:             lua,
	}
}

func (s *SplashConn) Download(addr string) ([]byte, error) {
	client := &http.Client{}
	//splashURL := fmt.Sprintf("%s%s?&url=%s&timeout=%d&resource_timeout=%d&wait=%d", s.host, s.renderHTMLURL, url.QueryEscape(addr), s.timeout, s.resourceTimeout, s.wait)

	splashURL := fmt.Sprintf(
		"%sexecute?url=%s&timeout=%d&resource_timeout=%d&wait=%d&lua_source=%s", s.host,
		url.QueryEscape(addr),
		s.timeout,
		s.resourceTimeout,
		s.wait,
		url.QueryEscape(s.lua))

	req, err := http.NewRequest("GET", splashURL, nil)
	req.SetBasicAuth(s.user, s.password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//response from Splash service
	if resp.StatusCode == 200 {
		var sResponse SplashResponse
		if err := json.Unmarshal(res, &sResponse); err != nil {
			logger.Println("Json Unmarshall error", err)
		}
		//status code returned by Splash
		statusCode := sResponse.Num1.Status
		if statusCode != 200 {
			return nil, fmt.Errorf("Error: %d. %s",
				statusCode,
				sResponse.Num1.StatusText)
		}
		decoded, err := base64.StdEncoding.DecodeString(sResponse.Num1.Content.Text)
		if err != nil {
			logger.Println("decode error:", err)
			return nil, fmt.Errorf(string(res))
		}
		return decoded, nil
	}
	return nil, fmt.Errorf(string(res))
}

/*
func (s *SplashConn) Download1(addr string) ([]byte, error) {
	client := &http.Client{}
	splashURL := fmt.Sprintf("%s%s?&url=%s&timeout=%d&resource_timeout=%d&wait=%d", s.host, s.renderHTMLURL, url.QueryEscape(addr), s.timeout, s.resourceTimeout, s.wait)
	req, err := http.NewRequest("GET", splashURL, nil)
	req.SetBasicAuth(s.user, s.password)
	req.Header.Set("Content-Type", "text/plain")
	//fmt.Println(req.Header.Get("Content-Type"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//resp.Header.Set("Content-Type", "text/plain")
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)

	//fmt.Println("CONTENT", http.DetectContentType(res))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		return res, nil
	}
	return nil, fmt.Errorf(string(res))
}
*/