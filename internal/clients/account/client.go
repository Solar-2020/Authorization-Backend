package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	GetUserByUid(userID int) (user models.User, err error)
	GetUserByEmail(email string) (user models.User, err error)
	GetYandexUser(userToken string) (user models.User, err error)
	CreateUser(request models.User) (userID int, err error)
}

type client struct {
	host   string
	secret string
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret}
}

type httpError struct {
	Error string `json:"error"`
}

func (c *client) GetUserByUid(userID int) (user models.User, err error) {
	req, err := http.NewRequest(http.MethodGet, c.host+fmt.Sprintf("/api/internal/account/by-user/%s", strconv.Itoa(userID)), nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response models.User
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return user, errors.New(httpErr.Error)
	default:
		return user, errors.New("Unexpected Server Error")
	}
}

func (c *client) GetUserByEmail(email string) (user models.User, err error) {
	req, err := http.NewRequest(http.MethodGet, c.host+fmt.Sprintf("/api/internal/account/by-email/%s", email), nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response models.User
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return user, errors.New(httpErr.Error)
	default:
		return user, errors.New("Unexpected Server Error")
	}
}

func (c *client) GetYandexUser(userToken string) (user models.User, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod("GET")
	req.Header.Add("Authorization", c.secret)

	tempURI, err := url.ParseRequestURI(c.host)
	if err != nil {
		return user, err
	}

	req.URI().SetScheme("http")
	req.URI().SetHost(tempURI.Host)
	req.URI().SetPath("api/internal/account/yandex/" + userToken)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var response models.User
		err = json.Unmarshal(resp.Body(), &response)
		return response, err
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return
		}
		return user, errors.New(httpErr.Error)
	default:
		return user, errors.New("Unexpected Server Error")
	}
}

type CreateUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Surname   string `json:"surname"`
	AvatarURL string `json:"avatarURL"`
}

func (c *client) CreateUser(request models.User) (userID int, err error) {
	createUser := CreateUser{
		ID:        request.ID,
		Email:     request.Email,
		Name:      request.Name,
		Surname:   request.Surname,
		AvatarURL: request.AvatarURL,
	}
	body, err := json.Marshal(createUser)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, c.host+"/api/internal/account/user", bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response CreateUser
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response.ID, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return userID, errors.New(httpErr.Error)
	default:
		return userID, errors.New("Unexpected Server Error")
	}
}
