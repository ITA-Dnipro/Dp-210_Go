package doctor

import (
	"context"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	ds "github.com/ITA-Dnipro/Dp-210_Go/appointment/proto/doctors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Client struct {
	client ds.DoctorsServiceClient
	conn   *grpc.ClientConn
}

func NewDoctorClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address,
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return &Client{}, err
	}
	return &Client{client: ds.NewDoctorsServiceClient(conn), conn: conn}, nil
}

func (c *Client) GetByID(ctx context.Context, id uuid.UUID) (entity.Doctor, error) {
	r, err := c.client.GetByID(ctx, &ds.GetByIDReq{DoctorID: id.String()})
	if err != nil {
		return entity.Doctor{}, err
	}

	rd := r.GetDoctor()
	doctorID, err := uuid.Parse(rd.GetDoctorID())
	if err != nil {
		return entity.Doctor{}, err
	}

	d := entity.Doctor{
		ID:         doctorID,
		FirstName:  rd.GetFirstName(),
		LastName:   rd.GetLastName(),
		Speciality: rd.GetSpeciality(),
		StartAt:    rd.GetStartAt().AsTime(),
		EndAt:      rd.GetEndAt().AsTime(),
	}
	//TODO: add validation
	return d, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
