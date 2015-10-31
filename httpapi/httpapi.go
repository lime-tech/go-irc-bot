package httpapi

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"go-irc-bot/client"
	"go-irc-bot/config"
	"net/http"
	"strings"
)

type Router struct {
	mux        *mux.Router
	client     *client.Client
	httpConfig *config.HttpApi
}

func getContent(req *http.Request) (string, error) {
	err := req.ParseForm()
	if err != nil {
		return "", err
	}
	var content []string
	if str, ok := req.PostForm["content"]; ok {
		content = str
	}

	return strings.Join(content, " "), nil
}

func (r *Router) PostToChannel(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	content, err := getContent(req)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	(&client.Message{
		Kind:    client.MSG_KIND_CHAN,
		Content: content,
		Channel: "#" + vars["channel"],
	}).Send(r.client.Conn)

	w.WriteHeader(http.StatusOK)
}

func (r *Router) PostToPM(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	content, err := getContent(req)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug(req.PostForm)
	(&client.Message{
		Kind:    client.MSG_KIND_PRIV,
		Content: content,
		Nick:    vars["nick"],
	}).Send(r.client.Conn)

	w.WriteHeader(http.StatusOK)
}

func NewRouter(httpConfig *config.HttpApi, c *client.Client) *mux.Router {
	r := &Router{
		mux:        mux.NewRouter(),
		client:     c,
		httpConfig: httpConfig,
	}

	r.
		mux.HandleFunc("/channel/{channel:[a-zA-Z0-9_-]+}", r.PostToChannel).
		Methods("POST")
	r.
		mux.HandleFunc("/pm/{nick:[a-zA-Z0-9_-]+}", r.PostToPM).
		Methods("POST")

	return r.mux
}
