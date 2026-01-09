package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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

	// Load configuration
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

	// Connect to database
	dbConnectionURL := getPostgresConnectionURL(config.DB)
	db, err := pgxpool.New(ctx, dbConnectionURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	querier := repo.New(db)
	scanner := bufio.NewScanner(os.Stdin)

	// Main menu loop
	for {
		printMainMenu()
		choice := getUserInput(scanner, "Enter your choice: ")

		switch choice {
		case "1":
			err := listQuizzes(ctx, querier, scanner)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "2":
			err := takeQuiz(ctx, querier, scanner)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "3":
			err := viewLeaderboard(ctx, querier, scanner)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "4":
			err := viewMyHistory(ctx, querier, scanner)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "5":
			err := viewGlobalStats(ctx, querier)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "6":
			fmt.Println("\n Thanks for playing! Goodbye!")
			return nil
		default:
			fmt.Println("\n❌ Invalid choice. Please try again.")
		}
	}
}

func printMainMenu() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" QUIZ APPLICATION")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1.  Browse available quizzes")
	fmt.Println("2.  Take a quiz")
	fmt.Println("3.  View leaderboard")
	fmt.Println("4.  View my history")
	fmt.Println("5.  Global statistics")
	fmt.Println("6.  Exit")
	fmt.Println(strings.Repeat("=", 50))
}

func listQuizzes(ctx context.Context, querier repo.Querier, scanner *bufio.Scanner) error {
	quizzes, err := querier.ListQuizzes(ctx)
	if err != nil {
		return err
	}

	if len(quizzes) == 0 {
		fmt.Println("\n📭 No quizzes available yet.")
		return nil
	}

	fmt.Println("\n Available Quizzes:")
	fmt.Println(strings.Repeat("=", 50))
	for i, quiz := range quizzes {
		// Get question count for this quiz
		questions, _ := querier.GetQuestionsByQuizID(ctx, quiz.ID)
		questionCount := len(questions)

		// Get attempt count
		attempts, _ := querier.GetQuizAttemptsByQuizID(ctx, quiz.ID)
		attemptCount := len(attempts)

		fmt.Printf("\n%d. %s\n", i+1, quiz.Title)
		if quiz.Description != "" {
			fmt.Printf("    %s\n", quiz.Description)
		}
		fmt.Printf("   ❓ Questions: %d\n", questionCount)
		fmt.Printf("    Attempts: %d\n", attemptCount)
		if quiz.CreatedAt.Valid {
			fmt.Printf("    Created: %s\n", quiz.CreatedAt.Time.Format("Jan 02, 2006"))
		}
	}
	fmt.Println(strings.Repeat("=", 50))

	return nil
}

