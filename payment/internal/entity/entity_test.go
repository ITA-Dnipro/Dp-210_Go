package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*---TestNewBill---*/
type Test struct {
	in   Bill
	in2  []byte
	out  *Bill
	out2 error
}

var tests = []Test{
	{ // Case 1.
		in: Bill{
			DoctorID:  uuid.MustParse("03d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"),
			PatientID: uuid.MustParse("30d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"),
			Price:     123,
		},
		out: &Bill{
			DoctorID:  uuid.MustParse("03d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"),
			PatientID: uuid.MustParse("30d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"),
			Price:     123,
		},
	},
	{ // Case 2.
		in2: []byte("wrong bill"),
	},
}

func TestNewBill(t *testing.T) {
	for _, test := range tests {
		var billBytes []byte
		if test.in2 == nil {
			billBytes, _ = json.Marshal(test.in)
		} else {
			billBytes = test.in2
		}

		result, err := NewBill(billBytes)
		if err != nil {
			assert.NotNil(t, err)
			return
		}
		assert.EqualValues(t, *test.out, *result)
	}
}
