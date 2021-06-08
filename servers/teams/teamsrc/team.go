package teamsrc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/user"
	"strconv"
	"strings"

	"github.com/info441-sp21/assignments-tomgerber/servers/gateway/models/users"
)

//Team represents the pokemon team a trainier will have
type Team struct {
	TeamID   int64      `json:"teamID"`
	Trainer  user.User  `json:"trainer"`
	Pokemons []*Pokemon `json:"pokemons"`
}

//Pokemon represents the specific Pokemon
type Pokemon struct {
	PokemonID int64  `json:"pokemonID"`
	Species   string `json:"species"`
	Type1     string `json:"type1"`
	Type2     string `json:"type2"`
	Sprite    string `json:"sprite"`
}

// user struct as it will be needed by our handlers
type User struct {
	TrainerID int64  `json:"trainerID"`
	Email     string `json:"-"`
	PassHash  []byte `json:"-"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//from teamstore.go for SQL query function
type TeamContext struct {
	TeamStore *TeamSQLStore
}

// this is for displaying all of the team the user have (with Pokemon as well?)
// "/v1/teams"
func (tc *TeamContext) AllTeamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User") != "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
	}
	u := users.User{}
	err := json.NewDecoder(strings.NewReader(r.Header.Get("X-User"))).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		teamList, err := tc.TeamStore.AllTeamsGetByName(u.UserName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		//writing response containing all channels user is a part of
		resp, err := json.Marshal(teamList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)

	}
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "the request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		t := Team{}
		err := json.NewDecoder(r.Body).Decode(&t)
		if t.Trainer.Username == "" {
			http.Error(w, "Trainers doesn't exit", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		teamInsert, err := tc.TeamStore.MakeNewTeam(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(teamInsert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	} else {
		http.Error(w, "AllTeam: only GET and POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

}

// managing the team
// "/v1/{userID}"
func (tc *TeamContext) TeamManageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User") != "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
	}
	u := users.User{}
	err := json.NewDecoder(strings.NewReader(r.Header.Get("X-User"))).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// deleting a team
	if r.Method == "DELETE" {
		urlID := strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]
		teamid, err := strconv.ParseInt(urlID, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = tc.TeamStore.TeamGetByID(teamid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		errDelete := tc.TeamStore.DeleteTeam(teamid, u.ID)
		if errDelete != nil {
			if errDelete.Error() == "Team not found" {
				http.Error(w, errDelete.Error(), http.StatusForbidden)
				return
			}
			http.Error(w, errDelete.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Team successfully deleted")
	}
	// creating a team using JSON from request body
	/***** use a POST METHOD***/
	if r.Method == http.MethodPost {
		dec := json.NewDecoder(r.Body)
		newTeam := &Team{}
		if err := dec.Decode(newTeam); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newTeam, err := tc.TeamStore.MakeNewTeam(u.TrainerID)
		if err != nil {
			fmt.Printf("error inserting channel: %v\n", err)
		}
		// response with JSON formatted newly created team
		enc := json.NewEncoder(w)
		if err := enc.Encode(newTeam); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application-json")
		return
	}
}

// TeamBuilderHandler handles pokemon in the current team
// "/v1/teams/{teamID}"
func (tc *TeamContext) TeamBuilderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User") != "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
	}
	u := users.User{}
	err := json.NewDecoder(strings.NewReader(r.Header.Get("X-User"))).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	urlID := strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]
	teamid, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	curTeam, err := tc.TeamStore.TeamGetByID(teamid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// POST: Add pokemon to the current team
	if r.Method == http.MethodPost {
		dec := json.NewDecoder(r.Body)
		pokemonToAdd := &Pokemon{}
		if err := dec.Decode(pokemonToAdd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		curTeam.Pokemons = append(curTeam.Pokemons, pokemonToAdd)
		err := tc.TeamStore.AddPokemonToTeam(curTeam.TeamID, pokemonToAdd.Species)
		if err != nil {
			fmt.Printf("failed to add pokemon to current team: %v", err.Error())
			return
		}
		w.Write([]byte("pokemon was successfully added to the team"))
		w.WriteHeader(http.StatusCreated)
		return
	}
	// DELETE: Remove pokemon given by the request body from the current team
	if r.Method == http.MethodDelete {
		dec := json.NewDecoder(r.Body)
		pokemonToDel := &Pokemon{}
		if err := dec.Decode(pokemonToDel); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err := tc.TeamStore.DeletePokemonFromTeam(curTeam.TeamID, pokemonToDel.PokemonID)
		if err != nil {
			fmt.Printf("failed to delete pokemon from current team: %v", err.Error())
			return
		}
		w.Write([]byte("the user was successfully removed from the team"))
		w.WriteHeader(200)
		return
	}
}
