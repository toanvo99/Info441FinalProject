package teamsrc

type TeamStore interface {

	//show all of the team
	AllTeamsGetByID(id int64) ([]*Pokemon, error)

	//make a new Pokemon team
	MakeNewTeam(id int64) (*Team, error)

	//delete a pokemon team
	DeleteTeam(id int64) error

	//delete pokemon from team
	DeletePokemonFromTeam(teamID int64, pid int64) (*Team, error)

	//add pokemon to team
	AddPokemonToTeam(teamID int64, pokemon string) error
}
