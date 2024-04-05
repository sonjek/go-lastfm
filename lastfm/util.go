package lastfm

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func requireAuth(params *apiParams) (err error) {
	if params.sk == "" {
		err = newLibError(
			ErrorAuthRequired,
			Messages[ErrorAuthRequired],
		)
	}
	return
}

func constructUrl(base string, params url.Values) (uri string) {
	p := params.Encode()
	uri = base + "?" + p
	return
}

func toString(val interface{}) (str string, err error) {
	switch v := val.(type) {
	case string:
		str = v
	case int:
		str = strconv.Itoa(v)
	case int64:
		str = strconv.FormatInt(v, 10)
	case []string:
		ss := v
		if len(ss) > 10 {
			ss = ss[:10]
		}
		str = strings.Join(ss, ",")
	default:
		err = newLibError(
			ErrorInvalidTypeOfArgument,
			Messages[ErrorInvalidTypeOfArgument],
		)
	}
	return
}

func parseResponse(body []byte, result interface{}) (err error) {
	var base Base
	err = xml.Unmarshal(body, &base)
	if err != nil {
		return
	}
	if base.Status == ApiResponseStatusFailed {
		var errorDetail ApiError
		err = xml.Unmarshal(base.Inner, &errorDetail)
		if err != nil {
			return
		}
		err = newApiError(&errorDetail)
		return
	} else if result == nil {
		return
	}
	err = xml.Unmarshal(base.Inner, result)
	return
}

func getSignature(params map[string]string, secret string) (sig string) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sigPlain string
	for _, k := range keys {
		sigPlain += k + params[k]
	}
	sigPlain += secret

	hasher := md5.New()
	hasher.Write([]byte(sigPlain))
	sig = hex.EncodeToString(hasher.Sum(nil))
	return
}

func formatArgs(args, rules P) (result map[string]string, err error) {
	result = make(map[string]string)
	if _, ok := rules["indexing"]; ok {
		for _, p := range rules["indexing"].([]string) {
			if valI, ok := args[p]; ok {
				switch v := valI.(type) {
				case string:
					key := p + "[0]"
					result[key] = v
				case int:
					key := p + "[0]"
					result[key] = strconv.Itoa(v)
				case int64: // timestamp
					key := p + "[0]"
					result[key] = strconv.FormatInt(v, 10)
				case []string: // with indexes
					for i, val := range v {
						key := fmt.Sprintf("%s[%d]", p, i)
						result[key] = val
					}
				default:
					err = newLibError(
						ErrorInvalidTypeOfArgument,
						Messages[ErrorInvalidTypeOfArgument],
					)
				}

				if err != nil {
					break
				}
			} else if _, ok := args[p+"[0]"]; ok {
				for i := 0; ; i++ {
					key := fmt.Sprintf("%s[%d]", p, i)
					if valI, ok := args[key]; ok {
						var val string
						val, err = toString(valI)
						result[key] = val
					}

					if err != nil {
						break
					}
				}
			}

			if err != nil {
				break
			}
		}
	}

	if err != nil {
		return
	}

	if _, ok := rules["plain"]; ok {
		for _, key := range rules["plain"].([]string) {
			if valI, ok := args[key]; ok {
				var val string
				val, err = toString(valI)
				result[key] = val
			}

			if err != nil {
				break
			}
		}
	}

	return
}

// ///////////
// GET API  //
// ///////////
func callGet(apiMethod string, params *apiParams, args map[string]interface{}, result interface{}, rules P) (err error) {
	urlParams := url.Values{}
	urlParams.Add("method", apiMethod)
	urlParams.Add("api_key", params.apikey)

	formated, err := formatArgs(args, rules)
	if err != nil {
		return
	}
	for k, v := range formated {
		urlParams.Add(k, v)
	}

	uri := constructUrl(UriApiSecBase, urlParams)

	client := http.DefaultClient
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	if params.useragent != "" {
		req.Header.Set("User-Agent", params.useragent)
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode/100 == 5 { // only 5xx class errors
		err = newLibError(res.StatusCode, res.Status)
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = parseResponse(body, result)
	return
}

// ////////////
// POST API  //
// ////////////
func callPost(apiMethod string, params *apiParams, args P, result interface{}, rules P) (err error) {
	if err = requireAuth(params); err != nil {
		return
	}

	urlParams := url.Values{}
	uri := constructUrl(UriApiSecBase, urlParams)

	// post data
	postData := url.Values{}
	postData.Add("method", apiMethod)
	postData.Add("api_key", params.apikey)
	postData.Add("sk", params.sk)

	tmp := make(map[string]string)
	tmp["method"] = apiMethod
	tmp["api_key"] = params.apikey
	tmp["sk"] = params.sk

	formated, err := formatArgs(args, rules)
	if err != nil {
		return
	}

	for k, v := range formated {
		tmp[k] = v
		postData.Add(k, v)
	}

	sig := getSignature(tmp, params.secret)
	postData.Add("api_sig", sig)

	client := http.DefaultClient
	req, err := http.NewRequest("POST", uri, strings.NewReader(postData.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if params.useragent != "" {
		req.Header.Set("User-Agent", params.useragent)
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = parseResponse(body, result)
	return
}

func callPostWithoutSession(apiMethod string, params *apiParams, args P, result interface{}, rules P) (err error) {
	urlParams := url.Values{}
	uri := constructUrl(UriApiSecBase, urlParams)

	// post data
	postData := url.Values{}
	postData.Add("method", apiMethod)
	postData.Add("api_key", params.apikey)

	tmp := make(map[string]string)
	tmp["method"] = apiMethod
	tmp["api_key"] = params.apikey

	formated, err := formatArgs(args, rules)
	if err != nil {
		return
	}

	for k, v := range formated {
		tmp[k] = v
		postData.Add(k, v)
	}

	sig := getSignature(tmp, params.secret)
	postData.Add("api_sig", sig)

	// call API
	res, err := http.PostForm(uri, postData)
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = parseResponse(body, result)
	return
}