func takeQuiz(ctx context.Context, querier repo.Querier, scanner *bufio.Scanner) error {
	// List available quizzes first
	quizzes, err := querier.ListQuizzes(ctx)
	if err != nil {
		return err
	}

	if len(quizzes) == 0 {
		fmt.Println("\n No quizzes available yet.")
		return nil
	}

	fmt.Println("\n Available Quizzes:")
	fmt.Println(strings.Repeat("-", 50))
	for i, quiz := range quizzes {
		questions, _ := querier.GetQuestionsByQuizID(ctx, quiz.ID)
		fmt.Printf("%d. %s (%d questions)\n", i+1, quiz.Title, len(questions))
	}

	// Get quiz selection
	choiceStr := getUserInput(scanner, "\nSelect a quiz (enter number): ")
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(quizzes) {
		return fmt.Errorf("invalid quiz selection")
	}

	selectedQuiz := quizzes[choice-1]

	// Get questions
	questions, err := querier.GetQuestionsByQuizID(ctx, selectedQuiz.ID)
	if err != nil {
		return err
	}

	if len(questions) == 0 {
		fmt.Println("\n❌ This quiz has no questions yet!")
		return nil
	}

	// Get user name
	userName := getUserInput(scanner, "\nEnter your name: ")
	if userName == "" {
		userName = "Anonymous"
	}

	fmt.Printf("\n Starting Quiz: %s\n", selectedQuiz.Title)
	fmt.Printf(" Total Questions: %d\n", len(questions))
	fmt.Println(strings.Repeat("=", 50))

	// Track timing
	startTime := time.Now()

	// Ask questions and collect answers
	score := 0

	for i, q := range questions {
		fmt.Printf("\n❓ Question %d of %d\n", i+1, len(questions))
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println(q.QuestionText)
		fmt.Printf("  A) %s\n", q.OptionA)
		fmt.Printf("  B) %s\n", q.OptionB)
		fmt.Printf("  C) %s\n", q.OptionC)
		fmt.Printf("  D) %s\n", q.OptionD)

		var answer string
		for {
			answer = strings.ToUpper(getUserInput(scanner, "\nYour answer (A/B/C/D): "))
			if answer == "A" || answer == "B" || answer == "C" || answer == "D" {
				break
			}
			fmt.Println("❌ Invalid answer. Please enter A, B, C, or D.")
		}

		// Get full question with correct answer
		fullQuestion, err := querier.GetQuestionByID(ctx, q.ID)
		if err != nil {
			continue
		}

		if answer == fullQuestion.CorrectAnswer {
			score++
			fmt.Println("✅ Correct!")
		} else {
			fmt.Printf("❌ Wrong! The correct answer was %s\n", fullQuestion.CorrectAnswer)
		}
	}

	duration := time.Since(startTime)

	// Save attempt
	_, err = querier.CreateQuizAttempt(ctx, repo.CreateQuizAttemptParams{
		QuizID:         selectedQuiz.ID,
		UserName:       userName,
		Score:          int32(score),
		TotalQuestions: int32(len(questions)),
	})
	if err != nil {
		return err
	}

	// Display results
	percentage := float64(score) / float64(len(questions)) * 100
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" QUIZ COMPLETED!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf(" Player: %s\n", userName)
	fmt.Printf(" Score: %d/%d (%.1f%%)\n", score, len(questions), percentage)
	fmt.Printf("⏱  Time: %s\n", formatDuration(duration))
	
	// Show performance message
	if percentage == 100 {
		fmt.Println(" Perfect score! Outstanding!")
	} else if percentage >= 80 {
		fmt.Println(" Great job! Excellent performance!")
	} else if percentage >= 60 {
		fmt.Println(" Good effort! Keep it up!")
	} else {
		fmt.Println("Keep practicing! You'll improve!")
	}

	// Show ranking
	attempts, _ := querier.GetQuizAttemptsByQuizID(ctx, selectedQuiz.ID)
	for i, attempt := range attempts {
		if attempt.UserName == userName && attempt.Score == int32(score) {
			fmt.Printf(" Your rank: #%d out of %d attempts\n", i+1, len(attempts))
			break
		}
	}

	fmt.Println(strings.Repeat("=", 50))

	return nil
}

func viewLeaderboard(ctx context.Context, querier repo.Querier, scanner *bufio.Scanner) error {
	// List available quizzes
	quizzes, err := querier.ListQuizzes(ctx)
	if err != nil {
		return err
	}

	if len(quizzes) == 0 {
		fmt.Println("\n No quizzes available yet.")
		return nil
	}

	fmt.Println(" Select a quiz to view leaderboard:")
	fmt.Println(strings.Repeat("-", 50))
	for i, quiz := range quizzes {
		fmt.Printf("%d. %s\n", i+1, quiz.Title)
	}

	choiceStr := getUserInput(scanner, "\nSelect a quiz (enter number, or 0 for all): ")
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 0 || choice > len(quizzes) {
		return fmt.Errorf("invalid quiz selection")
	}

	if choice == 0 {
		// Show combined leaderboard
		return showGlobalLeaderboard(ctx, querier)
	}

	selectedQuiz := quizzes[choice-1]

	// Get attempts
	attempts, err := querier.GetQuizAttemptsByQuizID(ctx, selectedQuiz.ID)
	if err != nil {
		return err
	}

	if len(attempts) == 0 {
		fmt.Println(" No attempts yet for this quiz!")
		return nil
	}

	fmt.Printf("\n LEADERBOARD - %s\n", selectedQuiz.Title)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-5s %-25s %-12s %-15s %-10s\n", "Rank", "Player", "Score", "Percentage", "Date")
	fmt.Println(strings.Repeat("-", 70))

	for i, attempt := range attempts {
		percentage := float64(attempt.Score) / float64(attempt.TotalQuestions) * 100
		rank := fmt.Sprintf("#%d", i+1)
		medal := ""
		switch i {
			case 0:
				medal = "🥇"
			case 1:
				medal = "🥈"
			case 2:
				medal = "🥉"
		}


		// Truncate long names
		displayName := attempt.UserName
		if len(displayName) > 25 {
			displayName = displayName[:22] + "..."
		}

		fmt.Printf("%-5s %-25s %-12s %-15s %-10s %s\n",
			rank,
			displayName,
			fmt.Sprintf("%d/%d", attempt.Score, attempt.TotalQuestions),
			fmt.Sprintf("%.1f%%", percentage),
			func() string {
				if attempt.CreatedAt.Valid {
					return attempt.CreatedAt.Time.Format("Jan 02")
				}
				return "N/A"
			}(),
			medal,
		)
	}
	fmt.Println(strings.Repeat("=", 70))

	return nil
}

