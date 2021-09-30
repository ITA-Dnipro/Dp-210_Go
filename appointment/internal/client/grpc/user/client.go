package user

import (
	"context"
	"log"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	us "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/users"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Client struct {
	client us.UsersServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address,
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return &Client{}, err
	}
	return &Client{client: us.NewUsersServiceClient(conn), conn: conn}, nil
}

func (c *Client) GetByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	r, err := c.client.GetByID(ctx, &us.GetByIDReq{UserID: id.String()})
	if err != nil {
		return entity.User{}, err
	}

	ru := r.GetUser()
	userID, err := uuid.Parse(ru.GetUserID())
	if err != nil {
		return entity.User{}, err
	}

	u := entity.User{
		ID:             userID,
		Name:           ru.GetName(),
		Email:          ru.GetEmail(),
		PermissionRole: ru.GetRole(),
	}
	return u, nil
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		log.Println(err)
	}
}
