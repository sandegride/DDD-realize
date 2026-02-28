package discount

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/internal/core/ports"
	"basket-service/internal/generated/clients/discountsrv/discountpb"
	"basket-service/internal/pkg/errs"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var _ ports.DiscountClient = &client{}

type client struct {
	conn            *grpc.ClientConn
	pbDiscountClint discountpb.DiscountClient
	timeout         time.Duration
}

func NewClient(host string) (ports.DiscountClient, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbDiscountClint := discountpb.NewDiscountClient(conn)

	return &client{
		conn:            conn,
		pbDiscountClint: pbDiscountClint,
		timeout:         5 * time.Second,
	}, nil
}

func (c *client) GetDiscount(ctx context.Context, basket *basket.Basket) (kernel.Discount, error) {
	req := &discountpb.GetDiscountRequest{
		Items: func() []*discountpb.Item {
			result := make([]*discountpb.Item, 0, len(basket.Items()))
			for _, item := range basket.Items() {
				result = append(result, &discountpb.Item{Id: item.ID().String()})
			}
			return result
		}(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.pbDiscountClint.GetDiscount(ctx, req)
	if err != nil {
		return kernel.Discount{}, err
	}

	return kernel.NewDiscount(resp.Value)
}

func (c *client) Close() error {
	return c.conn.Close()
}
