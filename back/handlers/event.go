package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real-time-forum/services"
	"real-time-forum/utils"
)

type ReceiveEvent struct {
	Type        string `json:"type"`         //comment or post
	ContentType string `json:"content_type"` // like dislikec
	Id          int    `json:"id"`           // ID post
}

func Event(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Missing cookie")
		return
	}
	checkSession, userID := utils.CheckSession(db, cookie.Value)
	if !checkSession {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Session invalid")
		return
	}

	var event ReceiveEvent
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		utils.SendErrorResponse(w, 400, "Erreur de parsing du JSON")
		return
	}

	if event.Id <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	switch event.Type {
	case "comment":
		handleCommentEvent(db, w, userID, event)
	case "post":
		handlePostEvent(db, w, userID, event)
	default:
		utils.SendErrorResponse(w, http.StatusBadRequest, "Unknown event type")
	}
}

func handleCommentEvent(db *sql.DB, w http.ResponseWriter, userID string, event ReceiveEvent) {
	var err error

	switch event.ContentType {
	case "like":
		err = services.CommentEventLike(db, userID, event.Id)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = services.CreateNotification(db, userID, "like_comment", event.Id)
	case "dislike":
		err = services.CommentEventDislike(db, userID, event.Id)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = services.CreateNotification(db, userID, "dislike_comment", event.Id)
	default:
		utils.SendErrorResponse(w, http.StatusBadRequest, "Unknown event type")
		return
	}

	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handlePostEvent(db *sql.DB, w http.ResponseWriter, userID string, event ReceiveEvent) {
	var err error

	switch event.ContentType {
	case "like":
		err = services.InsertPostLike(db, userID, event.Id)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "dislike":
		err = services.InsertPostDislike(db, userID, event.Id)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		utils.SendErrorResponse(w, http.StatusBadRequest, "Unknown event type")
		return
	}

	err = services.CreateNotification(db, userID, event.ContentType, event.Id)

	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
