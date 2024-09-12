package external

import (
	"time"

	"github.com/go-resty/resty/v2"
)

var RestyClient *resty.Client

func InitResty() *resty.Client {
	if RestyClient == nil {
		RestyClient = resty.New()

		// global configurations for the RestyClient here
		RestyClient.SetHeader("Content-Type", "application/json")
		// global timeout
		RestyClient.SetTimeout(time.Duration(30) * time.Second)
	}
	return RestyClient
}
