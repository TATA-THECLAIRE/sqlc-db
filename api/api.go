package api

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	querier repo.Querier
}

func NewQuizHandler(querier repo.Querier) *QuizHandler {
	return &QuizHandler{
		querier: querier,
	}
}

func (h *QuizHandler) WireHttpHandler() http.Handler {
	r := gin.Default()

	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// Quiz endpoints
	r.GET("/quizzes", h.handleListQuizzes)
	r.POST("/quizzes", h.handleCreateQuiz)
	r.GET("/quizzes/:id", h.handleGetQuiz)
	r.GET("/quizzes/:id/questions", h.handleGetQuizQuestions)

	// Question endpoints
	r.POST("/questions", h.handleCreateQuestion)

	// Attempt endpoints
	r.POST("/attempts", h.handleCreateAttempt)
	r.GET("/leaderboard/:quiz_id", h.handleLeaderboard)

	return r
}

// Quiz handlers
func (h *QuizHandler) handleListQuizzes(c *gin.Context) {
	quizzes, err := h.querier.ListQuizzes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

func (h *QuizHandler) handleCreateQuiz(c *gin.Context) {
	var req repo.CreateQuizParams

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz, err := h.querier.CreateQuiz(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (h *QuizHandler) handleGetQuiz(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	quiz, err := h.querier.GetQuizByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (h *QuizHandler) handleGetQuizQuestions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	questions, err := h.querier.GetQuestionsByQuizID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// Question handlers
func (h *QuizHandler) handleCreateQuestion(c *gin.Context) {
	var req repo.CreateQuestionParams

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.querier.CreateQuestion(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// Attempt handlers
func (h *QuizHandler) handleCreateAttempt(c *gin.Context) {
	var req struct {
		QuizID   string            `json:"quiz_id"`
		UserName string            `json:"user_name"`
		Answers  map[string]string `json:"answers"`
	}

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.QuizID == "" || req.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id and user_name are required"})
		return
	}

	// Get all questions for the quiz
	questions, err := h.querier.GetQuestionsByQuizID(c, req.QuizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found for this quiz"})
		return
	}

	// Calculate score
	score := 0
	results := make([]map[string]interface{}, 0)

	for _, q := range questions {
		// Get the full question with correct answer
		fullQuestion, err := h.querier.GetQuestionByID(c, q.ID)
		if err != nil {
			continue
		}

		userAnswer := req.Answers[q.ID]
		isCorrect := userAnswer == fullQuestion.CorrectAnswer

		if isCorrect {
			score++
		}

		results = append(results, map[string]interface{}{
			"question_id": q.ID,
			"correct":     isCorrect,
		})
	}

	// Save attempt
	attempt, err := h.querier.CreateQuizAttempt(c, repo.CreateQuizAttemptParams{
		QuizID:         req.QuizID,
		UserName:       req.UserName,
		Score:          int32(score),
		TotalQuestions: int32(len(questions)),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attempt": attempt,
		"results": results,
	})
}

func (h *QuizHandler) handleLeaderboard(c *gin.Context) {
	id := c.Param("quiz_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	attempts, err := h.querier.GetQuizAttemptsByQuizID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attempts)
}