package models

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type Phone struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Provider    string `json:"provider"`
}

type PhoneModel struct {
	DB *sql.DB
}

func (m PhoneModel) GetAll(oddEven string) ([]*Phone, error) {
	oddEvenValue := 0
	if oddEven == "odd" {
		oddEvenValue = 1
	}

	query := `
		SELECT id, phone_number, provider 
		FROM phones
		WHERE (
		    CAST(SUBSTRING(phone_number FROM LENGTH(phone_number) FOR 1) AS INTEGER) % 2 = $2 OR $1 = ''
		)
	`
	args := []interface{}{oddEven, oddEvenValue}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var phones []*Phone
	for rows.Next() {
		var phone Phone
		err := rows.Scan(&phone.ID, &phone.PhoneNumber, &phone.Provider)
		if err != nil {
			return nil, err
		}
		phones = append(phones, &phone)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return phones, err
}

func (m PhoneModel) Insert(phone *Phone) error {
	stmt := `
		INSERT INTO phones (phone_number, provider)
		VALUES($1, $2)
		RETURNING id
	`
	args := []interface{}{phone.PhoneNumber, phone.Provider}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&phone.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m PhoneModel) GetRandPhoneNumber() (string, string, error) {
	for {
		provider := GetRandomProvider()
		phoneNumber := GetRandomPhoneNumber(provider)

		var exist bool
		stmt := `SELECT EXISTS(SELECT true FROM phones WHERE provider = $1 AND phone_number = $2)`
		args := []interface{}{provider, phoneNumber}
		err := m.DB.QueryRow(stmt, args...).Scan(&exist)
		if err != nil {
			return "", "", err
		}
		if !exist {
			return provider, phoneNumber, nil
		}
	}
}

func GetRandomProvider() string {
	providers := []string{"telkomsel", "xl", "indosat", "tri", "smartfreen"}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(providers))

	return providers[randomIndex]
}

func GetRandomPhoneNumber(provider string) string {
	var prefix string

	switch provider {
	case "telkomsel":
		prefix = "0811"
	case "xl":
		prefix = "0817"
	case "indosat":
		prefix = "0814"
	case "tri":
		prefix = "0896"
	case "smartfren":
		prefix = "0881"
	default:
		prefix = "0812" // Default to a common prefix
	}

	rand.Seed(time.Now().UnixNano())
	randomDigits := fmt.Sprintf("%07d", rand.Intn(10000000))

	phoneNumber := prefix + randomDigits

	return phoneNumber
}
