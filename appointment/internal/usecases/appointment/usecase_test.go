package appointmen

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/usecases/appointment/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var doctorID = uuid.New()
var userID = uuid.New()

func TestUsecases_CreateRequest(t *testing.T) {
	type fields struct {
		ar       *mocks.MockAppointmentsRepository
		doctors  *mocks.MockDoctorsClient
		users    *mocks.MockUsersClient
		producer *mocks.MockProducer
	}
	type args struct {
		ctx context.Context
		a   *entity.Appointment
	}
	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "testCase 1 - time in past",
			prepare:   nil,
			args:      args{a: &entity.Appointment{From: time.Time{}}},
			assertion: assert.Error,
		},
		{
			name: "testCase 2 - doctor client error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						entity.Doctor{}, errors.New("test"),
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:     time.Now().Add(20 * time.Minute).UTC(),
					DoctorID: doctorID,
				}},
			assertion: assert.Error,
		},
		{
			name: "testCase 3 - doctor working hours after error",
			prepare: func(f *fields) {
				doctor := entity.Doctor{
					StartAt: time.Now().Add(30 * time.Minute).UTC(),
				}
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						doctor, nil,
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:     time.Now().Add(20 * time.Minute).UTC(),
					DoctorID: doctorID,
				}},
			assertion: assert.Error,
		},
		{
			name: "testCase 4 - user client error",
			prepare: func(f *fields) {
				doctor := entity.Doctor{
					StartAt: time.Now().Add(5 * time.Minute).UTC(),
					EndAt:   time.Now().Add(60 * time.Minute).UTC(),
				}
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						doctor, nil,
					),
					f.users.EXPECT().GetByID(
						gomock.Any(), userID,
					).Return(
						entity.User{}, errors.New("test"),
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:      time.Now().Add(10 * time.Minute).UTC(),
					DoctorID:  doctorID,
					PatientID: userID,
				}},
			assertion: assert.Error,
		},
		{
			name: "testCase 5 - user not a patient",
			prepare: func(f *fields) {
				doctor := entity.Doctor{
					StartAt: time.Now().Add(5 * time.Minute).UTC(),
					EndAt:   time.Now().Add(60 * time.Minute).UTC(),
				}
				user := entity.User{
					PermissionRole: "test",
				}
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						doctor, nil,
					),
					f.users.EXPECT().GetByID(
						gomock.Any(), userID,
					).Return(
						user, nil,
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:      time.Now().Add(10 * time.Minute).UTC(),
					DoctorID:  doctorID,
					PatientID: userID,
				}},
			assertion: assert.Error,
		},
		{
			name: "testCase 6 - get appointments on same time error",
			prepare: func(f *fields) {
				doctor := entity.Doctor{
					StartAt: time.Now().Add(5 * time.Minute).UTC(),
					EndAt:   time.Now().Add(60 * time.Minute).UTC(),
				}
				user := entity.User{
					PermissionRole: "patient",
				}
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						doctor, nil,
					),
					f.users.EXPECT().GetByID(
						gomock.Any(), userID,
					).Return(
						user, nil,
					),
					f.ar.EXPECT().GetByDoctorID(
						gomock.Any(), doctorID, gomock.Any(),
					).Return(
						nil, "", errors.New("test"),
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:      time.Now().Add(10 * time.Minute).UTC(),
					DoctorID:  doctorID,
					PatientID: userID,
				}},
			assertion: assert.Error,
		},
		{
			name: "testCase 7 - have appointment appointments on same time error",
			prepare: func(f *fields) {
				doctor := entity.Doctor{
					StartAt: time.Now().Add(5 * time.Minute).UTC(),
					EndAt:   time.Now().Add(60 * time.Minute).UTC(),
				}
				user := entity.User{
					PermissionRole: "patient",
				}
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						doctor, nil,
					),
					f.users.EXPECT().GetByID(
						gomock.Any(), userID,
					).Return(
						user, nil,
					),
					f.ar.EXPECT().GetByDoctorID(
						gomock.Any(), doctorID, gomock.Any(),
					).Return(
						[]entity.Appointment{{}, {}}, "", nil,
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:      time.Now().Add(10 * time.Minute).UTC(),
					DoctorID:  doctorID,
					PatientID: userID,
				}},
			assertion: assert.Error,
		},
		{
			name: "testCase 7 - have appointment appointments on same time error",
			prepare: func(f *fields) {
				doctor := entity.Doctor{
					StartAt: time.Now().Add(5 * time.Minute).UTC(),
					EndAt:   time.Now().Add(60 * time.Minute).UTC(),
				}
				user := entity.User{
					PermissionRole: "patient",
				}
				gomock.InOrder(
					f.doctors.EXPECT().GetByID(
						gomock.Any(), doctorID,
					).Return(
						doctor, nil,
					),
					f.users.EXPECT().GetByID(
						gomock.Any(), userID,
					).Return(
						user, nil,
					),
					f.ar.EXPECT().GetByDoctorID(
						gomock.Any(), doctorID, gomock.Any(),
					).Return(
						[]entity.Appointment{}, "", nil,
					),
					f.producer.EXPECT().SendAppointment(
						gomock.Any(),
					),
				)
			},
			args: args{
				a: &entity.Appointment{
					From:      time.Now().Add(10 * time.Minute).UTC(),
					DoctorID:  doctorID,
					PatientID: userID,
				}},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ar:       mocks.NewMockAppointmentsRepository(ctrl),
				doctors:  mocks.NewMockDoctorsClient(ctrl),
				users:    mocks.NewMockUsersClient(ctrl),
				producer: mocks.NewMockProducer(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			uc := &Usecases{
				ar:       f.ar,
				doctors:  f.doctors,
				users:    f.users,
				producer: f.producer,
			}
			tt.assertion(t, uc.CreateRequest(tt.args.ctx, tt.args.a))
		})
	}
}

