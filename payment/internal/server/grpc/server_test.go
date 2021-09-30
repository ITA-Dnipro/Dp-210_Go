package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/handlers/stathand"
	stat "github.com/ITA-Dnipro/Dp-210_Go/payment/internal/proto/statistics"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

type Test struct {
	in  interface{}
	out error
}

var tests = []Test{
	{ // Case 1.
		in: []entity.Doctor{
			{DoctorId: uuid.MustParse("01d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 250},
			{DoctorId: uuid.MustParse("02d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 1000},
		},
		out: nil,
	},
	{ // Case 2.
		in:  "some error",
		out: errors.New("unmarshal error > json: cannot unmarshal string into Go value of type []entity.Doctor"),
	},
}

func TestDocStat(t *testing.T) {
	logger, _ := zap.NewProduction()
	svr := NewGRPCServer(stathand.NewHandler(logger))
	for _, test := range tests {
		doctorsBytes, _ := json.Marshal(test.in)
		_, err := svr.DocStat(context.Background(), &stat.DocRequest{DocsBytesArr: doctorsBytes})

		if err != nil {
			assert.NotNil(t, err)
			assert.EqualError(t, err, test.out.Error())
			return
		}
		assert.Nil(t, err)
	}
}
