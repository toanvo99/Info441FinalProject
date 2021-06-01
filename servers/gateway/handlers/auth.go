package handlers

import (
	"Info441FinalProject/servers/gateway/models"
	"Info441FinalProject/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

// TODO
// Similar to our past assignments, this .go file will handle any necessary user/session authentication
// for preforming requests to our API.

// Handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the trainer store
type HandlerContext struct {
	SignKey      string
	SessionStore sessions.Store
	TrainerStore models.Store
}

type SessionState struct {
	BeginTime time.Time
	User      *models.User
}

// This function handles requests for the trainer resources. Only accept "POST" method for creating
// new trainer accounts
func (handlerContext *HandlerContext) TrainersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("The request body must be in JSON"))
			return
		}
		dec := json.NewDecoder(r.Body)
		newTrainer := &models.NewUser{}
		if err := dec.Decode(newTrainer); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		validatedTrainer, err := newTrainer.ToUser()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		handlerContext.TrainerStore.Insert(validatedTrainer)

		/* TODO: Do we need to log the signed in trainers? */

		// if len(r.Header.Get("X-Forwarded-For")) > 0 {
		// 	ips := strings.Split(r.Header.Get("X-Forwarded-For"), ", ")
		// 	handlerContext.UserStore.InsertLog(validatedUser.ID, ips[0])
		// } else {
		// 	handlerContext.UserStore.InsertLog(validatedUser.ID, r.RemoteAddr)
		// }

		signKey := handlerContext.SignKey
		sessionStore := handlerContext.SessionStore
		sessionState := &SessionState{User: validatedTrainer, BeginTime: time.Now()}
		sessions.BeginSession(signKey, sessionStore, sessionState, w)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err := enc.Encode(validatedTrainer); err != nil {
			w.Write([]byte("Failed to encode user to JSON"))
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// This function handles requests for a specific trainer. The resource path will be "/v1/trainer/{trainderID}"
// "/v1/trainer/me" refers to the currently-authenticated trainer.
func (handlerContext *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// The current user must be authenticated first before proceeding
	signKey := handlerContext.SignKey
	sessionStore := handlerContext.SessionStore
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, signKey, sessionStore, sessionState)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	currentTrainer := sessionState.User
	if r.Method == http.MethodGet {
		idString := path.Base(r.URL.Path)
		userID, _ := strconv.ParseInt(idString, 10, 64)
		user, err := handlerContext.TrainerStore.GetByID(userID)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			if err := enc.Encode(user); err != nil {
				w.Write([]byte("Failed to encode user to JSON"))
				return
			}
			return
		}
	} else if r.Method == http.MethodPatch {
		idString := path.Base(r.URL.Path)
		if idString != "me" && strconv.FormatInt(currentTrainer.TrainerID, 10) != idString {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		dec := json.NewDecoder(r.Body)
		updates := &models.Updates{}
		if err := dec.Decode(updates); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		// TODO: Get current user and update it
		err := currentTrainer.ApplyUpdates(updates)
		if err != nil {
			w.Write([]byte("Failed to update user"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(currentTrainer); err != nil {
			w.Write([]byte("Failed to encode user to JSON"))
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (handlerContext *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		dec := json.NewDecoder(r.Body)
		credentials := &models.Credentials{}
		if err := dec.Decode(credentials); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		userStore := handlerContext.TrainerStore
		currentUser, err := userStore.GetByEmail(credentials.Email)
		if err != nil {
			fakePassword := "fakefakefakefake"
			_ = currentUser.Authenticate(fakePassword)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			password := credentials.Password
			err := currentUser.Authenticate(password)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		signKey := handlerContext.SignKey
		sessionStore := handlerContext.SessionStore
		sessionState := &SessionState{User: currentUser, BeginTime: time.Now()}
		sessions.BeginSession(signKey, sessionStore, sessionState, w)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err := enc.Encode(currentUser); err != nil {
			w.Write([]byte("Failed to encode user to JSON"))
			return
		}
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (handlerContext *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		if path.Base(r.URL.Path) != "mine" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Unauthorized status"))
			return
		}
		signKey := handlerContext.SignKey
		sessionStore := handlerContext.SessionStore
		_, err := sessions.EndSession(r, signKey, sessionStore)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Signed out"))
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
