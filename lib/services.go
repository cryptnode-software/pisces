package lib

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//Services ...
type Services struct {
	ProductService ProductService
	UploadService  UploadService
	PaypalService  PaypalService
	OrderService   OrderService
	AuthService    AuthService
	CartService    CartService
	S3Client       *s3.Client
}
