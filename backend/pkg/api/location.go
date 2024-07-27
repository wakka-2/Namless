package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/wakka-2/Namless/backend/pkg/models"
)

// RequestAllLocations replies with all locations.
func (r *RESTAPI) RequestAllLocations(writer http.ResponseWriter, req *http.Request) {
	result, err := r.locationService.GetAll(req.Context())
	if err != nil {
		r.handleError(writer, "could not retrieve entry", http.StatusInternalServerError)
		return
	}

	asJSON, err := json.Marshal(result)
	if err != nil {
		r.handleError(writer, "could not marshal", http.StatusInternalServerError)
		return
	}

	err = write(writer, asJSON, http.StatusOK)
	if err != nil {
		log.Default().Printf("could not write: %s", err)
	}
}

// RequestLocation will return the Location with a given ID.
//
//nolint:dupl
func (r *RESTAPI) RequestLocation(writer http.ResponseWriter, req *http.Request) {
	key := req.PathValue("id")
	if key == "" {
		r.handleError(writer, "missing key", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(key)
	if err != nil {
		r.handleError(writer, "could not convert to int", http.StatusBadRequest)
		return
	}

	result, err := r.locationService.Get(req.Context(), id)
	if err != nil {
		r.handleError(writer, "could not retrieve location", http.StatusInternalServerError)
		return
	}

	asJSON, err := json.Marshal(result)
	if err != nil {
		r.handleError(writer, "could not marshal", http.StatusInternalServerError)
		return
	}

	err = write(writer, asJSON, http.StatusOK)
	if err != nil {
		log.Default().Printf("could not write: %s", err)
	}
}

// CreateLocation creates a new locattion.
func (r *RESTAPI) CreateLocation(writer http.ResponseWriter, req *http.Request) {
	input := models.Location{}

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.locationService.Add(req.Context(), input)
	if err != nil {
		r.handleError(writer, "could not create location", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}
