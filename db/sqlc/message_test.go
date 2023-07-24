package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateMessage(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contact between users
	contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	testQueries.AcceptContact(context.Background(), contact.ID)
	testQueries.AcceptContact(context.Background(), contact.ID)

	// Creating a chat for the contact
	chat, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, chat)

	// Sending messages for the chat
	arg := CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: chat.FromUserID,
		ToUserID:   chat.ToUserID,
		Body:       "Hello!",
	}
	message1, err := testQueries.CreateMessage(context.Background(), arg)
	// We must also update the chat
	testQueries.UpdateChat(context.Background(), message1.ChatID)
	chat, err = testQueries.GetChat(context.Background(), message1.ChatID)

	require.NoError(t, err)
	require.NotEmpty(t, message1)
	require.Equal(t, arg.ChatID, message1.ChatID)
	require.Equal(t, arg.FromUserID, message1.FromUserID)
	require.Equal(t, arg.ToUserID, message1.ToUserID)
	require.Equal(t, arg.Body, message1.Body)
	require.WithinDuration(t, time.Now(), message1.SentAt, time.Second)
	require.WithinDuration(t, chat.LastMessageReceivedAt.Time, message1.SentAt, time.Second)

	// Sending messages for the chat
	arg = CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: chat.ToUserID,
		ToUserID:   chat.FromUserID,
		Body:       "Hi, there!",
	}
	message2, err := testQueries.CreateMessage(context.Background(), arg)
	// We must also update the chat
	testQueries.UpdateChat(context.Background(), message2.ChatID)
	chat, err = testQueries.GetChat(context.Background(), message2.ChatID)

	require.NoError(t, err)
	require.NotEmpty(t, message1)
	require.Equal(t, arg.ChatID, message2.ChatID)
	require.Equal(t, arg.FromUserID, message2.FromUserID)
	require.Equal(t, arg.ToUserID, message2.ToUserID)
	require.Equal(t, arg.Body, message2.Body)
	require.WithinDuration(t, time.Now(), message2.SentAt, time.Second)
	require.WithinDuration(t, chat.LastMessageReceivedAt.Time, message2.SentAt, time.Second)
}

func TestGetMessage(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contact between users
	contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	testQueries.AcceptContact(context.Background(), contact.ID)
	testQueries.AcceptContact(context.Background(), contact.ID)

	// Creating a chat for the contact
	chat, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, chat)

	// Sending a message for the chat
	arg := CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: chat.FromUserID,
		ToUserID:   chat.ToUserID,
		Body:       "Hello!",
	}
	message, err := testQueries.CreateMessage(context.Background(), arg)
	// We must also update the chat
	testQueries.UpdateChat(context.Background(), message.ChatID)
	chat, err = testQueries.GetChat(context.Background(), message.ChatID)

	require.NoError(t, err)
	require.NotEmpty(t, message)
	require.Equal(t, arg.ChatID, message.ChatID)
	require.Equal(t, arg.FromUserID, message.FromUserID)
	require.Equal(t, arg.ToUserID, message.ToUserID)
	require.Equal(t, arg.Body, message.Body)
	require.WithinDuration(t, time.Now(), message.SentAt, time.Second)
	require.WithinDuration(t, chat.LastMessageReceivedAt.Time, message.SentAt, time.Second)

	// Trying to get the newly created message
	queriedMessage, err := testQueries.GetMessage(context.Background(), message.ID)
	require.NoError(t, err)
	require.NotEmpty(t, queriedMessage)
	require.Equal(t, message.ChatID, queriedMessage.ChatID)
	require.Equal(t, message.FromUserID, queriedMessage.FromUserID)
	require.Equal(t, message.ToUserID, queriedMessage.ToUserID)
	require.Equal(t, message.Body, queriedMessage.Body)
	require.Equal(t, message.SentAt, queriedMessage.SentAt)
	require.WithinDuration(t, chat.LastMessageReceivedAt.Time, queriedMessage.SentAt, time.Second)
}

