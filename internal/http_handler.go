package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"moengage/pkg/moengage/models"
	"net/http"
	"net/url"
	"runtime"
)

type HTTPHandler struct {
	AppId      string
	APIKey     string
	BaseURL    string
	HTTPClient http.Client
}

type QueryParameter struct {
	Name  string
	Value string
}

func (h *HTTPHandler) createReq(
	ctx context.Context,
	method string,
	resourcePath string,
	body io.Reader,
	queryParams []QueryParameter,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", h.BaseURL, resourcePath), body)
	if err != nil {
		return nil, err
	}
	req.Header = h.generateCommonHeaders()
	req.URL.RawQuery = generateQueryParams(queryParams)
	return req, nil
}

func (h *HTTPHandler) generateCommonHeaders() http.Header {
	msg := fmt.Sprintf("%s:%s", h.AppId, h.APIKey)
	header := http.Header{}
	header.Add("MOE-APPKEY", h.AppId)
	header.Add("Authorization", fmt.Sprintf("Basic "+base64.StdEncoding.EncodeToString([]byte(msg))))
	header.Add("Content-Type", "application/json")
	header.Add("Accept", "application/json")
	header.Add("User-Agent", "@wayni/go-sdk/v3"+" go/"+runtime.Version())
	return header
}

func generateQueryParams(params []QueryParameter) string {
	q := url.Values{}
	for _, param := range params {
		if param.Value != "" {
			q.Add(param.Name, param.Value)
		}
	}

	return q.Encode()
}

func (h *HTTPHandler) postRequest(
	ctx context.Context,
	payload *bytes.Buffer,
	respResource interface{},
	reqPath string,
	contentType string,
	queryParams []QueryParameter,
) (respDetails models.ResponseDetails, err error) {
	req, err := h.createReq(ctx, http.MethodPost, reqPath, payload, queryParams)
	if err != nil {
		return respDetails, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, parsedBody, err := h.executeReq(req) //nolint: bodyclose // closed in the method itself
	if err != nil {
		_ = json.Unmarshal(parsedBody, &respDetails.ErrorResponse)
		return respDetails, err
	}
	respDetails.HTTPResponse = *resp

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(parsedBody, &respResource)
	} else {
		_ = json.Unmarshal(parsedBody, &respDetails.ErrorResponse)
		// MMS 4xx/5xx responses use the same response as 2xx responses
		if _, ok := respResource.(*models.AlertSuccessResponse); ok {
			_ = json.Unmarshal(parsedBody, &respResource)
		}
	}
	return respDetails, err
}

func (h *HTTPHandler) PostJSONReq(
	ctx context.Context,
	postResource models.Validatable,
	respResource interface{},
	reqPath string,
) (respDetails models.ResponseDetails, err error) {
	err = postResource.Validate()
	if err != nil {
		return respDetails, err
	}
	payload, err := postResource.Marshal()
	if err != nil {
		return respDetails, err
	}
	return h.postRequest(ctx, payload, respResource, reqPath, "application/json", nil)
}

func (h *HTTPHandler) executeReq(
	req *http.Request,
) (resp *http.Response, respBody []byte, err error) {
	resp, err = h.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	parsedBody, err := io.ReadAll(resp.Body)
	if err != nil {
		parsedBody = nil
	}

	return resp, parsedBody, err
}
