package client

import (
	"encoding/json"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type Client interface {
	GetUserIDByCookie(sessionToken string) (userID int, err error)
	CompareSecret(inputSecret string) (err error)
	AuthorizeRequest(sessionToken string, lifetime int, req fasthttp.Request) (res fasthttp.Request, err error)
}

type client struct {
	host        string
	secret      string
	errorWorker errorWorker.ErrorWorker
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret, errorWorker: errorWorker.NewErrorWorker()}
}

func (c *client) GetUserIDByCookie(sessionToken string) (userID int, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().QueryArgs().Set("session_cookie", sessionToken)
	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/auth/cookie")

	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Authorization", c.secret)


	err = fasthttp.Do(req, resp)
	if err != nil {
		return userID, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var response GetUserIdByCookieResponse
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			return userID, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		}
		return response.UserID, nil
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return userID, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		}
		return userID, c.errorWorker.NewError(fasthttp.StatusBadRequest, errors.New(httpErr.Error), errors.New(httpErr.Error))
	default:
		return userID, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()))
	}
}

func (c *client) CompareSecret(inputSecret string) (err error) {
	if !strings.EqualFold(inputSecret, c.secret) {
		return c.errorWorker.NewError(fasthttp.StatusForbidden, ErrorInvalidSecretKey, ErrorInvalidSecretKey)
	}
	return
}

func (c *client) AuthorizeRequest(sessionToken string, lifetime int, request fasthttp.Request) (res fasthttp.Request, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().QueryArgs().Set("session_cookie", sessionToken)
	req.URI().QueryArgs().Set("lifetime",strconv.Itoa(lifetime))
	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/auth/cookie/dublicate")

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Authorization", c.secret)


	err = fasthttp.Do(req, resp)
	if err != nil {
		return res, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		cookie := resp.Header.PeekCookie("SessionToken")
		if cookie == nil || len(cookie) == 0 {
			return res, c.errorWorker.NewError(fasthttp.StatusInternalServerError, errors.New("empty session key gor"), err)
		}
		request.Header.SetCookie("SessionToken", string(cookie))
		return request, nil
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return res, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		}
		return res, c.errorWorker.NewError(fasthttp.StatusBadRequest, errors.New(httpErr.Error), errors.New(httpErr.Error))
	default:
		return res, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()))
	}
}