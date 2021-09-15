package users //userGRPCClient

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/role"
	us "github.com/ITA-Dnipro/Dp-210_Go/doctor/proto/users"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	client us.UsersServiceClient
	logger *zap.Logger
	conn   *grpc.ClientConn
}

func NewUserClient(cfg config.Config, logger *zap.Logger) (*Client, error) {
	conn, err := grpc.Dial(cfg.UserGRPCClient,
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

	d := entity.User{
		ID:             userID,
		Name:           ru.GetName(),
		Email:          ru.GetEmail(),
		PermissionRole: role.Role(ru.GetRole()),
	}
	//TODO: add validation
	return d, nil
}

func (c *Client) Update(ctx context.Context, u *entity.User) error {
	req := &us.UpdateReq{
		UserID: u.ID.String(),
		Name:   u.Name,
		Email:  u.Email,
		Role:   string(u.PermissionRole),
	}
	_, err := c.client.Update(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}
