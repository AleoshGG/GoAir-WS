package entities

type StatusApplication string

const (
	Requested StatusApplication = "requested"
	Pending   StatusApplication = "pending"
	Complete  StatusApplication = "complete"
)

type UserRequest struct {
	Destination        string            `json:"destination"`
	Id_application     int               `json:"id_application"`
	First_name         string            `json:"first_name"`
	Last_name          string            `json:"last_name"`
	Status_application StatusApplication `json:"status_application"`
	Id_user            int               `json:"id_user"`
}