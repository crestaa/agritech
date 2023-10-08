package model

type Message struct {
	MAC   string  `json:"m"`
	ID    int     `json:"i"`
	Value float32 `json:"v"`
	Type  string  `json:"t"`
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
	Nonce               int
	Valore              float32
	Data_ora            string
}
