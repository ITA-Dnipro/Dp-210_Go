package entity

import (
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type Patient struct {
	Id         uuid.UUID
	Name       string `json:"name" validate:"required,min=3,max=25"`
	Email      string `json:"email" validate:"required,email,min=5,max=25"`
	Gender     string `json:"gender" validate:"required,eq=Male|eq=Female"`
	Age        int    `json:"age" validate:"required,numeric,gt=0,lte=100"`
	Phone      string `json:"phone" validate:"required,e164"`
	Address    string `json:"address" validate:"required,max=150"`
	Disability bool   `json:"disability"`
	RegAt      time.Time
}

func NewPatient(dataArr []string) (p *Patient, err error) {
	p = &Patient{
		Name:    dataArr[0],
		Email:   dataArr[1],
		Gender:  dataArr[2],
		Phone:   dataArr[4],
		Address: dataArr[5],
	}
	if p.Age, err = strconv.Atoi(strings.TrimSpace(dataArr[3])); err != nil {
		return nil, err
	}
	if p.Disability, err = strconv.ParseBool(strings.TrimSpace(dataArr[6])); err != nil {
		return nil, err
	}
	p.AdditionalInfo()
	return p, nil
}

func (p *Patient) AdditionalInfo() {
	p.Id = uuid.New()
	p.RegAt = time.Now()
}