func TestListMessages(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contact between users
	contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	testQueries.AcceptContact(context.Background(), contact.ID)
	testQueries.AcceptContact(context.Background(), contact.ID)

	// Creating a chat for the contact
	chat, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, chat)

	// Sending messages for the chat
	arg := CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: chat.FromUserID,
		ToUserID:   chat.ToUserID,
		Body:       "Hello!",
	}
	message1, err := testQueries.CreateMessage(context.Background(), arg)
	// We must also update the chat
	testQueries.UpdateChat(context.Background(), message1.ChatID)
	chat, err = testQueries.GetChat(context.Background(), message1.ChatID)
	require.NoError(t, err)
	require.NotEmpty(t, message1)

	// Sending messages for the chat
	arg = CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: chat.ToUserID,
		ToUserID:   chat.FromUserID,
		Body:       "Hi, there!",
	}
	message2, err := testQueries.CreateMessage(context.Background(), arg)
	// We must also update the chat
	testQueries.UpdateChat(context.Background(), message2.ChatID)
	chat, err = testQueries.GetChat(context.Background(), message2.ChatID)
	require.NoError(t, err)
	require.NotEmpty(t, message1)

	// Listing the created messages
	messages, err := testQueries.ListMessages(context.Background(),
		ListMessagesParams{
			ChatID: chat.ID,
			Limit:  10,
			Offset: 0,
		})
	require.NoError(t, err)
	require.Len(t, messages, 2)
	for _, m := range messages {
		require.NotEmpty(t, m)
	}
}

func TestDeleteMessage(t *testing.T) {
	// Creating random users
	users := []User{}
	for i := 0; i < 2; i++ {
		user, _ := createRandomUser(t)
		users = append(users, user)
	}

	// Creating and accepting contact between users
	contact, _ := testQueries.CreateContact(context.Background(), CreateContactParams{
		FromUserID: users[0].ID,
		ToUserID:   users[1].ID,
	})
	testQueries.AcceptContact(context.Background(), contact.ID)
	testQueries.AcceptContact(context.Background(), contact.ID)

	// Creating a chat for the contact
	chat, err := testQueries.CreateChat(context.Background(), CreateChatParams{
		FromUserID: contact.FromUserID,
		ToUserID:   contact.ToUserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, chat)

	// Sending a message for the chat
	arg := CreateMessageParams{
		ChatID:     chat.ID,
		FromUserID: chat.FromUserID,
		ToUserID:   chat.ToUserID,
		Body:       "Hello!",
	}
	message, err := testQueries.CreateMessage(context.Background(), arg)
	// We must also update the chat
	testQueries.UpdateChat(context.Background(), message.ChatID)
	chat, err = testQueries.GetChat(context.Background(), message.ChatID)

	require.NoError(t, err)
	require.NotEmpty(t, message)
	require.Equal(t, arg.ChatID, message.ChatID)
	require.Equal(t, arg.FromUserID, message.FromUserID)
	require.Equal(t, arg.ToUserID, message.ToUserID)
	require.Equal(t, arg.Body, message.Body)
	require.WithinDuration(t, time.Now(), message.SentAt, time.Second)
	require.WithinDuration(t, chat.LastMessageReceivedAt.Time, message.SentAt, time.Second)

	// Trying to delete the chat
	// It shouldn't be possible, since there's a message associated to it
	err = testQueries.DeleteChat(context.Background(), chat.ID)
	require.Error(t, err)

	// Deleting the newly created message
	err = testQueries.DeleteMessage(context.Background(), message.ID)
	require.NoError(t, err)

	// Checking if message was deleted
	deletedMessage, err := testQueries.GetMessage(context.Background(), message.ID)
	require.Error(t, err)
	require.Empty(t, deletedMessage)

	// Now it should be possible to delete the chat
	err = testQueries.DeleteChat(context.Background(), chat.ID)
	require.NoError(t, err)

	// Checking if chat was deleted
	deletedChat, err := testQueries.GetChat(context.Background(), chat.ID)
	require.Error(t, err)
	require.Empty(t, deletedChat)
}