func TestUsecases_CreateFromEvent(t *testing.T) {
	type fields struct {
		ar       *mocks.MockAppointmentsRepository
		doctors  *mocks.MockDoctorsClient
		users    *mocks.MockUsersClient
		producer *mocks.MockProducer
	}
	type args struct {
		payload []byte
	}

	appointment, _ := json.Marshal(entity.Appointment{})

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "testCase 1 - unmarshal fail",
			prepare:   nil,
			args:      args{payload: []byte{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 1 - unmarshal fail",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().Create(gomock.Any(), gomock.Any()),
					f.producer.EXPECT().SendNotification(gomock.Any()),
				)
			},
			args:      args{payload: appointment},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ar:       mocks.NewMockAppointmentsRepository(ctrl),
				doctors:  mocks.NewMockDoctorsClient(ctrl),
				users:    mocks.NewMockUsersClient(ctrl),
				producer: mocks.NewMockProducer(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			uc := &Usecases{
				ar:       f.ar,
				doctors:  f.doctors,
				users:    f.users,
				producer: f.producer,
			}
			tt.assertion(t, uc.CreateFromEvent(tt.args.payload))
		})
	}
}

func TestUsecases_Create(t *testing.T) {
	type fields struct {
		ar       *mocks.MockAppointmentsRepository
		doctors  *mocks.MockDoctorsClient
		users    *mocks.MockUsersClient
		producer *mocks.MockProducer
	}
	type args struct {
		ctx context.Context
		a   *entity.Appointment
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "testCase 1 - create error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().Create(
						gomock.Any(), gomock.Any(),
					).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), a: &entity.Appointment{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 2 - notification error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().Create(
						gomock.Any(), gomock.Any(),
					),
					f.producer.EXPECT().SendNotification(
						gomock.Any(),
					).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), a: &entity.Appointment{}},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ar:       mocks.NewMockAppointmentsRepository(ctrl),
				doctors:  mocks.NewMockDoctorsClient(ctrl),
				users:    mocks.NewMockUsersClient(ctrl),
				producer: mocks.NewMockProducer(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			uc := &Usecases{
				ar:       f.ar,
				doctors:  f.doctors,
				users:    f.users,
				producer: f.producer,
			}
			tt.assertion(t, uc.Create(tt.args.ctx, tt.args.a))
		})
	}
}

