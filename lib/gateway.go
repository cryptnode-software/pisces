package lib

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cryptnode-software/pisces/lib/errors"
	proto "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
)

// NewGateway is the going to return a gateway i.e. "controller"
func NewGateway(env *Env, services *Services) (*Gateway, error) {

	if services.AuthService == nil {
		return nil, errors.ErrNoAuthService
	}

	if services.PaypalService == nil {
		return nil, errors.ErrNoPaypalService
	}

	if services.OrderService == nil {
		return nil, errors.ErrNoOrderService
	}

	if services.ProductService == nil {
		return nil, errors.ErrNoProductService
	}

	if services.CartService == nil {
		return nil, errors.ErrNoCartService
	}

	return &Gateway{
		services: services,
		Env:      env,
	}, nil
}

//Gateway represents the gateway structure that accepts requests
type Gateway struct {
	proto.UnimplementedPiscesServer
	services *Services
	Env      *Env
}

//SaveOrder creates an order that
func (g *Gateway) SaveOrder(ctx context.Context, req *proto.SaveOrderRequest) (res *proto.SaveOrderResponse, err error) {
	res = new(proto.SaveOrderResponse)
	order := convertOrder(req.Order)

	if order.InquiryID == 0 {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": "an inquiry is required to make a request please include the inquiry id",
			},
		}
	}

	order, err = g.services.OrderService.SaveOrder(ctx, order)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return nil, err
	}

	total, err := g.services.GetTotal(ctx, order)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return
	}
	order.Total = total

	res.Order = convertOrderToProto(order)
	return
}

//SaveCart saves the provided cart and
func (g *Gateway) SaveCart(ctx context.Context, req *proto.SaveCartRequest) (res *proto.SaveCartResponse, err error) {

	if req.Cart == nil || len(req.Cart.Contents) <= 0 {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": "no cart was provided, you can't save an empty or nil cart",
			},
		}
	}

	if req.Cart.OrderId == 0 {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": "no order id was provided, in order for a cart to be properly created it needs to have an associated order id",
			},
		}
	}

	cart, err := g.services.CartService.SaveCart(ctx, convertCart(req.Cart))

	return &proto.SaveCartResponse{
		Cart: convertCartToProto(cart),
	}, nil
}

//AuthorizeOrder handles the authorization of the provided order
func (g *Gateway) AuthorizeOrder(ctx context.Context, req *proto.AuthorizeOrderRequest) (res *proto.AuthorizeOrderResponse, err error) {
	order, err := g.services.OrderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order not found w/ order id %d", req.OrderId),
			},
		}
	}

	if order.Status != OrderStatusUserPending {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order w/ the order id %d isn't in a state of user_pending", req.OrderId),
			},
		}
	}

	cart, err := g.services.CartService.GetCart(ctx, order)

	if cart == nil || len(cart.Contents) <= 0 {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order w/ the order id %d doesn't have a cart associated with it, please create one", req.OrderId),
			},
		}
	}

	if order.ExtID != "" {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order w/ the order id %d already has an external id and has already been associated with a purchase", req.OrderId),
			},
		}
	}

	total, err := g.services.GetTotal(ctx, order)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return
	}
	order.Total = total

	switch order.PaymentMethod {
	case PaymentMethodPaypal:
		order, err = g.services.PaypalService.CreateOrder(ctx, order)
		if err != nil {
			g.Env.Log.Error(err.Error())
			return nil, err
		}
	}

	order.Status = OrderStatusAdminPending

	order, err = g.services.OrderService.SaveOrder(ctx, order)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return nil, err
	}

	res = &proto.AuthorizeOrderResponse{
		Order: convertOrderToProto(order),
	}

	return
}

//CreateUser route...
func (g *Gateway) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.JWT, error) {
	logger := g.Env.Log.With("function", "CreateUser", "ctx", ctx)

	user := convertUserFromProto(req.User)

	user, err := g.services.AuthService.CreateUser(ctx, user, req.Password)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	jwt, err := g.services.AuthService.GenerateJWT(ctx, user)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &proto.JWT{
		Jwt: jwt,
	}, nil
}

//Login route...
func (g *Gateway) Login(ctx context.Context, req *proto.LoginRequest) (*proto.JWT, error) {
	request := &LoginRequest{
		Email:    req.Username,
		Password: req.Password,
	}

	user, err := g.services.AuthService.Login(ctx, request)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return nil, err
	}

	token, err := g.services.AuthService.GenerateJWT(ctx, user)

	return &proto.JWT{
		Jwt: token,
	}, nil
}

//SaveInquiry creates an inquiry requests to a provided destination
func (g *Gateway) SaveInquiry(ctx context.Context, req *proto.Inquiry) (*proto.Inquiry, error) {

	inquiry := convertInquiry(req)

	inquiry, err := g.services.OrderService.SaveInquiry(ctx, inquiry)

	if err != nil {
		g.Env.Log.Error(err.Error())
		return nil, err
	}

	return convertInquiryToProto(inquiry), nil
}

