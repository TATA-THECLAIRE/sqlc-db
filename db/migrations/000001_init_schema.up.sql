CREATE TABLE quizzes (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE questions (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
    quiz_id VARCHAR(36) NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    option_a VARCHAR(255) NOT NULL,
    option_b VARCHAR(255) NOT NULL,
    option_c VARCHAR(255) NOT NULL,
    option_d VARCHAR(255) NOT NULL,
    correct_answer VARCHAR(1) NOT NULL CHECK (correct_answer IN ('A', 'B', 'C', 'D')),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE quiz_attempts (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
    quiz_id VARCHAR(36) NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    user_name VARCHAR(100) NOT NULL,
    score INTEGER NOT NULL,
    total_questions INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);