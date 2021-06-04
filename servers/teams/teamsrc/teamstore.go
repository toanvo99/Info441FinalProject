package teamsrc

import (
	"database/sql"
	"fmt"
)

type TeamSQLStore struct {
	Database *sql.DB
}

func NewTeamSQLStore(database *sql.DB) *TeamSQLStore {
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
	}
	return allTeams
}

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
func (ss *TeamSQLStore) MakeNewTeam(uid int64) (*Team, error) {
	insertQuery := "INSERT into Team (TrainerID) VALUES (?)"
	response, err := ss.Database.Exec(insertQuery, uid)

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
	_, err := ss.Database.Exec(insertQuery, id)
	if err != nil {
		return err
	}
	return nil
}

// Add a pokemon to the given team
func (ss *TeamSQLStore) AddPokemonToTeam(teamID int64, pid int64) error {
	insertQuery := "INSERT INTO PokemonTeam (PokemonID, PokemonTeamID) VALUES (?, ?)"
	resp, err := ss.Database.Exec(insertQuery, pid, teamID)
	if err != nil {
		return fmt.Errorf("failed to insert %v", err)
	}
	return nil
}

// Delete a pokemon from the given team
func (ss *TeamSQLStore) DeletePokemonFromTeam(teamID int64, pid int64) (*Team, error) {
	delq := "delete from PokemonTeam where PokemonTeamID=? and PokemonID=?"
	_, err := cs.db.Exec(delq, teamID, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete pokemon from team %v", err)
	}
	newTeam, err := cs.GetTeamByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated channel %v", err)
	}
	return newTeam, nil
}
