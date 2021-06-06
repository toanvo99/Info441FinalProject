package teamsrc

import (
	"database/sql"
	"fmt"
	"github.com/mtslzr/pokeapi-go"
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
func (ss *TeamSQLStore) AllTeamsGetByName(trainer string) ([]*Team, error) {
	var allTeams []*Team
	queryOne := "SELECT * FROM Team T join User U on T.TrainerID = U.TrainerID join PokemonTeam PT on T.TeamID = PT.TeamID join Pokemon P on P.PokemonID = PT.PokemonID where Team.UserName = ?"
	rows, err := ss.Database.Query(queryOne, trainer)
	if err != nil {
		return nil, err
	}
	pokemon := Pokemon{}
	defer rows.Close()
	for rows.Next() {
		if err2 := rows.Scan(
			&pokemon.PokemonID,
			&pokemon.Species,
			&pokemon.Type1,
			&pokemon.Type2,
		); err2 != nil {
			return nil, err2
		}
		allTeams = append(allTeams, &pokemon)
	}
	return allTeams, nil
}

// show a team based on team id
func (ss *TeamSQLStore) TeamGetByID(id int64) (*Team, error) {
	rows, err := ss.Database.Query("SELECT * FROM Team join PokemonTeam PT on T.TeamID = P.TeamID join Pokemon P WHERE TeamID=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	team := Team{}
	for rows.Next() {
		if err2 := rows.Scan(
			&team.TeamID,
			&team.Trainer,
			,
		); err2 != nil {
			return nil, err2
		}
	}
	return &team, nil
}

// add a team
func (ss *TeamSQLStore) MakeNewTeam(t *Team) (*Team, error) {
	insertQuery := "INSERT into Team (Trainer, Pokemon) VALUES (?,?)"
	response, err := ss.Database.Exec(insertQuery,
		t.Trainer,
		t.Pokemons,
	)
	if err != nil {
		return nil, err
	}

	id, err2 := response.LastInsertId()
	if err2 != nil {
		return nil, err
	}

	t.TeamID = id

	return t, nil
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
func (ss *TeamSQLStore) AddPokemonToTeam(teamID int64, pokemon string) error {
	insertQuery := "INSERT INTO PokemonTeam (PokemonID, PokemonTeamID) VALUES (?, ?)"
	_, err := ss.Database.Exec(insertQuery, pid, teamID)
	if err != nil {
		return fmt.Errorf("failed to insert %v", err)
	}
	return nil
}

// Delete a pokemon from the given team
func (ss *TeamSQLStore) DeletePokemonFromTeam(teamID int64, pid int64) (*Team, error) {
	delq := "delete from PokemonTeam where TeamID=? and PokemonID=?"
	_, err := ss.Database.Exec(delq, teamID, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete pokemon from team %v", err)
	}
	newTeam, err := ss.TeamGetByID(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated channel %v", err)
	}
	return newTeam, nil
}
