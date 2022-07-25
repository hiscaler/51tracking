package tracking51

import (
	"github.com/go-resty/resty/v2"
	"log"
)

type service struct {
	config     *config       // Config
	logger     *log.Logger   // Logger
	httpClient *resty.Client // HTTP client
}

// API Services
type services struct {
	Courier  courierService
	Tracking trackingService
}
