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

-- name: UpdateQuiz :one
UPDATE quizzes 
SET title = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteQuiz :exec
DELETE FROM quizzes
WHERE id = $1;

-- name: UpdateQuestion :one
UPDATE questions
SET question_text = $2, 
    option_a = $3, 
    option_b = $4, 
    option_c = $5, 
    option_d = $6, 
    correct_answer = $7
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE  FROM questions
WHERE id = $1;


-- name: ListQuizAttempts :many
SELECT * FROM quiz_attempts
WHERE quiz_id = $1
ORDER BY created_at DESC;

-- name: GetQuizStats :one
SELECT 
    COUNT(*) as attemps_count,
    AVG(score :: float / total_questions :: float * 100) as avg_score_percent,
    MAX(score) as highest_score,
    MIN(score) as lowest_score
FROM quiz_attempts
WHERE quiz_id = $1;
