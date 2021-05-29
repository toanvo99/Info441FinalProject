package models

//Team represents the pokemon team a trainier will have
type Team struct {
	TeamID    int64 `json:"teamID`
	TrainerID int64 `json:"trainerID"`
}
