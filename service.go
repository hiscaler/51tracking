package tracking51

import (
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/51tracking-go/config"
	"log"
)

type service struct {
	config     *config.Config // Config
	logger     *log.Logger    // Logger
	httpClient *resty.Client  // HTTP client
}

// API Services
type services struct {
	Account  accountService
	Courier  courierService
	Tracking trackingService
}