func TestUsecases_Delete(t *testing.T) {
	id := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dm := mocks.NewMockAppointmentsRepository(ctrl)
	dm.EXPECT().Delete(gomock.Any(), id)
	uc := &Usecases{
		ar:       dm,
		doctors:  mocks.NewMockDoctorsClient(ctrl),
		users:    mocks.NewMockUsersClient(ctrl),
		producer: mocks.NewMockProducer(ctrl),
	}
	assert.NoError(t, uc.Delete(context.Background(), id))
}

func TestUsecases_Update(t *testing.T) {
	type fields struct {
		ar       *mocks.MockAppointmentsRepository
		doctors  *mocks.MockDoctorsClient
		users    *mocks.MockUsersClient
		producer *mocks.MockProducer
	}
	type args struct {
		ctx context.Context
		a   *entity.Appointment
	}
	appointmen := &entity.Appointment{}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "testCase 1 - create error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().Update(
						gomock.Any(), gomock.Any(),
					).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), a: appointmen},
			assertion: assert.Error,
		},
		{
			name: "testCase 2 - send error",
			prepare: func(f *fields) {
				a := appointmen
				a.To = a.From.Add(30 * time.Minute)
				gomock.InOrder(
					f.ar.EXPECT().Update(
						gomock.Any(), a,
					),
					f.producer.EXPECT().SendNotification(
						a,
					).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), a: &entity.Appointment{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 3 - successful",
			prepare: func(f *fields) {
				a := appointmen
				a.To = a.From.Add(30 * time.Minute)
				gomock.InOrder(
					f.ar.EXPECT().Update(
						gomock.Any(), a,
					),
					f.producer.EXPECT().SendNotification(
						a,
					),
				)
			},
			args:      args{ctx: context.Background(), a: &entity.Appointment{}},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ar:       mocks.NewMockAppointmentsRepository(ctrl),
				doctors:  mocks.NewMockDoctorsClient(ctrl),
				users:    mocks.NewMockUsersClient(ctrl),
				producer: mocks.NewMockProducer(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			uc := &Usecases{
				ar:       f.ar,
				doctors:  f.doctors,
				users:    f.users,
				producer: f.producer,
			}
			tt.assertion(t, uc.Update(tt.args.ctx, tt.args.a))
		})
	}
}

