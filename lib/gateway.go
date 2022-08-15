package lib

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/google/uuid"
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

	order, err := convertOrder(req.Order)
	if err != nil {
		return nil, err
	}

	if order.Inquiry == nil {
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

	o, err := convertOrderToProto(order)
	if err != nil {
		return
	}

	res.Order = o
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

	if req.Cart.OrderId == "" {
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
	uuid, err := uuid.Parse(req.OrderId)
	if err != nil {
		return
	}

	order, err := g.services.OrderService.GetOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order not found w/ order id %s", req.OrderId),
			},
		}
	}

	if order.Status != OrderStatusUserPending {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order w/ the order id %s isn't in a state of user_pending", req.OrderId),
			},
		}
	}

	cart, err := g.services.CartService.GetCart(ctx, order)

	if cart == nil || len(cart.Contents) <= 0 {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order w/ the order id %s doesn't have a cart associated with it, please create one", req.OrderId),
			},
		}
	}

	if order.ExtID != "" {
		return nil, &errors.ErrInvalidRequest{
			Fields: map[string]string{
				"reason": fmt.Sprintf("order w/ the order id %s already has an external id and has already been associated with a purchase", req.OrderId),
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

	o, err := convertOrderToProto(order)
	if err != nil {
		return nil, err
	}

	res = &proto.AuthorizeOrderResponse{
		Order: o,
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

	if uuid, err := uuid.Parse(req.OrderId); err == nil {
		order, err := g.services.OrderService.GetOrder(ctx, uuid)
		if err != nil {
			return nil, err
		}

		total, err := g.services.GetTotal(ctx, order)
		if err != nil {
			return nil, err
		}

		order.Total = total

		o, err := convertOrderToProto(order)
		if err != nil {
			return nil, err
		}

		res = &proto.GetOrdersResponse{
			Order: o,
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

	o, err := convertOrdersToProto(orders)
	if err != nil {
		return nil, err
	}

	res = &proto.GetOrdersResponse{
		Orders: o,
	}

	return
}

//GetInquires gathers all of the inquires based off the conditions that are provided through
//the original rpc call
func (g *Gateway) GetInquires(ctx context.Context, req *proto.GetInquiresRequest) (res *proto.GetInquiresResponse, err error) {

	if uuid, err := uuid.Parse(req.InquiryId); err == nil {

		inqury, err := g.services.OrderService.GetInquiry(ctx, uuid)
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

func (g *Gateway) StartUpload(ctx context.Context, req *proto.StartUploadRequest) (res *proto.StartUploadResponse, err error) {

	res = new(proto.StartUploadResponse)

	{

		uuid := uuid.New().String() + req.Key

		client := s3.NewPresignClient(g.services.S3Client)

		req, err := client.PresignPutObject(ctx, &s3.PutObjectInput{
			ACL:    types.ObjectCannedACLPublicRead,
			Bucket: &g.Env.AWSEnv.Bucket,
			Key:    &uuid,
		},
		)

		if err != nil {
			return nil, err
		}

		res.PresignedUrl = req.URL
		res.Url = fmt.Sprintf("https://%s.%s.linodeobjects.com", g.Env.AWSEnv.Bucket, g.Env.AWSEnv.Region)

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