func showGlobalLeaderboard(ctx context.Context, querier repo.Querier) error {
	quizzes, err := querier.ListQuizzes(ctx)
	if err != nil {
		return err
	}

	type PlayerStats struct {
		Name          string
		TotalScore    int
		TotalQuestions int
		QuizzesTaken  int
	}

	playerMap := make(map[string]*PlayerStats)

	// Aggregate all attempts
	for _, quiz := range quizzes {
		attempts, err := querier.GetQuizAttemptsByQuizID(ctx, quiz.ID)
		if err != nil {
			continue
		}

		for _, attempt := range attempts {
			if stats, exists := playerMap[attempt.UserName]; exists {
				stats.TotalScore += int(attempt.Score)
				stats.TotalQuestions += int(attempt.TotalQuestions)
				stats.QuizzesTaken++
			} else {
				playerMap[attempt.UserName] = &PlayerStats{
					Name:          attempt.UserName,
					TotalScore:    int(attempt.Score),
					TotalQuestions: int(attempt.TotalQuestions),
					QuizzesTaken:  1,
				}
			}
		}
	}

	// Convert to slice and sort
	var players []*PlayerStats
	for _, stats := range playerMap {
		players = append(players, stats)
	}

	// Sort by percentage
	for i := 0; i < len(players)-1; i++ {
		for j := i + 1; j < len(players); j++ {
			pct1 := float64(players[i].TotalScore) / float64(players[i].TotalQuestions)
			pct2 := float64(players[j].TotalScore) / float64(players[j].TotalQuestions)
			if pct2 > pct1 {
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	fmt.Println("\n GLOBAL LEADERBOARD")
	fmt.Println(strings.Repeat("=", 75))
	fmt.Printf("%-5s %-25s %-12s %-15s %-10s\n", "Rank", "Player", "Score", "Avg %", "Quizzes")
	fmt.Println(strings.Repeat("-", 75))

	for i, stats := range players {
		if i >= 20 { // Show top 20
			break
		}

		percentage := float64(stats.TotalScore) / float64(stats.TotalQuestions) * 100
		rank := fmt.Sprintf("#%d", i+1)
		medal := ""
			switch i {
			case 0:
				medal = "🥇"
			case 1:
				medal = "🥈"
			case 2:
				medal = "🥉"
		}


		displayName := stats.Name
		if len(displayName) > 25 {
			displayName = displayName[:22] + "..."
		}

		fmt.Printf("%-5s %-25s %-12s %-15s %-10d %s\n",
			rank,
			displayName,
			fmt.Sprintf("%d/%d", stats.TotalScore, stats.TotalQuestions),
			fmt.Sprintf("%.1f%%", percentage),
			stats.QuizzesTaken,
			medal,
		)
	}
	fmt.Println(strings.Repeat("=", 75))

	return nil
}

func viewMyHistory(ctx context.Context, querier repo.Querier, scanner *bufio.Scanner) error {
	userName := getUserInput(scanner, "Enter your name: ")
	if userName == "" {
		fmt.Println("❌ Name cannot be empty")
		return nil
	}

	quizzes, err := querier.ListQuizzes(ctx)
	if err != nil {
		return err
	}

	type AttemptWithQuiz struct {
		QuizTitle      string
		Score          int32
		TotalQuestions int32
		CreatedAt      time.Time
	}

	var userAttempts []AttemptWithQuiz

	for _, quiz := range quizzes {
		attempts, err := querier.GetQuizAttemptsByQuizID(ctx, quiz.ID)
		if err != nil {
			continue
		}

		for _, attempt := range attempts {
			if attempt.UserName == userName {
				userAttempts = append(userAttempts, AttemptWithQuiz{
					QuizTitle:      quiz.Title,
					Score:          attempt.Score,
					TotalQuestions: attempt.TotalQuestions,
					CreatedAt:      attempt.CreatedAt.Time,
				})
			}
		}
	}

	if len(userAttempts) == 0 {
		fmt.Printf("\n No attempts found for '%s'\n", userName)
		return nil
	}

	fmt.Printf("\n QUIZ HISTORY - %s\n", userName)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-30s %-12s %-15s %-12s\n", "Quiz", "Score", "Percentage", "Date")
	fmt.Println(strings.Repeat("-", 70))

	totalScore := 0
	totalQuestions := 0

	for _, attempt := range userAttempts {
		percentage := float64(attempt.Score) / float64(attempt.TotalQuestions) * 100
		
		quizTitle := attempt.QuizTitle
		if len(quizTitle) > 30 {
			quizTitle = quizTitle[:27] + "..."
		}

		fmt.Printf("%-30s %-12s %-15s %-12s\n",
			quizTitle,
			fmt.Sprintf("%d/%d", attempt.Score, attempt.TotalQuestions),
			fmt.Sprintf("%.1f%%", percentage),
			attempt.CreatedAt.Format("Jan 02, 2006"),
		)

		totalScore += int(attempt.Score)
		totalQuestions += int(attempt.TotalQuestions)
	}

	fmt.Println(strings.Repeat("-", 70))
	overallPct := float64(totalScore) / float64(totalQuestions) * 100
	fmt.Printf("%-30s %-12s %-15s\n",
		"OVERALL",
		fmt.Sprintf("%d/%d", totalScore, totalQuestions),
		fmt.Sprintf("%.1f%%", overallPct),
	)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Total quizzes taken: %d\n", len(userAttempts))

	return nil
}

func viewGlobalStats(ctx context.Context, querier repo.Querier) error {
	quizzes, err := querier.ListQuizzes(ctx)
	if err != nil {
		return err
	}

	totalAttempts := 0
	totalScore := 0
	totalQuestions := 0
	uniquePlayers := make(map[string]bool)

	fmt.Println("\n GLOBAL STATISTICS")
	fmt.Println(strings.Repeat("=", 50))

	for _, quiz := range quizzes {
		attempts, err := querier.GetQuizAttemptsByQuizID(ctx, quiz.ID)
		if err != nil {
			continue
		}

		for _, attempt := range attempts {
			totalAttempts++
			totalScore += int(attempt.Score)
			totalQuestions += int(attempt.TotalQuestions)
			uniquePlayers[attempt.UserName] = true
		}
	}

	fmt.Printf(" Total Quizzes: %d\n", len(quizzes))
	fmt.Printf(" Unique Players: %d\n", len(uniquePlayers))
	fmt.Printf(" Total Attempts: %d\n", totalAttempts)
	
	if totalAttempts > 0 {
		avgPercentage := float64(totalScore) / float64(totalQuestions) * 100
		fmt.Printf(" Average Score: %.1f%%\n", avgPercentage)
		fmt.Printf(" Average Attempts per Quiz: %.1f\n", float64(totalAttempts)/float64(len(quizzes)))
	}

	fmt.Println(strings.Repeat("=", 50))

	// Show quiz popularity
	fmt.Println("\n MOST POPULAR QUIZZES:")
	fmt.Println(strings.Repeat("-", 50))

	type QuizStats struct {
		Title    string
		Attempts int
	}

	var quizStats []QuizStats
	for _, quiz := range quizzes {
		attempts, _ := querier.GetQuizAttemptsByQuizID(ctx, quiz.ID)
		quizStats = append(quizStats, QuizStats{
			Title:    quiz.Title,
			Attempts: len(attempts),
		})
	}

	// Sort by attempts
	for i := 0; i < len(quizStats)-1; i++ {
		for j := i + 1; j < len(quizStats); j++ {
			if quizStats[j].Attempts > quizStats[i].Attempts {
				quizStats[i], quizStats[j] = quizStats[j], quizStats[i]
			}
		}
	}

	for i, stats := range quizStats {
		if i >= 5 { // Top 5
			break
		}
		fmt.Printf("%d. %s (%d attempts)\n", i+1, stats.Title, stats.Attempts)
	}

	fmt.Println(strings.Repeat("=", 50))

	return nil
}

func getUserInput(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
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
