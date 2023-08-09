package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/token"
)

type createContactRequest struct {
	Username string `json:"username" binding:"required"`
}

type contactResponse struct {
	ID int64 `json:"id"`
	// The from/to order makes no difference here
	FromUserID       int64  `json:"from_user_id"`
	FromUserUsername string `json:"from_user_username"`
	// The from/to order makes no difference here
	ToUserID       int64  `json:"to_user_id"`
	ToUserUsername string `json:"to_user_username"`
	// Pending, Accepted or Rejected
	Status      string       `json:"status"`
	RequestedAt time.Time    `json:"requested_at"`
	AcceptedAt  sql.NullTime `json:"accepted_at"`
}

func newContactResponse(contact db.Contact, fromUser, toUser db.User) contactResponse {
	return contactResponse{
		ID:               contact.ID,
		FromUserID:       contact.FromUserID,
		FromUserUsername: fromUser.Username,
		ToUserID:         contact.ToUserID,
		ToUserUsername:   toUser.Username,
		Status:           contact.Status,
		RequestedAt:      contact.RequestedAt,
		AcceptedAt:       contact.AcceptedAt,
	}
}

func (server *Server) createContact(ctx *gin.Context) {
	var req createContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Getting user which requested the connection
	fromUser, _ := server.store.GetUserByUsername(ctx, authPayload.Username)

	// Checking if user exists
	toUser, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// If no item was found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Checking if user is trying to connect with himself
	if fromUser.Username == toUser.Username {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("a user cannot connect with itself, you fool!")))
		return
	}

	// Checking if there's already a connection between the users
	existingContact, err := server.store.CheckExistingContact(
		ctx,
		db.CheckExistingContactParams{
			FromUserID: fromUser.ID,
			ToUserID:   toUser.ID,
		},
	)
	if len(existingContact) > 0 {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("there's already a connection between the users")))
		return
	}

	arg := db.CreateContactParams{
		FromUserID: fromUser.ID,
		ToUserID:   toUser.ID,
	}

	contact, err := server.store.CreateContact(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	rsp := newContactResponse(contact, fromUser, toUser)
	ctx.JSON(http.StatusOK, rsp)
}

type listContactRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listContact(ctx *gin.Context) {
	var req listContactRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Querying the user item by the username
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	arg := db.ListContactsParams{
		FromUserID: user.ID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	contacts, err := server.store.ListContacts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

func (server *Server) listPendingContact(ctx *gin.Context) {
	var req listContactRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Querying the user item by the username
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	arg := db.ListPendingContactsParams{
		FromUserID: user.ID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	contacts, err := server.store.ListPendingContacts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

func (server *Server) listAcceptedContact(ctx *gin.Context) {
	var req listContactRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Querying the user item by the username
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	arg := db.ListAcceptedContactsParams{
		FromUserID: user.ID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	contacts, err := server.store.ListAcceptedContacts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

func (server *Server) listRejectedContact(ctx *gin.Context) {
	var req listContactRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Querying the user item by the username
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	arg := db.ListRejectedContactsParams{
		FromUserID: user.ID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	contacts, err := server.store.ListRejectedContacts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}

type acceptContactRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) acceptContact(ctx *gin.Context) {
	var req acceptContactRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contact, err := server.store.GetContact(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Getting current user
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Checking if the contact refers to the user trying to accept it
	if user.ID != contact.ToUserID {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("only the requested user can accept the contact request")))
		return
	}

	// Checking if the contact is not already accepted
	if contact.Status == "Accepted" {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("the contact request was already accepted")))
		return
	}

	acceptedContact, err := server.store.AcceptContact(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acceptedContact)
}

type rejectContactRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) rejectContact(ctx *gin.Context) {
	var req rejectContactRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contact, err := server.store.GetContact(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Getting current user
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Checking if the contact refers to the user trying to reject it
	if user.ID != contact.ToUserID {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("only the requested user can reject the contact request")))
		return
	}

	// Checking if the contact is not already rejected
	if contact.Status == "Rejected" {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("the contact request was already rejected")))
		return
	}

	rejectedContact, err := server.store.RejectContact(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rejectedContact)
}
