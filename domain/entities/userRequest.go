package entities

type StatusApplication string

const (
	Requested StatusApplication = "requested"
	Pending   StatusApplication = "pending"
	Complete  StatusApplication = "complete"
)

type UserRequest struct {
	Destination        string
	Id_application     int
	First_name         string
	Last_name          string
	Status_application StatusApplication
	Id_user            int
}