func TestUsecases_SendResult(t *testing.T) {
	type fields struct {
		ar       *mocks.MockAppointmentsRepository
		doctors  *mocks.MockDoctorsClient
		users    *mocks.MockUsersClient
		producer *mocks.MockProducer
	}
	type args struct {
		ctx context.Context
		v   *entity.Visit
	}

	v := &entity.Visit{
		AppointmentID: uuid.New(),
		DoctorID:      uuid.New(),
		PatientID:     uuid.New(),
	}

	a := entity.Appointment{
		ID:        v.AppointmentID,
		DoctorID:  v.DoctorID,
		PatientID: v.PatientID,
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "testCase 1 - appointment not exist",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().GetByID(
						gomock.Any(), gomock.Any(),
					).Return(a, errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), v: &entity.Visit{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 2 - appointment not exist",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().GetByID(
						gomock.Any(), gomock.Any(),
					).Return(a, nil),
					f.ar.EXPECT().Delete(
						gomock.Any(), a.ID,
					).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), v: &entity.Visit{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 3 - send bill error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().GetByID(
						gomock.Any(), gomock.Any(),
					).Return(a, nil),
					f.ar.EXPECT().Delete(
						gomock.Any(), a.ID,
					),
					f.producer.EXPECT().SendBill(
						entity.Bill{DoctorID: v.DoctorID, PatientID: v.PatientID},
					).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), v: &entity.Visit{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 4 - send notification error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().GetByID(
						gomock.Any(), gomock.Any(),
					).Return(a, nil),
					f.ar.EXPECT().Delete(
						gomock.Any(), a.ID,
					),
					f.producer.EXPECT().SendBill(
						entity.Bill{DoctorID: v.DoctorID, PatientID: v.PatientID},
					),
					f.producer.EXPECT().SendNotification(v).Return(errors.New("test")),
				)
			},
			args:      args{ctx: context.Background(), v: &entity.Visit{}},
			assertion: assert.Error,
		},
		{
			name: "testCase 5 - successful",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.ar.EXPECT().GetByID(
						gomock.Any(), gomock.Any(),
					).Return(a, nil),
					f.ar.EXPECT().Delete(
						gomock.Any(), a.ID,
					),
					f.producer.EXPECT().SendBill(
						entity.Bill{DoctorID: v.DoctorID, PatientID: v.PatientID},
					),
					f.producer.EXPECT().SendNotification(v),
				)
			},
			args:      args{ctx: context.Background(), v: &entity.Visit{}},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ar:       mocks.NewMockAppointmentsRepository(ctrl),
				doctors:  mocks.NewMockDoctorsClient(ctrl),
				users:    mocks.NewMockUsersClient(ctrl),
				producer: mocks.NewMockProducer(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			uc := &Usecases{
				ar:       f.ar,
				doctors:  f.doctors,
				users:    f.users,
				producer: f.producer,
			}
			tt.assertion(t, uc.SendResult(tt.args.ctx, tt.args.v))
		})
	}
}

func TestUsecases_GetAll(t *testing.T) {
	p := &entity.AppointmentsParam{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dm := mocks.NewMockAppointmentsRepository(ctrl)
	dm.EXPECT().GetAll(gomock.Any(), p)
	uc := &Usecases{
		ar:       dm,
		doctors:  mocks.NewMockDoctorsClient(ctrl),
		users:    mocks.NewMockUsersClient(ctrl),
		producer: mocks.NewMockProducer(ctrl),
	}
	_, _, err := uc.GetAll(context.Background(), p)
	assert.NoError(t, err)
}

func TestUsecases_GetByDoctorID(t *testing.T) {
	id := uuid.New()
	p := &entity.AppointmentsParam{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dm := mocks.NewMockAppointmentsRepository(ctrl)
	dm.EXPECT().GetByDoctorID(gomock.Any(), id, p)
	uc := &Usecases{
		ar:       dm,
		doctors:  mocks.NewMockDoctorsClient(ctrl),
		users:    mocks.NewMockUsersClient(ctrl),
		producer: mocks.NewMockProducer(ctrl),
	}
	_, _, err := uc.GetByDoctorID(context.Background(), id, p)
	assert.NoError(t, err)
}

func TestUsecases_GetByPatientID(t *testing.T) {
	id := uuid.New()
	p := &entity.AppointmentsParam{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dm := mocks.NewMockAppointmentsRepository(ctrl)
	dm.EXPECT().GetByPatientID(gomock.Any(), id, p)
	uc := &Usecases{
		ar:       dm,
		doctors:  mocks.NewMockDoctorsClient(ctrl),
		users:    mocks.NewMockUsersClient(ctrl),
		producer: mocks.NewMockProducer(ctrl),
	}
	_, _, err := uc.GetByPatientID(context.Background(), id, p)
	assert.NoError(t, err)
}

func TestUsecases_GetByID(t *testing.T) {
	id := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dm := mocks.NewMockAppointmentsRepository(ctrl)
	dm.EXPECT().GetByID(gomock.Any(), id)
	uc := &Usecases{
		ar:       dm,
		doctors:  mocks.NewMockDoctorsClient(ctrl),
		users:    mocks.NewMockUsersClient(ctrl),
		producer: mocks.NewMockProducer(ctrl),
	}
	_, err := uc.GetByID(context.Background(), id)
	assert.NoError(t, err)
}
