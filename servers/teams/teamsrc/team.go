package teamsrc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/user"
	"strconv"
	"strings"
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
	Species   string `json:species"`
	Type1     string `json:type1"`
	Type2     string `json:type2"`
}

//from teamstore.go for SQL query function
type TeamContext struct {
	TeamStore TeamSQLStore
}

// this is for displaying all of the team the user have (with Pokemon as well?)
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
		teamList, err := tc.TeamStore.AllTeamsGetByID(u.ID)
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
		if t.TeamID == "" {
			http.Error(w, "team ID doesn't exit", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		teamInsert, err := tc.TeamStore.AddChannel(&t)
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
	if r.Method != "DELETE" {
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
	// creating a team
	/***** use a POST METHOD***/
}
