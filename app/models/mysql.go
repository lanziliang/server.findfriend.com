package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/revel/revel"
)

var (
	Db, DbLocal *sql.DB
)

var InitDB func() = func() {
	var err error

	dbSpec := getParamString("db.spec", "")
	if Db, err = sql.Open("mysql", dbSpec); err != nil {
		revel.ERROR.Fatal(err)
	}

	dbLocalSpec := getParamString("dbLocal.spec", "")
	if DbLocal, err = sql.Open("mysql", dbLocalSpec); err != nil {
		revel.ERROR.Fatal(err)
	}
}

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}
