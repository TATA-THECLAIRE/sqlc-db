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
// Replace the quizzes slice in your main.go with this:

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
		Title:       "World Geography Challenge",
		Description: "Test your knowledge of world geography",
		Questions: []QuestionData{
			{
				QuestionText:  "What is the capital of Australia?",
				OptionA:       "Sydney",
				OptionB:       "Melbourne",
				OptionC:       "Canberra",
				OptionD:       "Brisbane",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which is the largest ocean on Earth?",
				OptionA:       "Atlantic Ocean",
				OptionB:       "Indian Ocean",
				OptionC:       "Arctic Ocean",
				OptionD:       "Pacific Ocean",
				CorrectAnswer: "D",
			},
			{
				QuestionText:  "How many continents are there?",
				OptionA:       "5",
				OptionB:       "6",
				OptionC:       "7",
				OptionD:       "8",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What is the longest river in the world?",
				OptionA:       "Amazon River",
				OptionB:       "Nile River",
				OptionC:       "Yangtze River",
				OptionD:       "Mississippi River",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which country has the most natural lakes?",
				OptionA:       "United States",
				OptionB:       "Russia",
				OptionC:       "Canada",
				OptionD:       "Brazil",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What is the smallest country in the world?",
				OptionA:       "Monaco",
				OptionB:       "Vatican City",
				OptionC:       "San Marino",
				OptionD:       "Liechtenstein",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which desert is the largest hot desert in the world?",
				OptionA:       "Gobi Desert",
				OptionB:       "Kalahari Desert",
				OptionC:       "Sahara Desert",
				OptionD:       "Arabian Desert",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Mount Everest is located in which mountain range?",
				OptionA:       "Alps",
				OptionB:       "Andes",
				OptionC:       "Himalayas",
				OptionD:       "Rockies",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which country is both in Europe and Asia?",
				OptionA:       "Russia",
				OptionB:       "Turkey",
				OptionC:       "Egypt",
				OptionD:       "Kazakhstan",
				CorrectAnswer: "A",
			},
			{
				QuestionText:  "What is the capital of Canada?",
				OptionA:       "Toronto",
				OptionB:       "Vancouver",
				OptionC:       "Montreal",
				OptionD:       "Ottawa",
				CorrectAnswer: "D",
			},
		},
	},
	{
		Title:       "Science and Nature Quiz",
		Description: "Explore the wonders of science and nature",
		Questions: []QuestionData{
			{
				QuestionText:  "What is the chemical symbol for gold?",
				OptionA:       "Go",
				OptionB:       "Au",
				OptionC:       "Gd",
				OptionD:       "Ag",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "How many bones are in the adult human body?",
				OptionA:       "186",
				OptionB:       "206",
				OptionC:       "226",
				OptionD:       "246",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the speed of light?",
				OptionA:       "300,000 km/s",
				OptionB:       "150,000 km/s",
				OptionC:       "450,000 km/s",
				OptionD:       "200,000 km/s",
				CorrectAnswer: "A",
			},
			{
				QuestionText:  "What is the largest organ in the human body?",
				OptionA:       "Heart",
				OptionB:       "Brain",
				OptionC:       "Liver",
				OptionD:       "Skin",
				CorrectAnswer: "D",
			},
			{
				QuestionText:  "What gas do plants absorb from the atmosphere?",
				OptionA:       "Oxygen",
				OptionB:       "Nitrogen",
				OptionC:       "Carbon Dioxide",
				OptionD:       "Hydrogen",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What is the hardest natural substance on Earth?",
				OptionA:       "Gold",
				OptionB:       "Iron",
				OptionC:       "Diamond",
				OptionD:       "Titanium",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "How many planets are in our solar system?",
				OptionA:       "7",
				OptionB:       "8",
				OptionC:       "9",
				OptionD:       "10",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the powerhouse of the cell?",
				OptionA:       "Nucleus",
				OptionB:       "Ribosome",
				OptionC:       "Mitochondria",
				OptionD:       "Chloroplast",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What is the boiling point of water at sea level?",
				OptionA:       "90°C",
				OptionB:       "100°C",
				OptionC:       "110°C",
				OptionD:       "120°C",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What type of animal is a dolphin?",
				OptionA:       "Fish",
				OptionB:       "Amphibian",
				OptionC:       "Mammal",
				OptionD:       "Reptile",
				CorrectAnswer: "C",
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
	{
		Title:       "Movies and Entertainment",
		Description: "Test your knowledge of movies and pop culture",
		Questions: []QuestionData{
			{
				QuestionText:  "Which movie won the first Academy Award for Best Picture?",
				OptionA:       "The Jazz Singer",
				OptionB:       "Wings",
				OptionC:       "Sunrise",
				OptionD:       "The Broadway Melody",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Who directed 'The Shawshank Redemption'?",
				OptionA:       "Steven Spielberg",
				OptionB:       "Frank Darabont",
				OptionC:       "Christopher Nolan",
				OptionD:       "Martin Scorsese",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What year was the first 'Star Wars' movie released?",
				OptionA:       "1975",
				OptionB:       "1977",
				OptionC:       "1979",
				OptionD:       "1980",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which actor played Iron Man in the Marvel Cinematic Universe?",
				OptionA:       "Chris Evans",
				OptionB:       "Chris Hemsworth",
				OptionC:       "Robert Downey Jr.",
				OptionD:       "Mark Ruffalo",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What is the highest-grossing film of all time (not adjusted for inflation)?",
				OptionA:       "Titanic",
				OptionB:       "Avatar",
				OptionC:       "Avengers: Endgame",
				OptionD:       "Star Wars: The Force Awakens",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which movie features the quote 'Here's looking at you, kid'?",
				OptionA:       "Gone with the Wind",
				OptionB:       "Casablanca",
				OptionC:       "The Maltese Falcon",
				OptionD:       "Citizen Kane",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Who composed the music for 'The Lion King'?",
				OptionA:       "Hans Zimmer",
				OptionB:       "John Williams",
				OptionC:       "Alan Menken",
				OptionD:       "Elton John",
				CorrectAnswer: "D",
			},
			{
				QuestionText:  "Which film won the most Oscars in a single ceremony?",
				OptionA:       "Titanic",
				OptionB:       "Ben-Hur",
				OptionC:       "The Lord of the Rings: Return of the King",
				OptionD:       "All of the above (tied at 11)",
				CorrectAnswer: "D",
			},
			{
				QuestionText:  "What is the name of the fictional African country in Black Panther?",
				OptionA:       "Zamunda",
				OptionB:       "Wakanda",
				OptionC:       "Genovia",
				OptionD:       "Latveria",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which director is known for movies like 'Pulp Fiction' and 'Kill Bill'?",
				OptionA:       "Quentin Tarantino",
				OptionB:       "Wes Anderson",
				OptionC:       "Paul Thomas Anderson",
				OptionD:       "David Fincher",
				CorrectAnswer: "A",
			},
		},
	},
	{
		Title:       "Sports Trivia",
		Description: "Test your sports knowledge across various disciplines",
		Questions: []QuestionData{
			{
				QuestionText:  "How many players are on a soccer team on the field?",
				OptionA:       "9",
				OptionB:       "10",
				OptionC:       "11",
				OptionD:       "12",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which country has won the most FIFA World Cups?",
				OptionA:       "Germany",
				OptionB:       "Argentina",
				OptionC:       "Italy",
				OptionD:       "Brazil",
				CorrectAnswer: "D",
			},
			{
				QuestionText:  "How many Grand Slam tournaments are there in tennis?",
				OptionA:       "3",
				OptionB:       "4",
				OptionC:       "5",
				OptionD:       "6",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the diameter of a basketball hoop in inches?",
				OptionA:       "16 inches",
				OptionB:       "18 inches",
				OptionC:       "20 inches",
				OptionD:       "22 inches",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "In which sport would you perform a 'Fosbury Flop'?",
				OptionA:       "Pole Vault",
				OptionB:       "Long Jump",
				OptionC:       "High Jump",
				OptionD:       "Triple Jump",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "How many rings are on the Olympic flag?",
				OptionA:       "4",
				OptionB:       "5",
				OptionC:       "6",
				OptionD:       "7",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the maximum score in a single frame of bowling?",
				OptionA:       "10",
				OptionB:       "20",
				OptionC:       "30",
				OptionD:       "40",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which athlete has won the most Olympic gold medals?",
				OptionA:       "Usain Bolt",
				OptionB:       "Michael Phelps",
				OptionC:       "Carl Lewis",
				OptionD:       "Simone Biles",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the length of a marathon?",
				OptionA:       "26.2 miles",
				OptionB:       "25 miles",
				OptionC:       "30 miles",
				OptionD:       "24.5 miles",
				CorrectAnswer: "A",
			},
			{
				QuestionText:  "In which sport is the term 'love' used?",
				OptionA:       "Cricket",
				OptionB:       "Tennis",
				OptionC:       "Golf",
				OptionD:       "Badminton",
				CorrectAnswer: "B",
			},
		},
	},
	{
		Title:       "World History Quiz",
		Description: "Journey through significant historical events",
		Questions: []QuestionData{
			{
				QuestionText:  "In what year did World War II end?",
				OptionA:       "1943",
				OptionB:       "1944",
				OptionC:       "1945",
				OptionD:       "1946",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Who was the first president of the United States?",
				OptionA:       "Thomas Jefferson",
				OptionB:       "John Adams",
				OptionC:       "George Washington",
				OptionD:       "Benjamin Franklin",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What year did the Berlin Wall fall?",
				OptionA:       "1987",
				OptionB:       "1989",
				OptionC:       "1991",
				OptionD:       "1993",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which ancient wonder is still standing today?",
				OptionA:       "Colossus of Rhodes",
				OptionB:       "Hanging Gardens of Babylon",
				OptionC:       "Great Pyramid of Giza",
				OptionD:       "Lighthouse of Alexandria",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Who was the first person to walk on the moon?",
				OptionA:       "Buzz Aldrin",
				OptionB:       "Neil Armstrong",
				OptionC:       "Yuri Gagarin",
				OptionD:       "John Glenn",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What year did the Titanic sink?",
				OptionA:       "1910",
				OptionB:       "1911",
				OptionC:       "1912",
				OptionD:       "1913",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which empire built Machu Picchu?",
				OptionA:       "Aztec Empire",
				OptionB:       "Mayan Empire",
				OptionC:       "Inca Empire",
				OptionD:       "Olmec Empire",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Who painted the Mona Lisa?",
				OptionA:       "Michelangelo",
				OptionB:       "Leonardo da Vinci",
				OptionC:       "Raphael",
				OptionD:       "Donatello",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What year did Christopher Columbus reach the Americas?",
				OptionA:       "1490",
				OptionB:       "1492",
				OptionC:       "1494",
				OptionD:       "1496",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Who was known as the 'Iron Lady'?",
				OptionA:       "Angela Merkel",
				OptionB:       "Golda Meir",
				OptionC:       "Margaret Thatcher",
				OptionD:       "Indira Gandhi",
				CorrectAnswer: "C",
			},
		},
	},
	{
		Title:       "Music Knowledge Test",
		Description: "Challenge your music and music history knowledge",
		Questions: []QuestionData{
			{
				QuestionText:  "Which band is known as the 'Fab Four'?",
				OptionA:       "The Rolling Stones",
				OptionB:       "The Beatles",
				OptionC:       "The Who",
				OptionD:       "Led Zeppelin",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Who is known as the 'King of Pop'?",
				OptionA:       "Elvis Presley",
				OptionB:       "Prince",
				OptionC:       "Michael Jackson",
				OptionD:       "David Bowie",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "How many strings does a standard guitar have?",
				OptionA:       "4",
				OptionB:       "5",
				OptionC:       "6",
				OptionD:       "7",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which musical term means to play loudly?",
				OptionA:       "Piano",
				OptionB:       "Forte",
				OptionC:       "Allegro",
				OptionD:       "Adagio",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the best-selling album of all time?",
				OptionA:       "Back in Black",
				OptionB:       "The Dark Side of the Moon",
				OptionC:       "Thriller",
				OptionD:       "The Bodyguard Soundtrack",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which instrument has 88 keys?",
				OptionA:       "Organ",
				OptionB:       "Piano",
				OptionC:       "Harpsichord",
				OptionD:       "Accordion",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Who composed the 'Four Seasons'?",
				OptionA:       "Bach",
				OptionB:       "Mozart",
				OptionC:       "Vivaldi",
				OptionD:       "Beethoven",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What does BPM stand for in music?",
				OptionA:       "Beats Per Minute",
				OptionB:       "Bass Per Measure",
				OptionC:       "Beat Pattern Method",
				OptionD:       "Baseline Per Melody",
				CorrectAnswer: "A",
			},
			{
				QuestionText:  "Which music streaming service was launched first?",
				OptionA:       "Apple Music",
				OptionB:       "Spotify",
				OptionC:       "Tidal",
				OptionD:       "YouTube Music",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What genre of music did Elvis Presley primarily perform?",
				OptionA:       "Jazz",
				OptionB:       "Rock and Roll",
				OptionC:       "Country",
				OptionD:       "Blues",
				CorrectAnswer: "B",
			},
		},
	},
	{
		Title:       "Food and Cuisine Quiz",
		Description: "Test your culinary knowledge from around the world",
		Questions: []QuestionData{
			{
				QuestionText:  "What is the main ingredient in guacamole?",
				OptionA:       "Tomato",
				OptionB:       "Avocado",
				OptionC:       "Lime",
				OptionD:       "Pepper",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "Which country is the origin of sushi?",
				OptionA:       "China",
				OptionB:       "Korea",
				OptionC:       "Japan",
				OptionD:       "Thailand",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What type of pasta is shaped like little ears?",
				OptionA:       "Penne",
				OptionB:       "Farfalle",
				OptionC:       "Orecchiette",
				OptionD:       "Rigatoni",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which spice is the most expensive by weight?",
				OptionA:       "Vanilla",
				OptionB:       "Saffron",
				OptionC:       "Cardamom",
				OptionD:       "Cinnamon",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What is the main ingredient in hummus?",
				OptionA:       "Lentils",
				OptionB:       "Black beans",
				OptionC:       "Chickpeas",
				OptionD:       "Kidney beans",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which cheese is traditionally used on pizza Margherita?",
				OptionA:       "Parmesan",
				OptionB:       "Cheddar",
				OptionC:       "Mozzarella",
				OptionD:       "Gouda",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "What is the base spirit in a Mojito?",
				OptionA:       "Vodka",
				OptionB:       "Tequila",
				OptionC:       "Rum",
				OptionD:       "Gin",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which fruit is used to make traditional wine?",
				OptionA:       "Apples",
				OptionB:       "Grapes",
				OptionC:       "Berries",
				OptionD:       "Peaches",
				CorrectAnswer: "B",
			},
			{
				QuestionText:  "What does 'al dente' mean in cooking?",
				OptionA:       "Fully cooked",
				OptionB:       "Undercooked",
				OptionC:       "Firm to the bite",
				OptionD:       "Crispy",
				CorrectAnswer: "C",
			},
			{
				QuestionText:  "Which country is famous for its chocolate?",
				OptionA:       "France",
				OptionB:       "Switzerland",
				OptionC:       "Germany",
				OptionD:       "Italy",
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