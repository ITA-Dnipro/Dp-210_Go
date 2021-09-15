package appointments //appointmentsGRPCClient

import (
	"context"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/entity"
	as "github.com/ITA-Dnipro/Dp-210_Go/doctor/proto/appointments"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	client as.AppointmentServiceClient
	logger *zap.Logger
	conn   *grpc.ClientConn
}

func NewAppointmentsClient(cfg config.Config, logger *zap.Logger) (*Client, error) {
	conn, err := grpc.Dial(cfg.GRPCHost, grpc.WithInsecure())
	if err != nil {
		logger.Error("client cant establish connection with server", zap.Error(err))
		return &Client{}, err
	}
	return &Client{client: as.NewAppointmentServiceClient(conn), logger: logger, conn: conn}, nil
}

func (client *Client) Close() error {
	err := client.conn.Close()
	if err != nil {
		client.logger.Error("cant close connection", zap.Error(err))
		return err
	}
	return nil
}

func (c *Client) GetByDoctorID(ctx context.Context, id string, from, till time.Time) ([]entity.Appointment, error) {
	frompb := timestamppb.New(from)
	tillpb := timestamppb.New(till)
	r, err := c.client.GetByDoctorID(ctx, &as.GetByDoctrorIDReq{DoctorID: id, From: frompb, Till: tillpb})
	if err != nil {
		return []entity.Appointment{}, err
	}
	ap := r.GetAppointments()

	var aparr []entity.Appointment
	for _, v := range ap {
		a := entity.Appointment{
			ID:        v.GetAppointmentID(),
			DoctorID:  v.GetDoctorID(),
			PatientID: v.GetPatientID(),
			Reason:    v.GetReason(),
			From:      v.GetFrom().AsTime(),
			To:        v.GetTo().AsTime(),
		}
		aparr = append(aparr, a)
	}
	//TODO: add validation
	return aparr, nil
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func Truncate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
