package model

type Message struct {
	MAC   string  `json:"mac"`
	ID    int     `json:"id"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

type Tipi_Misurazione struct {
	ID          int
	Nome        string
	UnitaMisura string
}

type Campi struct {
	ID   int
	Lat  float32
	Lon  float32
	Nome string
}

type Sensori struct {
	ID       int
	MAC      string
	ID_campo int
	Lat      float32
	Lon      float32
}

type Misurazioni struct {
	ID                  int
	ID_sensore          int
	ID_tipo_misurazione int
	Valore              float32
	Data_ora            string
}
