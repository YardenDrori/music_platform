package artists

import (
	"encoding/json"
	"fmt"
	"io"
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

func (h *handler) GetArtistsByNameOrAlias(w http.ResponseWriter, r *http.Request) error {
	queries := r.URL.Query()
	name := queries.Get("name")
	if name == "" {
		return apperrors.NewErrBadRequest("query missing name field")
	}
	artists, err := h.service.GetArtistsByNameOrAlias(r.Context(), name)
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
	artistIDRaw := r.PathValue("id")
	if artistIDRaw == "" {
		return apperrors.NewErrInternal().
			WithInternal("user accessed UpdateArtistDetails handler without \"id\" in the url path")
	}
	artistID, err := uuid.Parse(artistIDRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("invalid artist id").WithCause(err)
	}

	req := &UpdateArtistReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperrors.NewErrBadRequest("invalid request body")
	}
	req.ID = artistID

	err = h.service.UpdateArtistDetails(r.Context(), req)
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

const maxPfpBytes = 5 * 1024 * 1024 // 5MB
func (h *handler) AddProfilePicture(w http.ResponseWriter, r *http.Request) error {
	artistIDRaw := r.PathValue("id")
	if artistIDRaw == "" {
		return apperrors.NewErrInternal().
			WithInternal("user accessed AddProfilePicture handler without \"id\" in the url path")
	}
	artistID, err := uuid.Parse(artistIDRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("invalid artist id").WithCause(err)
	}

	// we read one byte more to be able to distuinguish from EOF to File cut off
	pictureBytes, err := io.ReadAll(io.LimitReader(r.Body, maxPfpBytes+1))
	if err != nil {
		return apperrors.NewErrInternal().WithCause(err)
	}
	if len(pictureBytes) > maxPfpBytes {
		return apperrors.NewErrBadRequest(fmt.Sprintf("file size over %d bytes", maxPfpBytes))
	}

	err = h.service.UploadArtistProfilePicture(r.Context(), pictureBytes, artistID)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}

const maxBannerBytes = 200 * 1024 * 1024 // 200MB
func (h *handler) AddBanner(w http.ResponseWriter, r *http.Request) error {
	artistIDRaw := r.PathValue("id")
	if artistIDRaw == "" {
		return apperrors.NewErrInternal().
			WithInternal("user accessed AddBanner handler without \"id\" in the url path")
	}
	artistID, err := uuid.Parse(artistIDRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("invalid artist id").WithCause(err)
	}

	// we read one byte more to be able to distuinguish from EOF to File cut off
	pictureBytes, err := io.ReadAll(io.LimitReader(r.Body, maxBannerBytes+1))
	if err != nil {
		return apperrors.NewErrInternal().WithCause(err)
	}
	if len(pictureBytes) > maxBannerBytes {
		return apperrors.NewErrBadRequest(fmt.Sprintf("file size over %d bytes", maxBannerBytes))
	}

	err = h.service.UploadArtistBannerPicture(r.Context(), pictureBytes, artistID)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}
