/*
Package api offers handlers for a REST API server.
*/
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wakka-2/Namless/backend/pkg/service"
	"github.com/wakka-2/Namless/backend/pkg/types"
)

// RESTAPI offers handlers.
type RESTAPI struct {
	dataService     *service.Data
	locationService *service.Location
}

// New builds a new REST API.
func New(
	dataService *service.Data,
	locationService *service.Location,
) *RESTAPI {
	return &RESTAPI{
		dataService:     dataService,
		locationService: locationService,
	}
}

// BuildMultiplexer adds handlers and middleware to the resulting multiplexer.
func (r *RESTAPI) BuildMultiplexer() http.Handler {
	multiplexer := http.NewServeMux()

	multiplexer.Handle("GET /data/{key}", http.HandlerFunc(r.Request))
	multiplexer.Handle("POST /data", http.HandlerFunc(r.Create))
	multiplexer.Handle("PUT /data", http.HandlerFunc(r.Update))
	multiplexer.Handle("DELETE /data/{key}", http.HandlerFunc(r.Delete))
	multiplexer.Handle("GET /location/{id}", http.HandlerFunc(r.RequestLocation))
	multiplexer.Handle("GET /location", http.HandlerFunc(r.RequestAllLocations))
	multiplexer.Handle("POST /location", http.HandlerFunc(r.CreateLocation))
	multiplexer.Handle("POST /token", http.HandlerFunc(r.CreateToken))
	multiplexer.Handle("GET /two/{name}", http.HandlerFunc(r.CreateToken2))

	return RecoverMiddleware(EnableCORS(multiplexer))
}

// Create will create a new data entry.
// @Summary      Create will create a new data entry.
// @Accept       json
// @Produce      json
// @Param        models.Data	payload		string				true	"Request Body"
// @Success      200		{object}	string
// @Failure      400		{object}	ErrorMessage
// @Router       /create	[post].
//
//nolint:dupl
func (r *RESTAPI) Create(writer http.ResponseWriter, req *http.Request) {
	input := types.Pair{}

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.dataService.Add(req.Context(), input.Key, input.Value)
	if err != nil {
		r.handleError(writer, "could not create entry", http.StatusInternalServerError)
		return
	}
}

// Request will retrieve the data entry with a given key.
// @Summary      Request will retrieve the data entry with a given key.
// @Accept       json
// @Produce      json
// @Param        key		path		string				true	"Request Path"
// @Success      200		{object}	string
// @Failure      400		{object}	ErrorMessage
// @Router       /create	[post].
func (r *RESTAPI) Request(writer http.ResponseWriter, req *http.Request) {
	key := req.PathValue("key")
	if key == "" {
		r.handleError(writer, "missing key", http.StatusBadRequest)
		return
	}

	result, err := r.dataService.Get(req.Context(), key)
	if err != nil {
		r.handleError(writer, "could not retrieve entry", http.StatusInternalServerError)
		return
	}

	err = write(writer, []byte(result), http.StatusOK)
	if err != nil {
		log.Default().Printf("could not write: %s", err)
	}
}

// Update will update an existing data entry.
// @Summary      Update will update an existing data entry.
// @Accept       json
// @Produce      json
// @Param        models.Data	payload		string				true	"Request Body"
// @Success      200		{object}	string
// @Failure      400		{object}	ErrorMessage
// @Router       /create	[post].
//
//nolint:dupl
func (r *RESTAPI) Update(writer http.ResponseWriter, req *http.Request) {
	input := types.Pair{}

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.dataService.Update(req.Context(), input.Key, input.Value)
	if err != nil {
		r.handleError(writer, "could not update entry", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

// Delete will delete an existing data entry.
// @Summary      Delete will delete an existing data entry.
// @Accept       json
// @Produce      json
// @Param        key		path		string				true	"Request Path"
// @Success      200		{object}	string
// @Failure      400		{object}	ErrorMessage
// @Router       /create	[post].
func (r *RESTAPI) Delete(writer http.ResponseWriter, req *http.Request) {
	key := req.PathValue("key")
	if key == "" {
		r.handleError(writer, "missing key", http.StatusBadRequest)
		return
	}

	err := r.dataService.Delete(req.Context(), key)
	if err != nil {
		r.handleError(writer, "could not delete entry", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// handleError wraps an error message in a struct and sends it.
func (r *RESTAPI) handleError(w http.ResponseWriter, message string, statusCode uint) {
	err := writeJSON(w, ErrorMessage{Error: message}, statusCode)
	if err != nil {
		log.Fatalf("Could not write error response: %s", err)
	}
}
