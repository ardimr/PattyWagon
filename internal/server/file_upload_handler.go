package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/observability"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.Tracer.Start(r.Context(), "handler.file_upload")
	defer span.End()

	if r.Method != "POST" {
		sendErrorResponse(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	if err := r.ParseMultipartForm(150 << 10); err != nil { // 150 KB
		log.Printf("error parsing multipart form: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error parsing multipart form: %v", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("invalid request: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer file.Close()

	uploadedFile, err := s.service.UploadFile(ctx, file, header.Filename, header.Size)
	if err != nil {
		switch err {
		case constants.ErrMaximumFileSize:
			log.Println("invalid file size")
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		case constants.ErrInvalidFileType:
			log.Println("invalid file type")
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			log.Printf("internal server error: %v", err)
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := FileUploadResponse{
		FileID:           strconv.FormatInt(uploadedFile.ID, 10),
		FileUri:          uploadedFile.Uri,
		FileThumbnailUri: uploadedFile.ThumbnailUri,
	}

	sendResponse(w, http.StatusOK, resp)
	defer r.Body.Close()

}
