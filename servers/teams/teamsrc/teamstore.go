package teamsrc

import (
	"database/sql"
)

type TeamSQLStore struct {
	Database *sql.DB
}

func newTeamSQLStore(database *sql.DB) *TeamSQLStore {
	return &TeamSQLStore{
		Database: database,
	}
}

//show all of the team 
func (ss *TeamSQLStore) AllTeamsGetByID(id int64) ([]*Team, error) {
	var allTeams []*Team
	queryOne := "SELECT t.TeamID, p.PokemonID FROM Team t join (Pokemon p, PokemonTeam pt, User u) ON (u.TrainierID == t.TrainerID AND t.TeamID == pt.MoveSetID AND p.PokemonID == pt.PokemonID where u.ID = ?"
	rows, err := ss.Database.Query(queryOne, id)
	if err != nil {
		return nil, err
	}
	pokemonteam := PokemonTeam{}
	defer rows.Close()
	for rows.Next() {
		if err2 := rows.Scan(
			&pokemonteam.PokemonID,
			&pokemonteam.MoveSetID,
		); err2 != nil {
			return nil, err2
		}
		allTeams = append(allTeam, &pokemonteam)
	return allTeams

// show a team based on team id
func (ss *TeamSQLStore) TeamGetByID(id int64) (*Team, error) {
	rows, err := ss.Database.Query("SELECT * FROM team WHERE teamID=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	team := Team{}
	for rows.Next() {
		if err2 := rows.Scan(
			&team.TeamID,
			&team.TrainerID,
		); err2 != nil {
			return nil, err2
		}
	}	
	return team, nil
}

// NOT SURE IF THIS WORKS
func (ss *TeamSQLStore) MakeNewTeam(user *User) (*Team, error) {
	insertQuery := "INSERT into Team (TrainerID) VALUES (?)"
	response, err := ss.Database.Exec(insertQuery, user.TrainerID)

	if err != nil {
		return nil, err
	}

	team := &Team{}

	id, err2 := response.LastInsertId()

	if err2 != nil {
		return nil, err2
	}

	team.TrainerID = id
	return team, nil

}

func (ss *TeamSQLStore) DeleteTeam(id int64) error {
	insertQuery := "DELETE FROM Team WHERE TeamID=?"
	_, err := ss.Database.Exec(insertQuery,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}