package sessions

package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	id, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	store.Save(id, sessionState)
	w.Header().Add(headerAuthorization, schemeBearer+id.String())
	return id, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	authHeader := r.Header.Get(headerAuthorization)
	if len(authHeader) == 0 {
		authHeader = r.URL.Query().Get(paramAuthorization)
	}
	splitHeader := strings.Fields(authHeader)
	var sessionId string
	if len(splitHeader) == 2 {
		sessionId = splitHeader[1]
		if splitHeader[0] != "Bearer" {
			return InvalidSessionID, fmt.Errorf("invalid scheme prefix")
		}
	} else {
		return InvalidSessionID, fmt.Errorf("missing scheme prefix")
	}
	id, err := ValidateID(sessionId, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	return id, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	id, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	err2 := store.Get(id, sessionState)
	if err2 != nil {
		return InvalidSessionID, err2
	}
	return id, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	id, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	store.Delete(id)
	return id, nil
}
