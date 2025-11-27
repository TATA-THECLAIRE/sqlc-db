# HOW TO RUN
**clone repo:** `git clone https://github.com/TATA-THECLAIRE/sqlc-db` 
**enter the folder:** `cd sqlc-db` 
**download dependencies:** `go mod tidy` 
**edit env:** `create database and change the the .env.example to .env and make sure the credential there  match that of the created database` 
**run server:** `go run cmd/api/main.go` 
**seed database:** `go run cmd/seed/main.go` 

# Quiz API — Postman Testing Guide

##  Setup

**Base URL:** `http://localhost:8085` 

> Make sure the API server is running:

```bash
go run cmd/api/main.go
```

---

## 1️⃣ Create a Quiz

**Method:** `POST`
**URL:** `{{base_url}}/quizzes`
**Headers:**

```
Content-Type: application/json
```

**Body (JSON):** `example as below`

```json
{
  "title": "General Knowledge Quiz",
  "description": "Test your general knowledge"
}
```

**Expected Response (200 OK):**

```json
{
  "id": "some-generated-id",
  "title": "General Knowledge Quiz",
  "description": "Test your general knowledge",
  "created_at": "2024-11-26T10:00:00Z"
}
```

**Note:** copy the `id` from the response — you'll need it for later requests.

---

## 2️⃣ Get All Quizzes

**Method:** `GET`
**URL:** `{{base_url}}/quizzes`

**Expected Response (200 OK):**

```json
[
  {
    "id": "quiz-id-1",
    "title": "General Knowledge Quiz",
    "description": "Test your general knowledge",
    "created_at": "2024-11-26T10:00:00Z"
  }
]
```

---

## 3️⃣ Get Single Quiz by ID

**Method:** `GET`
**URL:** `{{base_url}}/quizzes/{{quiz_id}}`

**Expected Response (200 OK):**

```json
{
  "id": "your-quiz-id-here",
  "title": "General Knowledge Quiz",
  "description": "Test your general knowledge",
  "created_at": "2024-11-26T10:00:00Z"
}
```

---

## 4️⃣ Create Questions

**Method:** `POST`
**URL:** `{{base_url}}/questions`
**Headers:**

```
Content-Type: application/json
```

**Body (JSON) — Example Question 1:**

```json
{
  "quiz_id": "{{quiz_id}}",
  "question_text": "What is the capital of France?",
  "option_a": "London",
  "option_b": "Paris",
  "option_c": "Berlin",
  "option_d": "Madrid",
  "correct_answer": "B"
}
```

**Expected Response (200 OK):**

```json
{
  "id": "question-id-1",
  "quiz_id": "{{quiz_id}}",
  "question_text": "What is the capital of France?",
  "option_a": "London",
  "option_b": "Paris",
  "option_c": "Berlin",
  "option_d": "Madrid",
  "correct_answer": "B",
  "created_at": "2024-11-26T10:05:00Z"
}
```

Create more questions by repeating the request with different bodies (examples shown below).

**Question 2:**

```json
{
  "quiz_id": "{{quiz_id}}",
  "question_text": "What is 2 + 2?",
  "option_a": "3",
  "option_b": "4",
  "option_c": "5",
  "option_d": "6",
  "correct_answer": "B"
}
```

**Question 3:**

```json
{
  "quiz_id": "{{quiz_id}}",
  "question_text": "Which planet is known as the Red Planet?",
  "option_a": "Venus",
  "option_b": "Jupiter",
  "option_c": "Mars",
  "option_d": "Saturn",
  "correct_answer": "C"
}
```

**Note:** Save question IDs from responses — you will use them for attempts.

---

## 5️⃣ Get Questions for a Quiz

**Method:** `GET`
**URL:** `{{base_url}}/quizzes/{{quiz_id}}/questions`

**Expected Response (200 OK):**

```json
[
  {
    "id": "question-id-1",
    "quiz_id": "{{quiz_id}}",
    "question_text": "What is the capital of France?",
    "option_a": "London",
    "option_b": "Paris",
    "option_c": "Berlin",
    "option_d": "Madrid",
    "created_at": "2024-11-26T10:05:00Z"
  },
  {
    "id": "question-id-2",
    "quiz_id": "{{quiz_id}}",
    "question_text": "What is 2 + 2?",
    "option_a": "3",
    "option_b": "4",
    "option_c": "5",
    "option_d": "6",
    "created_at": "2024-11-26T10:06:00Z"
  }
]
```

---

## 6️⃣ Submit Quiz Attempt (Take the Quiz)

**Method:** `POST`
**URL:** `{{base_url}}/attempts`
**Headers:**

```
Content-Type: application/json
```

**Body (JSON):**

```json
{
  "quiz_id": "{{quiz_id}}",
  "user_name": "John Doe",
  "answers": {
    "question-id-1": "B",
    "question-id-2": "B",
    "question-id-3": "C"
  }
}
```

**Expected Response (200 OK):**

```json
{
  "attempt": {
    "id": "attempt-id-1",
    "quiz_id": "{{quiz_id}}",
    "user_name": "John Doe",
    "score": 3,
    "total_questions": 3,
    "created_at": "2024-11-26T10:15:00Z"
  },
  "results": [
    { "question_id": "question-id-1", "correct": true },
    { "question_id": "question-id-2", "correct": true },
    { "question_id": "question-id-3", "correct": true }
  ]
}
```

**Test with wrong answers:**

```json
{
  "quiz_id": "{{quiz_id}}",
  "user_name": "Jane Smith",
  "answers": {
    "question-id-1": "A",
    "question-id-2": "B",
    "question-id-3": "A"
  }
}
```

---

## 7️⃣ Get Leaderboard

**Method:** `GET`
**URL:** `{{base_url}}/leaderboard/{{quiz_id}}`

**Expected Response (200 OK):**

```json
[
  {
    "id": "attempt-id-1",
    "quiz_id": "{{quiz_id}}",
    "user_name": "John Doe",
    "score": 3,
    "total_questions": 3,
    "created_at": "2024-11-26T10:15:00Z"
  },
  {
    "id": "attempt-id-2",
    "quiz_id": "{{quiz_id}}",
    "user_name": "Jane Smith",
    "score": 1,
    "total_questions": 3,
    "created_at": "2024-11-26T10:16:00Z"
  }
]
```

##  OPTIONAL
`To practice the database queries more i made a small cli app you can run it and actually do the quizes `
**run:** `go run cmd/main.go`
**sql operation used**
* COUNT Operations - Counting questions, attempts, unique players
* SUM Operations - Totaling scores and questions answered
* AVG Calculations - Average score percentages, attempts per quiz
* GROUP BY - Grouping attempts by player name
* ORDER BY / RANKING - Sorting leaderboards and showing ranks
* FILTER Operations - Filtering by user name and quiz ID
* MIN/MAX - Finding top performers (medals)
* JOIN Operations - Combining data from multiple tables


