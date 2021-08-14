package user

/*
	file > ./internal/usercardjob/patient_card.csv
*/

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/config"
	"github.com/jackc/pgx/v4"
	"os"
)

func StartDialog() error {
	userChoice := getChoice()
	if userChoice == 3 {
		return errors.New("exit patient dialog chosen")
	}
	userObj, err := handleChoice(userChoice)
	if err != nil {
		return err
	}
	fmt.Printf("User for db > %#v\n", userObj)
	/*if err = writeUserDataToDb(userObj); err != nil {
		return err
	}
	fmt.Println("Данные успешно внесены и ожидают проверку [Оператором].")*/
	return nil
}

func getChoice() (userChoice int) {
	greeting := "У Вас нет карточки в нашем лечебном учереждении.\nДля того, чтобы мы могли создать карточку пациента, Вам нужно заполнить о себе данные:\n" +
		"1 -> У меня есть карточка в файле CSV.\n" +
		"2 -> У меня нет карточки, готов(а) заполнить карточку на сайте.\n" +
		"3 -> Покинуть сайт.\n" +
		"Сделайте Ваш выбор > "

	for {
		fmt.Print(greeting)
		_, err := fmt.Scan(&userChoice)
		if err != nil || userChoice > 3 || userChoice < 1 {
			fmt.Println("Ошибка ввода.")
			continue
		}
		return
	}
}

func handleChoice(userChoice int) (*User, error) {
	if userChoice == 1 {
		userObj, err := readData(readCsvToUser)
		if err != nil {
			return nil, err
		}
		return userObj, nil
	}
	userObj, err := interactiveCardDialog()
	if err != nil {
		return nil, err
	}
	return userObj, nil
}

//goland:noinspection GoUnhandledErrorResult
func readData(f func(file *os.File) (*User, error)) (*User, error) {
	file, err := os.Open(getFilePath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	userObj, err := f(file)
	if err != nil {
		return nil, err
	}
	return userObj, nil
}

func readCsvToUser(file *os.File) (*User, error) {
	csvLines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	return NewUser(csvLines[1]), nil
}

func getFilePath() (path string) {
	for {
		fmt.Print("Введите пожалуйста путь к файлу > ")
		_, err := fmt.Scan(&path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if b := checkIfFileExist(path); !b {
			fmt.Println("Файл не существует.")
			continue
		}
		break
	}
	return
}

func checkIfFileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

func interactiveCardDialog() (*User, error) {
	for {
		fmt.Println("Пожалуйста введите Ваши данные без пробелов и через запятую:")
		fmt.Println("FirstName,LastName,Email,Gender(male/female),BirthDay,Phone,Address,JobInfo," +
			"Disability(yes/no),Allergies(yes/no)")
		csvArr, err := csv.NewReader(os.Stdin).Read()
		if err != nil {
			return nil, err
		}
		if len(csvArr) != 10 {
			fmt.Println("Данные введены не корректно или не все поля были заполнены.")
			continue
		}
		return NewUser(csvArr), nil
	}
}

//goland:noinspection GoUnhandledErrorResult
func writeUserDataToDb(u *User) error {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	var env config.Env
	conn, err := pgx.Connect(context.Background(), env.DatabaseStr())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	query := `
	INSERT INTO users(first_name, last_name, email, gender, birthday, phone, address, job_info, disability, allergies, reg_day, role) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	if _, err = conn.Exec(context.Background(), query,
		u.FirstName, u.LastName, u.Email, u.GetGender(), u.BirthDayStr, u.Phone, u.Address, u.JobInfo,
		u.GetDisability(), u.GetAllergies(), u.GetRegDayStr(), u.Role); err != nil {
		return err
	}
	return nil
}
