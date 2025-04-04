package models

type ProcessParameter struct {
	ID          int    `db:"id" json:"id"`
	EquipmentID int    `db:"id_equipment" json:"equipmentId"`
	Name        string `db:"name" json:"name"`
	Units       string `db:"units" json:"units"`
}

type ReferenceParameter struct {
	ID      int     `db:"id" json:"id"`
	ParamID int     `db:"id_param" json:"paramId"`
	Min     float64 `db:"min_value" json:"min"`
	Max     float64 `db:"max_value" json:"max"`
}

type CurrentParameter struct {
	Timestamp string  `db:"timestamp" json:"timestamp"`
	ParamID   int     `db:"id_param" json:"paramId"`
	Value     float64 `db:"value" json:"value"`
}
