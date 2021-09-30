package stathand

import (
	"encoding/json"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*---TestGetBest---*/
type Test struct {
	in  []entity.Doctor
	out []entity.Doctor
}

var tests = []Test{
	{ // Case 1.
		in: []entity.Doctor{
			{DoctorId: uuid.MustParse("01d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 250},
			{DoctorId: uuid.MustParse("02d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 350},
			{DoctorId: uuid.MustParse("03d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 550},
		},
		out: []entity.Doctor{
			{DoctorId: uuid.MustParse("03d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 550},
		},
	},
	{ // Case 2.
		in: []entity.Doctor{
			{DoctorId: uuid.MustParse("01d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 250},
		},
		out: []entity.Doctor{
			{DoctorId: uuid.MustParse("01d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 250},
		},
	},
}

func TestGetBest(t *testing.T) {
	h := NewHandler(nil)
	for _, test := range tests {
		result := h.GetBest(test.in)
		assert.EqualValues(t, test.out, result)
	}
}

/*---TestDoctorsUnmarshal---*/
type Test2 struct {
	in  []entity.Doctor
	in2 []byte
	out []entity.Doctor
}

var tests2 = []Test2{
	{ // Case 1.
		in: []entity.Doctor{
			{
				DoctorId: uuid.MustParse("03d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 550,
			},
		},
		out: []entity.Doctor{
			{
				DoctorId: uuid.MustParse("03d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 550,
			},
		},
	},
	{ // Case 2.
		in2: []byte("wrong input"),
	},
}

func TestDoctorsUnmarshal(t *testing.T) {
	h := NewHandler(nil)
	for _, test := range tests2 {
		var doctorsBytes []byte
		if test.in2 == nil {
			doctorsBytes, _ = json.Marshal(test.in)
		} else {
			doctorsBytes = test.in2
		}
		result, err := h.DoctorsUnmarshal(doctorsBytes)
		if err != nil {
			assert.NotNil(t, err)
			return
		}
		assert.EqualValues(t, test.out, result)
	}
}
