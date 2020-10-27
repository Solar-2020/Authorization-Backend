package api

import (
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	service "github.com/Solar-2020/GoUtils/http"
)

type AuthServiceInterface interface {
	ByCookie(cookie string, headers map[string]string) (int, error)
}

type AuthClient struct {
	service.Service
	Addr string
}
func (c *AuthClient) Address () string { return c.Addr }
func (c *AuthClient) ByCookie(cookie string, headers map[string]string) (uid int, err error) {
	endpoint := service.ServiceEndpoint{
		Service:   	 c,
		Endpoint:    "/auth/cookie",
		Method:      "POST",
		//ContentType: "application/json",
	}
	message := models.CheckAuthRequest{
		SessionToken: cookie,
	}
	resp := models.CheckAuthResponse{}
	err = endpoint.Send(message, &resp)
	if err != nil {
		return
	}
	uid = resp.Uid
	return
}
