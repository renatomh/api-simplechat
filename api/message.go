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

type createMessageRequest struct {
	ChatID int64  `json:"chat_id" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

func (server *Server) createMessage(ctx *gin.Context) {
	var req createMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Getting user which requested the connection
	user, _ := server.store.GetUserByUsername(ctx, authPayload.Username)

	// Checking if chat exists
	chat, err := server.store.GetChat(ctx, req.ChatID)
	if err != nil {
		// If no item was found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Checking if user is trying to send a message in a chat where it does not take part
	if (user.ID != chat.FromUserID) && (user.ID != chat.ToUserID) {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("cannot send a message to a chat that is not yours")))
		return
	}

	// Defining user who will receive the message
	var toUserId int64
	if user.ID == chat.FromUserID {
		toUserId = chat.ToUserID
	} else {
		toUserId = chat.FromUserID
	}
	arg := db.CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: user.ID,
		ToUserID:   toUserId,
		Body:       req.Body,
	}

	message, err := server.store.CreateMessage(ctx, arg)
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

	ctx.JSON(http.StatusOK, message)
}

type listMessageRequest struct {
	ChatID   int64 `form:"chat_id" binding:"required,min=1"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listMessage(ctx *gin.Context) {
	var req listMessageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Getting the user which made the request
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Querying the user item by the username
	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)

	// Checking if chat exists
	chat, err := server.store.GetChat(ctx, req.ChatID)
	if err != nil {
		// If no item was found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Checking if user is trying to read messages from a chat where it does not take part
	if (user.ID != chat.FromUserID) && (user.ID != chat.ToUserID) {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("cannot read messages from a chat that is not yours")))
		return
	}

	arg := db.ListMessagesParams{
		ChatID: chat.ID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	messages, err := server.store.ListMessages(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
