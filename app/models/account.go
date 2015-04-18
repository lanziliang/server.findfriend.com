package models

import (
	"github.com/revel/revel"
)

var (
	AccountModel = new(Account)
)

type Account struct {
	Id      int    `json:"id"`
	Sex     int    `json:"sex"`
	Phone   string `json:"phone"`
	Pwd     string `json:"pwd"`
	Name    string `json:"username"`
	Email   string `json:"email"`
	Avator  string `json:"headimg"`
	Evalute int    `json:"evalute"`
}

func (m *Account) GetAccountInfoByName(name string) (Account, error) {
	sql := `SELECT id,sex,phone,pwd,name,email,avator,evalute 
			FROM account 
			WHERE name = ?`

	var account Account
	err := DbLocal.QueryRow(sql, name).Scan(&account.Id, &account.Sex, &account.Phone, &account.Pwd, &account.Name, &account.Email, &account.Avator, &account.Evalute)
	if err != nil {
		revel.ERROR.Println(err)
		return account, err
	}

	return account, nil
}

func (m *Account) GetAccountInfoById(id int) (Account, error) {
	sql := `SELECT sex,phone,name,email,avator,evalute 
			FROM account 
			WHERE id = ?`

	var account Account
	err := DbLocal.QueryRow(sql, id).Scan(&account.Sex, &account.Phone, &account.Name, &account.Email, &account.Avator, &account.Evalute)
	if err != nil {
		revel.ERROR.Println(err)
		return account, err
	}

	return account, nil
}

func (m *Account) Insert(d Account) error {
	sql := "INSERT INTO account (name, email, phone, sex, pwd) VALUES (?, ?, ?, ?, ?)"

	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	_, err = stmt.Exec(d.Name, d.Email, d.Phone, d.Sex, d.Pwd)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}

func (m *Account) UpdateAvator(id int, avator string) error {
	sql := "UPDATE account SET avator=? WHERE id=?"

	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	_, err = stmt.Exec(avator, id)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}

type PersonalActivity struct {
	Id         int    `json:"activityId"`
	Title      string `json:"title"`
	Type       string `json:"messageType"` //1:发起的  2:参与的
	CreateDate string `json:"time"`
}

func (m *Account) GetActivityByAccountId(id int) ([]PersonalActivity, error) {
	var personalActivities []PersonalActivity
	sql := `SELECT id, title, 1 as type, create_date
			FROM activity 
			WHERE account_id=?`

	rows, err := DbLocal.Query(sql, id)
	if err != nil {
		revel.ERROR.Println(err)
		return personalActivities, err
	}

	for rows.Next() {
		var personalActivity PersonalActivity
		err = rows.Scan(&personalActivity.Id, &personalActivity.Title, &personalActivity.Type, &personalActivity.CreateDate)
		if err != nil {
			revel.ERROR.Println(err)
			return personalActivities, err
		}

		personalActivities = append(personalActivities, personalActivity)
	}

	sql = `SELECT A.activity_id, B.title, 2 AS type, A.create_date
			FROM activity_join AS A LEFT JOIN activity AS B ON A.activity_id=B.id 
			WHERE A.account_id = ?`

	rows, err = DbLocal.Query(sql, id)
	if err != nil {
		revel.ERROR.Println(err)
		return personalActivities, err
	}

	for rows.Next() {
		var personalActivity PersonalActivity
		err = rows.Scan(&personalActivity.Id, &personalActivity.Title, &personalActivity.Type, &personalActivity.CreateDate)
		if err != nil {
			revel.ERROR.Println(err)
			return personalActivities, err
		}

		personalActivities = append(personalActivities, personalActivity)
	}

	return personalActivities, nil
}
