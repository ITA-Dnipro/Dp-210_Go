package kafkahand

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/proto/statistics"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) InsertDataToDb(*entity.Bill) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) GetAllDataFromDb() (*entity.DocPat, error) {
	args := m.Called()
	return args.Get(0).(*entity.DocPat), args.Error(1)
}

type MockStatClient struct {
	mock.Mock
}

func (m *MockStatClient) DocStat(tx context.Context, in *statistics.DocRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	_ = m.Called()
	return nil, nil
}

/*---TestSendMonthlyReport---*/
type Test struct {
	out         *entity.DocPat
	outError    error
	callDocStat bool
}

var tests = []Test{
	{ // Case 1.
		out: &entity.DocPat{
			Doctors:  []entity.Doctor{{DoctorId: uuid.MustParse("01d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), DoctorTotal: 250}},
			Patients: []entity.Patient{{PatientId: uuid.MustParse("109d9717-b3a4-4a11-b2f1-67df7246cc21"), PatientTotal: 250}},
		},
		outError:    nil,
		callDocStat: true,
	},
	{ // Case 2.
		out: &entity.DocPat{
			Doctors:  []entity.Doctor{},
			Patients: []entity.Patient{},
		},
		outError:    errors.New("some db error"),
		callDocStat: false,
	},
	{ // Case 3.
		out: &entity.DocPat{
			Doctors:  []entity.Doctor{},
			Patients: []entity.Patient{},
		},
		outError:    nil,
		callDocStat: false,
	},
}

func TestSendMonthlyReport(t *testing.T) {
	mockRepository := new(MockRepository)
	mockStatClient := new(MockStatClient)
	testHandler := NewHandler(mockRepository, nil, mockStatClient)

	for _, test := range tests {
		mockRepository.On("GetAllDataFromDb").Once().Return(test.out, test.outError).Once()
		if test.callDocStat {
			mockStatClient.On("DocStat").Once().Return(nil, test.outError).Once()
		}

		result, err := testHandler.SendMonthlyReport(context.Background())
		mockRepository.AssertExpectations(t)
		if test.callDocStat {
			mockStatClient.AssertExpectations(t)
		}

		if err != nil {
			assert.NotNil(t, err)
		} else {
			patientsBytes, _ := json.Marshal(test.out.Patients)
			assert.EqualValues(t, string(patientsBytes), string(result))
		}
	}
}

type Test2 struct {
	in  *entity.Bill
	out error
}

var tests2 = []Test2{
	{ // Case 1.
		in:  &entity.Bill{DoctorID: uuid.MustParse("01d41c0f-f3ca-4e0b-80b0-dd5bd4a0b586"), PatientID: uuid.MustParse("109d9717-b3a4-4a11-b2f1-67df7246cc21"), Price: 250},
		out: nil,
	},
}

func TestInsertToDb(t *testing.T) {
	mockRepository := new(MockRepository)
	testHandler := NewHandler(mockRepository, nil, nil)
	for _, test := range tests2 {
		mockRepository.On("InsertDataToDb").Once().Return(test.out).Once()
		err := testHandler.InsertToDb(test.in)
		mockRepository.AssertExpectations(t)

		assert.Nil(t, err)
	}
}
