package apihttp

import (
	"database/sql"
	nethttp "net/http"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
	"github.com/chgc/golf_team_manager/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	playerService       *service.PlayerService
	sessionService      *service.SessionService
	registrationService *service.RegistrationService
}

func NewHandlers(database *sql.DB) *Handlers {
	playerRepository := repository.NewSQLitePlayerRepository(database)
	sessionRepository := repository.NewSQLiteSessionRepository(database)
	registrationRepository := repository.NewSQLiteRegistrationRepository(database)

	return &Handlers{
		playerService:       service.NewPlayerService(playerRepository),
		sessionService:      service.NewSessionService(sessionRepository),
		registrationService: service.NewRegistrationService(playerRepository, sessionRepository, registrationRepository),
	}
}

func (h *Handlers) CreatePlayer(c *gin.Context) {
	var input domain.PlayerWriteDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	player, err := h.playerService.Create(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusCreated, player)
}

func (h *Handlers) ListPlayers(c *gin.Context) {
	players, err := h.playerService.List(
		c.Request.Context(),
		c.Query("query"),
		c.Query("status"),
	)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, players)
}

func (h *Handlers) GetPlayerByID(c *gin.Context) {
	player, err := h.playerService.GetByID(c.Request.Context(), c.Param("playerId"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, player)
}

func (h *Handlers) CreateSession(c *gin.Context) {
	var input domain.SessionWriteDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	session, err := h.sessionService.Create(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusCreated, session)
}

func (h *Handlers) ListSessions(c *gin.Context) {
	sessions, err := h.sessionService.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, sessions)
}

func (h *Handlers) GetSessionByID(c *gin.Context) {
	session, err := h.sessionService.GetByID(c.Request.Context(), c.Param("sessionId"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, session)
}

func (h *Handlers) UpdatePlayer(c *gin.Context) {
	var input domain.PlayerWriteDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	player, err := h.playerService.Update(c.Request.Context(), c.Param("playerId"), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, player)
}

func (h *Handlers) UpdateSession(c *gin.Context) {
	var input domain.SessionWriteDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	session, err := h.sessionService.Update(c.Request.Context(), c.Param("sessionId"), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, session)
}

func (h *Handlers) CreateRegistration(c *gin.Context) {
	var input domain.RegistrationWriteDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	input.SessionID = c.Param("sessionId")

	registration, err := h.registrationService.Create(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusCreated, registration)
}

func (h *Handlers) ListRegistrationsBySession(c *gin.Context) {
	registrations, err := h.registrationService.ListBySession(c.Request.Context(), c.Param("sessionId"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, registrations)
}

func (h *Handlers) NotImplemented(c *gin.Context) {
	c.JSON(nethttp.StatusNotImplemented, ErrorResponse{
		Error: APIError{
			Code:    "not_implemented",
			Message: "route reserved for a future phase",
		},
	})
}
