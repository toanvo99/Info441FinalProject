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
func (ss *TeamSQLStore) AllTeamsGetByName(trainer string) ([]*Team, error) {
	var allTeams []*Team
	queryOne := "SELECT T.TeamID FROM Team T join User U on T.TrainerID = U.TrainerID where U.UserName = ?"
	rows, err := ss.Database.Query(queryOne, trainer)
	if err != nil {
		return nil, err
	}
	var teamid int64
	defer rows.Close()
	for rows.Next() {
		if err2 := rows.Scan(
			&teamid,
		); err2 != nil {
			return nil, err2
		}
		teamadd, err := ss.TeamGetByID(teamid)
		if err != nil {
			return nil, err
		}
		allTeams = append(allTeams, teamadd)
	}
	return allTeams, nil
}

// show a team based on team id
func (ss *TeamSQLStore) TeamGetByID(id int64) (*Team, error) {
	rows, err := ss.Database.Query("SELECT P.PokemonID, S.SpeciesName, T1.TypeName, T2.TypeName FROM Team join PokemonTeam PT on T.TeamID = PT.TeamID join Pokemon P ON P.PokemonID = PT.PokemonID join Species S ON S.SpeciesID = P.SpeciesID join [Type] T1 ON T1.TypeID = S.Type1ID join [Type] T2 ON T2.TypeID = S.Type2ID WHERE TeamID=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	team := Team{}
	pokemonlist := Pokemon{}
	for rows.Next() {
		if err2 := rows.Scan(
			&pokemonlist.PokemonID,
			&pokemonlist.Species,
			&pokemonlist.Type1,
			&pokemonlist.Type2,
			&pokemonlist.Sprite,
		); err2 != nil {
			return nil, err2
		}
		team.Pokemons = append(team.Pokemons, &pokemonlist)
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

	insertQuery := "SELECT S.SpeciesID FROM Species S WHERE S.SpeciesName = ?"
	resp1, err := ss.Database.Query(insertQuery, pokemon)
	if err != nil {
		return fmt.Errorf("failed to find Pokemon %v", err)
	}
	defer resp1.Close()
	var pid int64
	for resp1.Next() {
		if err2 := resp1.Scan(
			&pid,
		); err2 != nil {
			return err2
		}
	}
	insertQuery2 := "INSERT INTO Pokemon (SpeciesID) VALUES (?)"
	resp2, err := ss.Database.Exec(insertQuery2, pid)
	if err != nil {
		return fmt.Errorf("failed to insert %v", err)
	}
	id, err2 := resp2.LastInsertId()
	if err2 != nil {
		return fmt.Errorf("failed to insert %v", err2)
	}
	insertQuery3 := "INSERT INTO PokemonTeam (PokemonID, PokemonTeamID) VALUES (?, ?)"
	_, err = ss.Database.Exec(insertQuery3, id, teamID)

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
