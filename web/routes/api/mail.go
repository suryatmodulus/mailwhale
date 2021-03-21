package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/types/dto"
	"github.com/muety/mailwhale/util"
	"github.com/muety/mailwhale/web/handlers"
	"net/http"
)

const routeMail = "/api/mail"

type MailHandler struct {
	config          *conf.Config
	sendService     *service.SendService
	templateService *service.TemplateService
	clientService   *service.ClientService
	userService     *service.UserService
}

func NewMailHandler() *MailHandler {
	return &MailHandler{
		config:          conf.Get(),
		sendService:     service.NewSendService(),
		templateService: service.NewTemplateService(),
		clientService:   service.NewClientService(),
		userService:     service.NewUserService(),
	}
}

func (h *MailHandler) Register(router *mux.Router) {
	r := router.PathPrefix(routeMail).Subrouter()
	r.Use(
		handlers.NewAuthMiddleware(h.clientService, h.userService, []string{types.PermissionSendMail}),
	)
	r.Path("/").Methods(http.MethodPost).HandlerFunc(h.post)
}

func (h *MailHandler) post(w http.ResponseWriter, r *http.Request) {
	var payload dto.MailSendRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	client := r.Context().Value(conf.KeyClient).(*types.Client)

	sender := payload.From
	if sender == "" {
		sender = client.DefaultSender
	}

	if !client.AllowsSender(sender) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("sender not allowed"))
		return
	}

	mail := &types.Mail{
		From:    sender,
		To:      payload.To,
		Subject: payload.Subject,
	}

	// Case 1: Template ID is given
	if payload.TemplateId != "" {
		template, err := h.templateService.GetById(payload.TemplateId)
		if err != nil {
			util.RespondError(w, r, http.StatusNotFound, err)
			return
		}
		templateData := template.FillContent(payload.TemplateVars)
		if template.GuessIsHtml() {
			mail.WithHTML(templateData)
		} else {
			mail.WithText(templateData)
		}
	}

	// Case 2: Plain text is given
	if mail.Body == "" && payload.Text != "" {
		mail.WithText(payload.Text)
	}

	// Case 3: HTML text is given
	if mail.Body == "" && payload.Html != "" {
		mail.WithHTML(payload.Html)
	}

	if mail.Body == "" {
		util.RespondError(w, r, http.StatusBadRequest, errors.New("empty mail"))
	}

	if err := h.sendService.Send(mail); err != nil {
		util.RespondError(w, r, http.StatusInternalServerError, err)
		return
	}

	util.RespondEmpty(w, r, 0)
}
