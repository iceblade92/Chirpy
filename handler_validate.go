package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	texts := strings.Split(params.Body, " ")
	for i, word := range texts {
		a_word := strings.ToLower(word)
		if a_word == "kerfuffle" || a_word == "sharbert" || a_word == "fornax" {
			texts[i] = "****"
		}
	}
	cleanedText := strings.Join(texts, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedText,
	})
}
