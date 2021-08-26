package entity

import (
	"github.com/google/uuid"
	"time"
)

type PatientCard struct {
	uuid        uuid.UUID
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	BirthDayStr string `json:"birthday_str"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	JobInfo     string `json:"job_info"`
	Disability  string `json:"disability"`
	DisabilityB bool
	Allergies   string `json:"allergies"`
	AllergiesB  bool
	RegDayTime  time.Time
	regDayStr   string
	Role        int
}

func New(dataArr []string) *PatientCard {
	patientCard := &PatientCard{
		FirstName:   dataArr[0],
		LastName:    dataArr[1],
		Email:       dataArr[2],
		Gender:      dataArr[3],
		BirthDayStr: dataArr[4],
		Phone:       dataArr[5],
		Address:     dataArr[6],
		JobInfo:     dataArr[7],
		Disability:  dataArr[8],
		Allergies:   dataArr[9],
	}
	patientCard.CorrectData()
	return patientCard
}

func (pc *PatientCard) CorrectData() {
	pc.uuid = uuid.New()
	if pc.Gender != "male" && pc.Gender != "female" && pc.Gender != "not specified" {
		pc.Gender = "not specified"
	}
	if pc.Disability != "yes" && pc.Disability != "no" {
		pc.Disability = "no"
	}
	if pc.Disability == "yes" {
		pc.DisabilityB = true
	}
	if pc.Allergies == "yes" {
		pc.AllergiesB = true
	}
	if pc.Allergies != "yes" && pc.Allergies != "no" {
		pc.Allergies = "no"
	}
	pc.RegDayTime = time.Now()
	pc.regDayStr = pc.RegDayTime.Format("2006-01-02")
}
