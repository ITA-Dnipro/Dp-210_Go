package doctor

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	ds "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/doctors"
	"google.golang.org/grpc"
)

type Client struct {
	client ds.DoctorsServiceClient
	conn   *grpc.ClientConn
}

func NewDoctorClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return &Client{}, err
	}
	return &Client{client: ds.NewDoctorsServiceClient(conn), conn: conn}, nil
}

func (c *Client) GetByID(ctx context.Context, id string) (entity.Doctor, error) {
	r, err := c.client.GetByID(ctx, &ds.GetByIDReq{DoctorID: id})
	if err != nil {
		return entity.Doctor{}, err
	}
	d := entity.Doctor{
		ID:         r.Doctor.DoctorID,
		FirstName:  r.Doctor.FirstName,
		LastName:   r.Doctor.LastName,
		Speciality: r.Doctor.Speciality,
		StartAt:    r.Doctor.EndAt.AsTime(),
		EndAt:      r.Doctor.EndAt.AsTime(),
	}
	//TODO: add validation
	return d, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
