package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	golog "github.com/Vladroon22/GoLog"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/database"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/service"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/utils"
	"github.com/gorilla/mux"
)

type Handlers struct {
	logger *golog.Logger
	repo   *database.Repo
	srv    *service.Service
}

func NewHandlers(l *golog.Logger, r *database.Repo, s *service.Service) *Handlers {
	return &Handlers{
		repo:   r,
		srv:    s,
		logger: l,
	}
}

func (h *Handlers) SignUP(w http.ResponseWriter, r *http.Request) {
	user := h.srv.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}
	if err := utils.Valid(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorln(err)
		return
	}

	if err := h.repo.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{"Response": "User successfully created"})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	user := h.srv.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorln(err)
		return
	}

	if err := utils.Valid(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorln(err)
		return
	}

	id, err := h.repo.Login(user.Password, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	token, err := h.repo.GenerateJWT(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		h.logger.Errorln(err)
		return
	}

	SetCookie(w, "JWT", token, database.TTLofJWT)

	WriteJSON(w, http.StatusOK, map[string]interface{}{"Response": "OK"})
}

func (h *Handlers) Post(w http.ResponseWriter, r *http.Request) {
	in := h.srv.Post

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorln(err)
		return
	}

	if err := h.repo.Posting(&in); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, "Data was posted")
}

func (h *Handlers) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])

	post, err := h.repo.GetPost(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, post)
}

func (h *Handlers) Comment(w http.ResponseWriter, r *http.Request) {
	in := h.srv.Comment

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorln(err)
		return
	}

	if err := h.repo.Comment(&in); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, "Comment left successfully")
}

func (h *Handlers) ReadComm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])

	comment, err := h.repo.GetComment(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, comment)
}

func (h *Handlers) Like(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])
	typeOfCont := ""

	if strings.Contains(r.URL.Path, "/post/") {
		typeOfCont = "post"
	} else if strings.Contains(r.URL.Path, "/comment/") {
		typeOfCont = "comment"
	} else {
		http.Error(w, "Invalid post/comment action", http.StatusBadRequest)
		h.logger.Errorln("Invalid post/comment action")
		return
	}
	h.logger.Infoln(typeOfCont)

	var name struct {
		liker     string `json:"liker"`
		likeOrDis string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorln(err)
		return
	}

	if err := h.repo.LikeIt(name.likeOrDis, name.likeOrDis, typeOfCont, ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, typeOfCont+" was "+name.likeOrDis+"ed")
}

func (h *Handlers) GetLiker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])

	like, err := h.repo.GetLiker(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, like)

}

func (h *Handlers) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("JWT")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Errorln(err)
			return
		}
		claims, err := database.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			h.logger.Errorln(err)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "id", claims.UserID))

		next.ServeHTTP(w, r)
	})
}

func SetCookie(w http.ResponseWriter, cookieName string, cookies string, ttl time.Duration) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookies,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(ttl),
	}
	http.SetCookie(w, cookie)
}

func WriteJSON(w http.ResponseWriter, status int, a any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(a)
}
