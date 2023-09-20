package database

import (
	constants "agritech/server/constants"
	"agritech/server/model"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", constants.DB_USER+":"+constants.DB_PASS+"@tcp("+constants.DB_HOST+":"+constants.DB_PORT+")/"+constants.DB_NAME)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()
	for err != nil {
		fmt.Println(err)
		time.Sleep(4 * time.Second)
		err = db.Ping()
	}

	fmt.Println("MySQL connection successful")
	return db
}

func GetTipiMisurazione(db *sql.DB) ([]model.Tipi_Misurazione, error) {
	var ret []model.Tipi_Misurazione

	rows, err := db.Query("SELECT * FROM Tipi_Misurazione")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Tipi_Misurazione
		err := rows.Scan(&a.ID, &a.Nome, &a.UnitaMisura)
		if err != nil {
			return ret, err
		}
		ret = append(ret, a)
	}
	return ret, nil
}

func GetCampi(db *sql.DB) ([]model.Campi, error) {
	var a model.Campi
	var ret []model.Campi

	rows, err := db.Query("SELECT * FROM Campi")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.Nome, &a.Lat, &a.Lon)
		if err != nil {
			return ret, err
		}
		ret = append(ret, a)
	}
	return ret, nil
}

func GetSensori(db *sql.DB) ([]model.Sensori, error) {
	var a model.Sensori
	var ret []model.Sensori

	rows, err := db.Query("SELECT * FROM Sensori")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.MAC, &a.ID_campo, &a.Lat, &a.Lon)
		if err != nil {
			return ret, err
		}
		ret = append(ret, a)
	}
	return ret, nil
}

func GetMisurazioni(db *sql.DB) ([]model.Misurazioni, error) {
	var a model.Misurazioni
	var ret []model.Misurazioni

	rows, err := db.Query("SELECT * FROM Misurazioni")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.ID_sensore, &a.ID_tipo_misurazione, &a.Valore, &a.Data_ora)
		if err != nil {
			return ret, err
		}
		ret = append(ret, a)
	}
	return ret, nil
}
