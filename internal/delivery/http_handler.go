// internal/delivery/http_handler.go
package delivery

import (
	"context"
	"net/http"

	"telegram-chat-analyzer/internal/domain"
	"telegram-chat-analyzer/internal/repository"
	"telegram-chat-analyzer/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	usecase    usecase.MessageUsecase
	repo       repository.MongoRepository
	collection string
}

func NewMessageHandler(router *gin.Engine, uc usecase.MessageUsecase, repo repository.MongoRepository, collection string) {
	handler := &MessageHandler{
		usecase:    uc,
		repo:       repo,
		collection: collection,
	}
	router.POST("/topSixWords", handler.ProcessMessages)                                     // return top 6 frequent words
	router.POST("/countMessages", handler.CountMessages)                                     // count total messages sent by each person
	router.POST("/countWords", handler.countWords)                                           // return shared interests
	router.POST("/totalDaysTalked", handler.totalDaysTalked)                                 // return total active days
	router.POST("/messagesPerDay", handler.MessagesPerDay)                                   // return each day messages
	router.POST("/averageMessagesPerDay", handler.AverageMessagesPerDay)                     // return average messages per day for each person
	router.POST("/weeklyStats", handler.WeeklyStats)                                         // return number of messages on each day of the week
	router.POST("/hourlyStats", handler.hourlyStats)                                         // return number of messages on each hour of the day
	router.POST("/mostActiveDayOfWeek", handler.MostActiveDayOfWeek)                         // return most active day of the week
	router.POST("/messageLengthStatistics", handler.MessageLengthStatistics)                 // return average char per text total char max and min
	router.POST("/replyTimeAnalysis", handler.ReplyTimeAnalysis)                             // return the average time taken to reply
	router.POST("/countConversationStartersPerDay", handler.CountConversationStartersPerDay) // return the number of conversation starters per day
	router.POST("/countConsecutiveDays", handler.CountConsecutiveDays)                       // return the number of consecutive days talked
	router.POST("/relationshipScore", handler.RelationshipScore)
	router.POST("/currentStreak", handler.CurrentStreak)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Telegram Chat Analyzer!"})
	})
}

func (h *MessageHandler) ProcessMessages(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}
	topSixWords, err := h.usecase.CountWords(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process messages: " + err.Error()})
		return
	}
	result := gin.H{
		"topSixWords": topSixWords,
	}

	// Save to MongoDB
	if err := h.repo.SaveProcessedData(context.Background(), h.collection, chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database: " + err.Error()})
		return
	}

	// Respond with success and data
	c.JSON(http.StatusOK, gin.H{
		"message":       "Successfully processed messages",
		"processedData": result,
	})
}

func (h *MessageHandler) CountMessages(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}
	personOne, personTwo := h.usecase.GetPersons(chat)
	totalMessageCount, personOneMessageCount, personTwoMessageCount := h.usecase.CountMessages(chat)
	result := gin.H{
		"totalMessageCount": totalMessageCount,
		personOne:           personOneMessageCount,
		personTwo:           personTwoMessageCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Successfully counted messages",
		"processedData": result,
	})
}

func (h *MessageHandler) totalDaysTalked(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	totalDays := h.usecase.TotalDaysTalked(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Successfully calculated total days talked",
		"totalDays": totalDays,
	})
}

func (h *MessageHandler) MessagesPerDay(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	messagesPerDay := h.usecase.MessagesPerDay(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Successfully calculated messages per day",
		"messagesPerDay": messagesPerDay,
	})
}

func (h *MessageHandler) WeeklyStats(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	weeklyStat := h.usecase.WeeklyStats(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully calculated weekly stats",
		"weeklyStats": weeklyStat,
	})
}

func (h *MessageHandler) hourlyStats(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	HourlyStats := h.usecase.HourlyStats(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully calculated hourly stats",
		"hourlyStats": HourlyStats,
	})
}

func (h *MessageHandler) MostActiveDayOfWeek(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	mostActiveDay := h.usecase.MostActiveDayOfWeek(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Successfully found the most active day of the week",
		"mostActiveDay": mostActiveDay,
	})
}

func (h *MessageHandler) MessageLengthStatistics(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	messageLengthStats := h.usecase.MessageLengthStatistics(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":              "Successfully calculated message length statistics",
		"charactersPerMessage": messageLengthStats,
	})
}

func (h *MessageHandler) ReplyTimeAnalysis(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	totalReplyTimes := h.usecase.ReplyTimeAnalysis(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Successfully analyzed reply times",
		"totalReplyTimes": totalReplyTimes,
	})
}

func (h *MessageHandler) CountConversationStartersPerDay(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	conversationStarters, err := h.usecase.CountConversationStartersPerDay(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process messages: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "Successfully counted conversation starters per day",
		"conversationStarters": conversationStarters,
	})
}

func (h *MessageHandler) CountConsecutiveDays(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	consecutiveDays, err := h.usecase.CountConsecutiveDays(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process messages: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":         "Successfully counted consecutive days talked",
		"consecutiveDays": consecutiveDays,
	})
}

func (h *MessageHandler) GetSharedInterests(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	interests := h.usecase.GetSharedInterests(chat)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Successfully found shared interests",
		"interests": interests,
	})
}

func (h *MessageHandler) AverageMessagesPerDay(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	averageMessages := h.usecase.AverageMessagesPerDay(chat)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Successfully calculated average messages per day",
		"averageMessages": averageMessages,
	})
}

func (h *MessageHandler) countWords(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	count, average, err := h.usecase.CountWord(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process messages: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully counted words",
		"count":   count,
		"average": average,
	})
}

func (h *MessageHandler) RelationshipScore(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	score, err := h.usecase.RelationshipScore(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process messages: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully calculated relationship score",
		"score":   score,
	})
}

func (h *MessageHandler) CurrentStreak(c *gin.Context) {
	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	if len(chat.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No messages provided in the input data"})
		return
	}

	streak, err := h.usecase.CurrentStreak(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process messages: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Successfully calculated current streak",
		"currentStreak": streak,
	})
}
