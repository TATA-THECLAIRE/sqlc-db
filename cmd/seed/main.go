package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/ardanlabs/conf/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBUser      string `conf:"env:DB_USER,required"`
	DBPassword  string `conf:"env:DB_PASSWORD,required,mask"`
	DBHost      string `conf:"env:DB_HOST,required"`
	DBPort      uint16 `conf:"env:DB_PORT,required"`
	DBName      string `conf:"env:DB_Name,required"`
	TLSDisabled bool   `conf:"env:DB_TLS_DISABLED"`
}

type Config struct {
	DB DBConfig
}

type QuizData struct {
	Title       string
	Description string
	Questions   []QuestionData
}

type QuestionData struct {
	QuestionText  string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectAnswer string
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	config := Config{}

	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	_, err := conf.Parse("", &config)
	if err != nil {
		return err
	}

	dbConnectionURL := getPostgresConnectionURL(config.DB)
	db, err := pgxpool.New(ctx, dbConnectionURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	querier := repo.New(db)

	// Define quiz data
	quizzes := []QuizData{
		{
			Title:       "Go Programming Basics",
			Description: "Test your knowledge of Go fundamentals",
			Questions: []QuestionData{
				{
					QuestionText:  "What is a goroutine?",
					OptionA:       "A function",
					OptionB:       "A lightweight thread managed by Go runtime",
					OptionC:       "A package",
					OptionD:       "A data structure",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "Which keyword is used to define a constant in Go?",
					OptionA:       "const",
					OptionB:       "var",
					OptionC:       "let",
					OptionD:       "final",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "What does the 'defer' keyword do?",
					OptionA:       "Delays execution permanently",
					OptionB:       "Executes a function after surrounding function returns",
					OptionC:       "Cancels function execution",
					OptionD:       "Creates a new thread",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "How do you create a slice in Go?",
					OptionA:       "var s []int",
					OptionB:       "var s [int]",
					OptionC:       "slice s int[]",
					OptionD:       "new slice(int)",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "What is the zero value of a pointer in Go?",
					OptionA:       "0",
					OptionB:       "null",
					OptionC:       "nil",
					OptionD:       "undefined",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "Which of these is NOT a valid Go data type?",
					OptionA:       "int64",
					OptionB:       "float32",
					OptionC:       "decimal",
					OptionD:       "complex128",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "What does the 'range' keyword do?",
					OptionA:       "Creates a range of numbers",
					OptionB:       "Iterates over elements in various data structures",
					OptionC:       "Defines a numeric range type",
					OptionD:       "Limits variable scope",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "How do you check if a key exists in a map?",
					OptionA:       "value := map[key]",
					OptionB:       "value, exists := map[key]",
					OptionC:       "exists := map.has(key)",
					OptionD:       "value, ok := map[key]",
					CorrectAnswer: "D",
				},
				{
					QuestionText:  "What is the purpose of the 'interface{}' type?",
					OptionA:       "To define abstract methods",
					OptionB:       "To represent any type (empty interface)",
					OptionC:       "To create network interfaces",
					OptionD:       "To define GUI interfaces",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "Which command builds a Go program?",
					OptionA:       "go compile",
					OptionB:       "go make",
					OptionC:       "go build",
					OptionD:       "go create",
					CorrectAnswer: "C",
				},
			},
		},
		{
			Title:       "PostgreSQL Fundamentals",
			Description: "Test your PostgreSQL database knowledge",
			Questions: []QuestionData{
				{
					QuestionText:  "What does SQL stand for?",
					OptionA:       "Structured Query Language",
					OptionB:       "Simple Question Language",
					OptionC:       "System Query Logic",
					OptionD:       "Standard Quality Language",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "Which command is used to retrieve data from a database?",
					OptionA:       "GET",
					OptionB:       "FETCH",
					OptionC:       "SELECT",
					OptionD:       "RETRIEVE",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "What is a PRIMARY KEY?",
					OptionA:       "The first column in a table",
					OptionB:       "A unique identifier for a record",
					OptionC:       "The most important column",
					OptionD:       "A password field",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "Which JOIN returns all records from both tables?",
					OptionA:       "INNER JOIN",
					OptionB:       "LEFT JOIN",
					OptionC:       "RIGHT JOIN",
					OptionD:       "FULL OUTER JOIN",
					CorrectAnswer: "D",
				},
				{
					QuestionText:  "What does the DISTINCT keyword do?",
					OptionA:       "Removes duplicate rows",
					OptionB:       "Sorts the results",
					OptionC:       "Counts unique values",
					OptionD:       "Creates a unique index",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "Which clause filters grouped results?",
					OptionA:       "WHERE",
					OptionB:       "HAVING",
					OptionC:       "FILTER",
					OptionD:       "GROUP BY",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is a FOREIGN KEY?",
					OptionA:       "A key from another database",
					OptionB:       "A key that references a primary key in another table",
					OptionC:       "An encrypted key",
					OptionD:       "A backup key",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "Which command adds a new record to a table?",
					OptionA:       "ADD",
					OptionB:       "INSERT",
					OptionC:       "CREATE",
					OptionD:       "APPEND",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What does ACID stand for in database transactions?",
					OptionA:       "Atomic, Consistent, Isolated, Durable",
					OptionB:       "Automatic, Controlled, Indexed, Dynamic",
					OptionC:       "Advanced, Combined, Integrated, Defined",
					OptionD:       "All, Clear, Important, Data",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "Which command removes all records from a table?",
					OptionA:       "DELETE *",
					OptionB:       "REMOVE ALL",
					OptionC:       "TRUNCATE",
					OptionD:       "DROP TABLE",
					CorrectAnswer: "C",
				},
			},
		},
		{
			Title:       "General Programming Concepts",
			Description: "Test your general programming knowledge",
			Questions: []QuestionData{
				{
					QuestionText:  "What is Big O notation used for?",
					OptionA:       "Measuring code quality",
					OptionB:       "Describing algorithm time/space complexity",
					OptionC:       "Defining variable scope",
					OptionD:       "Error handling",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is recursion?",
					OptionA:       "A loop that runs forever",
					OptionB:       "A function that calls itself",
					OptionC:       "A type of variable",
					OptionD:       "An error handling technique",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What does API stand for?",
					OptionA:       "Application Programming Interface",
					OptionB:       "Advanced Program Integration",
					OptionC:       "Automated Process Interaction",
					OptionD:       "Application Process Interface",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "What is the purpose of version control?",
					OptionA:       "To control software versions for sale",
					OptionB:       "To track and manage changes to code",
					OptionC:       "To verify code correctness",
					OptionD:       "To control access to applications",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is a stack data structure?",
					OptionA:       "First In First Out (FIFO)",
					OptionB:       "Last In First Out (LIFO)",
					OptionC:       "Random access",
					OptionD:       "Sorted access",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What does REST stand for in RESTful APIs?",
					OptionA:       "Remote Execution State Transfer",
					OptionB:       "Representational State Transfer",
					OptionC:       "Resource Execution State Transfer",
					OptionD:       "Remote State Execution Tool",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is the purpose of unit testing?",
					OptionA:       "To test the entire application",
					OptionB:       "To test individual components in isolation",
					OptionC:       "To test user interface",
					OptionD:       "To test database connections",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is a binary search?",
					OptionA:       "Searching through binary files",
					OptionB:       "Searching by dividing sorted data in half repeatedly",
					OptionC:       "Searching using two variables",
					OptionD:       "Searching with 0s and 1s",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What does DRY principle mean?",
					OptionA:       "Don't Repeat Yourself",
					OptionB:       "Do Review Yearly",
					OptionC:       "Debug Run Yearly",
					OptionD:       "Don't Rush Yourself",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "What is polymorphism in OOP?",
					OptionA:       "Multiple inheritance",
					OptionB:       "Objects taking many forms/behaving differently",
					OptionC:       "Multiple constructors",
					OptionD:       "Dynamic memory allocation",
					CorrectAnswer: "B",
				},
			},
		},
		{
			Title:       "Web Development Basics",
			Description: "Test your web development knowledge",
			Questions: []QuestionData{
				{
					QuestionText:  "What does HTTP stand for?",
					OptionA:       "HyperText Transfer Protocol",
					OptionB:       "High Transfer Text Protocol",
					OptionC:       "HyperText Transmission Process",
					OptionD:       "Home Tool Transfer Protocol",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "Which HTTP method is used to retrieve data?",
					OptionA:       "POST",
					OptionB:       "PUT",
					OptionC:       "GET",
					OptionD:       "DELETE",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "What is JSON?",
					OptionA:       "JavaScript Object Notation",
					OptionB:       "Java Standard Object Notation",
					OptionC:       "JavaScript Online Network",
					OptionD:       "Java Serialized Object Network",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "What HTTP status code indicates success?",
					OptionA:       "404",
					OptionB:       "500",
					OptionC:       "200",
					OptionD:       "301",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "What does CSS stand for?",
					OptionA:       "Computer Style Sheets",
					OptionB:       "Cascading Style Sheets",
					OptionC:       "Creative Style System",
					OptionD:       "Colorful Style Sheets",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is CORS?",
					OptionA:       "Cross-Origin Resource Sharing",
					OptionB:       "Central Origin Resource System",
					OptionC:       "Cross-Object Reference System",
					OptionD:       "Core Origin Resource Sharing",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "Which tag is used for the largest heading in HTML?",
					OptionA:       "<heading>",
					OptionB:       "<h6>",
					OptionC:       "<h1>",
					OptionD:       "<head>",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "What does DOM stand for?",
					OptionA:       "Document Object Model",
					OptionB:       "Data Object Management",
					OptionC:       "Digital Online Media",
					OptionD:       "Document Oriented Model",
					CorrectAnswer: "A",
				},
				{
					QuestionText:  "What is a cookie in web development?",
					OptionA:       "A sweet snack",
					OptionB:       "Small piece of data stored in the browser",
					OptionC:       "A type of server",
					OptionD:       "A JavaScript library",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is HTTPS?",
					OptionA:       "HTTP with extra speed",
					OptionB:       "HTTP with security (SSL/TLS)",
					OptionC:       "High Transfer Protocol System",
					OptionD:       "HTTP with special features",
					CorrectAnswer: "B",
				},
			},
		},
		{
			Title:       "Data Structures Quiz",
			Description: "Test your knowledge of data structures",
			Questions: []QuestionData{
				{
					QuestionText:  "What is the time complexity of accessing an array element by index?",
					OptionA:       "O(n)",
					OptionB:       "O(log n)",
					OptionC:       "O(1)",
					OptionD:       "O(n²)",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "Which data structure uses LIFO?",
					OptionA:       "Queue",
					OptionB:       "Stack",
					OptionC:       "Linked List",
					OptionD:       "Tree",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is a hash table?",
					OptionA:       "A sorted array",
					OptionB:       "A data structure using key-value pairs with hash function",
					OptionC:       "A type of tree",
					OptionD:       "A linear list",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "In a binary tree, how many children can each node have?",
					OptionA:       "At most 1",
					OptionB:       "At most 2",
					OptionC:       "Exactly 2",
					OptionD:       "Unlimited",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is a linked list?",
					OptionA:       "An array with links",
					OptionB:       "A sequence of nodes where each node points to the next",
					OptionC:       "A list stored in links",
					OptionD:       "A circular array",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is a priority queue?",
					OptionA:       "A queue sorted by time",
					OptionB:       "A queue where elements are served based on priority",
					OptionC:       "The first queue in a system",
					OptionD:       "A fast queue implementation",
					CorrectAnswer: "B",
				},
				{
					QuestionText:  "What is the worst-case time complexity of quicksort?",
					OptionA:       "O(n)",
					OptionB:       "O(n log n)",
					OptionC:       "O(n²)",
					OptionD:       "O(log n)",
					CorrectAnswer: "C",
				},
				{
					QuestionText:  "What is a graph?",
					OptionA:       "A chart showing data",
					OptionB:       "A collection of nodes connected by edges",
					OptionC:       "A type of tree",
					OptionD:       "A sorted array",
					CorrectAnswer: "B",
				},
			},
		},
	}

	// Seed the database
	fmt.Println("Starting database seeding...")
	for i, quizData := range quizzes {
		fmt.Printf("\n[%d/%d] Creating quiz: %s\n", i+1, len(quizzes), quizData.Title)

		// Create quiz
		quiz, err := querier.CreateQuiz(ctx, repo.CreateQuizParams{
			Title:       quizData.Title,
			Description: quizData.Description,
		})
		if err != nil {
			return fmt.Errorf("failed to create quiz: %w", err)
		}

		fmt.Printf("  ✓ Quiz created with ID: %s\n", quiz.ID)

		// Create questions
		for j, q := range quizData.Questions {
			_, err := querier.CreateQuestion(ctx, repo.CreateQuestionParams{
				QuizID:        quiz.ID,
				QuestionText:  q.QuestionText,
				OptionA:       q.OptionA,
				OptionB:       q.OptionB,
				OptionC:       q.OptionC,
				OptionD:       q.OptionD,
				CorrectAnswer: q.CorrectAnswer,
			})
			if err != nil {
				return fmt.Errorf("failed to create question: %w", err)
			}
			fmt.Printf("  ✓ Question %d/%d created\n", j+1, len(quizData.Questions))
		}
	}

	fmt.Printf("\n✅ Successfully seeded %d quizzes!\n", len(quizzes))
	return nil
}

func getPostgresConnectionURL(config DBConfig) string {
	queryValues := url.Values{}
	if config.TLSDisabled {
		queryValues.Add("sslmode", "disable")
	} else {
		queryValues.Add("sslmode", "require")
	}

	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.DBUser, config.DBPassword),
		Host:     fmt.Sprintf("%s:%d", config.DBHost, config.DBPort),
		Path:     config.DBName,
		RawQuery: queryValues.Encode(),
	}

	return dbURL.String()
}