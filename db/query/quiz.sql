-- name: CreateQuiz :one
INSERT INTO quizzes (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetQuizByID :one
SELECT * FROM quizzes
WHERE id = $1;

-- name: ListQuizzes :many
SELECT * FROM quizzes
ORDER BY created_at DESC;

-- name: CreateQuestion :one
INSERT INTO questions (quiz_id, question_text, option_a, option_b, option_c, option_d, correct_answer)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetQuestionsByQuizID :many
SELECT id, quiz_id, question_text, option_a, option_b, option_c, option_d, created_at
FROM questions
WHERE quiz_id = $1
ORDER BY created_at;

-- name: GetQuestionByID :one
SELECT * FROM questions
WHERE id = $1;

-- name: CreateQuizAttempt :one
INSERT INTO quiz_attempts (quiz_id, user_name, score, total_questions)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetQuizAttemptsByQuizID :many
SELECT * FROM quiz_attempts
WHERE quiz_id = $1
ORDER BY score DESC, created_at DESC;