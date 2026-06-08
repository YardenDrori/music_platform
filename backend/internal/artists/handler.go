package artists

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	httputils "github.com/YardenDrori/music-platform/internal/http_utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) NewArtist(w http.ResponseWriter, r *http.Request) error {
	req := &NewArtistReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperrors.NewErrBadRequest("invalid request body")
	}

	err := h.service.NewArtist(r.Context(), *req)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusCreated, nil)
}

func (h *handler) GetArtistsByName(w http.ResponseWriter, r *http.Request) error {
	queries := r.URL.Query()
	name := queries.Get("name")
	if name == "" {
		return apperrors.NewErrBadRequest("query missing name field")
	}
	artists, err := h.service.GetArtistsByName(r.Context(), name)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, artists)
}

func (h *handler) GetArtistByID(w http.ResponseWriter, r *http.Request) error {
	queries := r.URL.Query()
	idRaw := queries.Get("id")
	if idRaw == "" {
		return apperrors.NewErrBadRequest("query missing id field")
	}
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("id provided is not a valid UUID")
	}

	artist, err := h.service.GetArtistByID(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, artist)
}

func (h *handler) UpdateArtistDetails(w http.ResponseWriter, r *http.Request) error {
	req := &UpdateArtistReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperrors.NewErrBadRequest("invalid request body")
	}

	err := h.service.UpdateArtistDetails(r.Context(), req)
	if err != nil {
		return err
	}

	return httputils.WriteResponse(w, http.StatusOK, nil)
}

func (h *handler) SoftDeleteArtist(w http.ResponseWriter, r *http.Request) error {
	queries := r.URL.Query()
	idRaw := queries.Get("id")
	if idRaw == "" {
		return apperrors.NewErrBadRequest("query missing id field")
	}
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("id provided is not a valid UUID")
	}

	err = h.service.SoftDeleteArtist(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}

func (h *handler) HardDeleteArtist(w http.ResponseWriter, r *http.Request) error {
	queries := r.URL.Query()
	idRaw := queries.Get("id")
	if idRaw == "" {
		return apperrors.NewErrBadRequest("query missing id field")
	}
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("id provided is not a valid UUID")
	}

	err = h.service.HardDeleteArtist(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}
