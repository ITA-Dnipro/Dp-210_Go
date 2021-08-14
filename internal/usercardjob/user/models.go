package user

import "time"

type User struct {
	FirstName   string
	LastName    string
	Email       string
	gender      string
	BirthDayStr string
	Phone       string
	Address     string
	JobInfo     string
	disability  bool
	allergies   bool
	RegDayTime  time.Time
	regDayStr   string
	Role        int
}

func NewUser(csvArr []string) *User {
	userObj := &User{
		FirstName:   csvArr[0],
		LastName:    csvArr[1],
		Email:       csvArr[2],
		gender:      csvArr[3],
		BirthDayStr: csvArr[4],
		Phone:       csvArr[5],
		Address:     csvArr[6],
		JobInfo:     csvArr[7],
		RegDayTime:  time.Now(),
		Role:        0,
	}
	userObj.SetDisability(csvArr[8])
	userObj.SetAllergies(csvArr[9])
	return userObj
}

func (u *User) GetRegDayStr() string {
	if u.regDayStr == "" {
		u.regDayStr = u.RegDayTime.Format("2006-01-02")
	}
	return u.regDayStr
}

func (u *User) GetGender() string {
	if u.gender != "male" && u.gender != "female" && u.gender != "not specified" {
		u.gender = "not specified"
	}
	return u.gender
}

func (u *User) SetDisability(info string) {
	if info == "yes" {
		u.disability = true
	}
}

func (u *User) SetAllergies(info string) {
	if info == "yes" {
		u.allergies = true
	}
}

func (u *User) GetDisability() bool {
	return u.disability
}

func (u *User) GetAllergies() bool {
	return u.allergies
}
