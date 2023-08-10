package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/token"
)

type createChatRequest struct {
	ContactID int64 `json:"contact_id" binding:"required"`
}

func (server *Server) createChat(ctx *gin.Context) {
	var req createChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Getting user which requested the connection
	user, _ := server.store.GetUserByUsername(ctx, authPayload.Username)

	// Checking if contact exists
	contact, err := server.store.GetContact(ctx, req.ContactID)
	if err != nil {
		// If no item was found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Checking if user is trying to start a chat with someone else's contact
	if (user.ID != contact.FromUserID) && (user.ID != contact.ToUserID) {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("cannot start chat with contact not in your list")))
		return
	}

	// Checking if user is trying to start a chat with a contact which wasn't accepted yet
	if contact.Status != "Accepted" {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("cannot start chat with pending contact")))
		return
	}

	// Checking if chat already exists
	existingChat, err := server.store.GetChatByUserIDs(
		ctx,
		db.GetChatByUserIDsParams{
			FromUserID: contact.FromUserID,
			ToUserID:   contact.ToUserID,
		})
	if err == nil {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("there's already an existing chat with this contact %d", existingChat.ID)))
		return
	}

	arg := db.CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	}

	chat, err := server.store.CreateChat(ctx, arg)
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

	ctx.JSON(http.StatusOK, chat)
}

type listChatRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listChat(ctx *gin.Context) {
	var req listChatRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Querying the user item by the username
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	arg := db.ListChatsParams{
		FromUserID: user.ID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}
	chats, err := server.store.ListChats(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, chats)
}
