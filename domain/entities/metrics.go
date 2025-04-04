package entities

type Sensor struct {
	Id_sensor   []string `json:"id_sensor"`
	Air_quality int      `json:"air_quality"`
	Temperature float64  `json:"temperature"`
	Humidity    float64  `json:"humidity"`
	Id_device   string   `json:"id_device"`
	Ventilador  string   `json:"ventilador"`
	Id_place    int      `json:"id_place"`
}