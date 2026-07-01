package main

import (
	"net/http"
	"time"

	"github.com/iceblade92/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token in refresh handler", err)
		return
	}

	dbtoken, err := cfg.db.GetRefreshToken(r.Context(), refreshtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token from data base", err)
		return
	}

	if dbtoken.ExpiresAt.Before(time.Now()) || dbtoken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token from data base", err)
		return
	}

	token, err := auth.MakeJWT(dbtoken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshtoken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token in revoke handler", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshtoken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token in revoke handler", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