//GetOrders gathers all of the orders and sorts them depending on the request received
func (g *Gateway) GetOrders(ctx context.Context, req *proto.GetOrdersRequest) (res *proto.GetOrdersResponse, err error) {

	if req.OrderId != 0 {
		order, err := g.services.OrderService.GetOrder(ctx, req.OrderId)
		if err != nil {
			return nil, err
		}

		total, err := g.services.GetTotal(ctx, order)
		if err != nil {
			return nil, err
		}

		order.Total = total

		res = &proto.GetOrdersResponse{
			Order: convertOrderToProto(order),
		}

		return res, nil
	}

	//everything beyond this is admin only
	_, err = g.AuthenticateAdmin(ctx)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return
	}

	conditions := &OrderConditions{
		Status: convertOrderStatus(req.Status),
		SortBy: OrdersSortByDueDescending,
	}

	orders, err := g.services.OrderService.GetOrders(ctx, conditions)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return nil, err
	}

	for i, order := range orders {
		total, err := g.services.GetTotal(ctx, order)

		if err != nil {
			g.Env.Log.Error(err.Error())
			return nil, err
		}

		orders[i].Total = total
	}

	res = &proto.GetOrdersResponse{
		Orders: convertOrdersToProto(orders),
	}

	return
}

//GetInquires gathers all of the inquires based off the conditions that are provided through
//the original rpc call
func (g *Gateway) GetInquires(ctx context.Context, req *proto.GetInquiresRequest) (res *proto.GetInquiresResponse, err error) {

	if req.InquiryId != 0 {

		inqury, err := g.services.OrderService.GetInquiry(ctx, req.InquiryId)
		if err != nil {
			return nil, err
		}

		return &proto.GetInquiresResponse{
			Inquiry: convertInquiryToProto(inqury),
		}, nil

	}

	//everything beyond this is admin only
	_, err = g.AuthenticateAdmin(ctx)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return
	}

	inquires, err := g.services.OrderService.GetInquires(ctx, &GetInquiryConditions{
		WithoutOrder: req.WithoutOrder,
	})

	if err != nil {
		return
	}

	return &proto.GetInquiresResponse{
		Inquires: convertInquiresToProto(inquires),
	}, nil
}

func (g *Gateway) GetSignedURL(ctx context.Context, req *proto.GetSignedURLRequest) (res *proto.GetSignedURLResponse, err error) {

	res = new(proto.GetSignedURLResponse)

	{
		client := s3.NewPresignClient(g.services.S3Client,
			s3.WithPresignExpires(3600),
		)

		output, err := g.services.S3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
			ACL:    types.ObjectCannedACLPublicRead,
			Bucket: &g.Env.AWSEnv.Bucket,
			Key:    &req.FileName,
		})

		if err != nil {
			return nil, err
		}

		req, err := client.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:   &g.Env.AWSEnv.Bucket,
			UploadId: output.UploadId,
			Key:      &req.FileName,
		})

		if err != nil {
			return nil, err
		}

		res.Url = req.URL
	}

	return
}

//GeneratePaypalClientToken generates a returns a unique paypal client token in order to create
func (g *Gateway) GeneratePaypalClientToken(ctx context.Context, req *proto.GeneratePaypalClientTokenRequest) (*proto.GeneratePaypalClientTokenResponse, error) {
	token, err := g.services.PaypalService.GenerateClientToken(ctx)
	if err != nil {
		g.Env.Log.Error(err.Error())
		return nil, err
	}

	return &proto.GeneratePaypalClientTokenResponse{
		Token: token.Token,
	}, nil
}

//CheckJWT checks to see if a jwt token is valid and whether or not it has been tampered
//with the method that this uses `ValidateJWT` within the auth  service is one that will
//be used to
func (g *Gateway) CheckJWT(ctx context.Context, req *proto.JWT) (*proto.JWT, error) {
	if req.Jwt == "" {
		jwt, err := GetAuthFromContext(ctx)
		if err != nil {
			return nil, err
		}
		if jwt == "" {
			return nil, errors.ErrNoMetadata
		}
	}

	_, err := g.services.AuthService.DecodeJWT(ctx, req.Jwt)
	if err != nil {
		return nil, err
	}

	return req, nil
}

//AuthenticateAdmin is a export by pass to allow us to directly communicate
//with the auth service from out of the base Pisces library. We need this
//for every route the is considered an admin only route. The `auth` header
// that is signed by a valid user must be set with a valid JWT token in
// order to be approved for any route requires this check.
func (g *Gateway) AuthenticateAdmin(ctx context.Context) (*User, error) {
	return g.services.AuthService.AuthenticateAdmin(ctx)
}

//AuthenticateToken is a export by pass to allow us to directly communicate
//with the auth service from out of the base Pisces library. We need this
//for every route the is considered an user only route. The `auth` header
//must be set with a valid JWT token in order to be approved for any route
//requires this check.
func (g *Gateway) AuthenticateToken(ctx context.Context) (*User, error) {
	return g.services.AuthService.AuthenticateToken(ctx)
}
