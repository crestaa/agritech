package database

import (
	constants "agritech/server/constants"
	"agritech/server/model"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var DB *sql.DB

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

/******************************************************************************************/
/****************				FUNCTIONS FOR API					***********************/
/******************************************************************************************/

func GetCampi() ([]model.Campi, error) {
	var a model.Campi
	var ret []model.Campi

	rows, err := DB.Query("SELECT * FROM Campi")
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

func GetCampo(id int) (model.Campi, error) {
	var a model.Campi
	err := DB.QueryRow("SELECT * FROM Campi WHERE id_campo = ? LIMIT 1", id).Scan(&a.ID, &a.Nome, &a.Lat, &a.Lon)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Campi{}, fmt.Errorf("no field found with id: %d", id)
		}
		return model.Campi{}, err
	}
	return a, nil
}

func GetReadings(id int) ([]model.Misurazioni, error) {
	var a model.Misurazioni
	var ret []model.Misurazioni

	rows, err := DB.Query("SELECT Misurazioni.* FROM Misurazioni JOIN Sensori ON Misurazioni.id_sensore = Sensori.id_sensore JOIN Campi ON Sensori.id_campo = Campi.id_campo WHERE Campi.id_campo = ?", id)
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.ID_sensore, &a.Nonce, &a.ID_tipo_misurazione, &a.Valore, &a.Data_ora)
		if err != nil {
			return ret, err
		}
		ret = append(ret, a)
	}
	return ret, nil
}

/******************************************************************************************/
/****************	FUNCTIONS FOR LISTENER ON MQTT + SAVING DATA	***********************/
/******************************************************************************************/

func GetSensori() ([]model.Sensori, error) {
	var a model.Sensori
	var ret []model.Sensori

	rows, err := DB.Query("SELECT * FROM Sensori")
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

func GetMisurazioni() ([]model.Misurazioni, error) {
	var a model.Misurazioni
	var ret []model.Misurazioni

	rows, err := DB.Query("SELECT * FROM Misurazioni")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.ID_sensore, &a.Nonce, &a.ID_tipo_misurazione, &a.Valore, &a.Data_ora)
		if err != nil {
			return ret, err
		}
		ret = append(ret, a)
	}
	return ret, nil
}

func GetSensorID(mac string) (int, error) {
	var ret int
	err := DB.QueryRow("SELECT id_sensore FROM Sensori WHERE mac = ? LIMIT 1", mac).Scan(&ret)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("no sensor found with MAC: %s", mac)
		}
		return -1, err
	}
	return ret, nil
}

func GetMeasurementTypeID(name string) (int, error) {
	var ret int
	err := DB.QueryRow("SELECT id_tipo_misurazione FROM Tipi_Misurazione WHERE nome = ? LIMIT 1", name).Scan(&ret)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("no measurement type found with name: %s", name)
		}
		return -1, err
	}
	return ret, nil
}

// returns false only if there are no doubles in Misurazioni table (based on id_sensore & nonce)
func CheckDoubles(nonce int, sensor int) (bool, error) {
	var tmp int
	err := DB.QueryRow("SELECT id_misurazione FROM Misurazioni WHERE id_sensore = ? AND nonce = ? LIMIT 1", sensor, nonce).Scan(&tmp)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, err
	}
	return true, fmt.Errorf("found double")
}

func SaveMisurazione(data model.Misurazioni) error {
	query := "INSERT INTO Misurazioni (id_sensore, tipo_misurazione, nonce, valore) VALUES (?, ?, ?, ?)"
	_, err := DB.Exec(query, data.ID_sensore, data.ID_tipo_misurazione, data.Nonce, data.Valore)
	return err
}

func GetTipiMisurazione() ([]model.Tipi_Misurazione, error) {
	var ret []model.Tipi_Misurazione

	rows, err := DB.Query("SELECT * FROM Tipi_Misurazione")
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
