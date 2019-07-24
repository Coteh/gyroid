package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	pocketActions "gyroid/lib/actions"
	pocketConnector "gyroid/lib/connector"
	models "gyroid/lib/models"
	utils "gyroid/lib/utils"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type PocketAuth struct {
	AccessToken string `json:"access_token"`
}

func loadEnvVars(isVerbose bool) {
	err := godotenv.Load()
	if err != nil && isVerbose {
		log.Println("Error loading .env file. Will rely on system environment variables instead.")
	}
}

func initializePocketConnection() *pocketConnector.PocketClient {
	consumerKey := os.Getenv("CONSUMER_KEY")
	redirectURI := os.Getenv("REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "localhost:8000"
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal("Could not load auth.json file")
	}

	configDir := filepath.Join(usr.HomeDir, ".config", "pocket")

	err = os.MkdirAll(configDir, 0777)
	if err != nil {
		log.Fatal("Could not create directories for config dir")
	}

	pocketAuth := &PocketAuth{}

	var accessToken string

	authFilePath := filepath.Join(configDir, "auth.json")

	err = loadFromJSON(authFilePath, &pocketAuth)
	if err != nil {
		fmt.Println("Could not read auth.json file.")
	} else {
		accessToken = pocketAuth.AccessToken
	}

	pocketClient := pocketConnector.CreatePocketClient(consumerKey, accessToken)

	if accessToken == "" {
		fmt.Println("Creating new access token")
		accessToken, err = pocketActions.PerformAuth(pocketClient, 3000, redirectURI, utils.OpenBrowser)
		if err != nil {
			log.Fatal("Authentication with Pocket failed", err)
		}
		pocketClient.SetAccessToken(accessToken)
		pocketAuth.AccessToken = accessToken
		err := saveToJSON(authFilePath, &pocketAuth)
		if err != nil {
			log.Fatal("Could not save auth.json file")
		}
	}

	return pocketClient
}

func loadPocketArticles(pocketClient *pocketConnector.PocketClient, articlesList *[]models.ArticleResult) {
	var mut sync.Mutex

	// Count should be high to prevent rate limiting for user
	count := 200

	err := pocketActions.GetUntaggedArticles(pocketClient, 0, count, articlesList, &mut)
	if err != nil {
		log.Fatal(err)
	}

	i := count
	for i < 1000 {
		go func() {
			err := pocketActions.GetUntaggedArticles(pocketClient, i, count, articlesList, &mut)
			if err != nil {
				log.Fatal(err)
			}
		}()
		i += count
	}
}

func runArticleLoop(pocketClient *pocketConnector.PocketClient, articlesList *[]models.ArticleResult) {
	articles := *articlesList

	if len(articles) == 0 {
		fmt.Println("No untagged articles")
		return
	}

	i := 0
	isFav := false
	isNext := false
	userMarkedFav := false
	for i < len(articles) {
		article := articles[i]
		if !userMarkedFav && article.Favorite == 1 {
			isFav = true
		}

		unfavStr := ""
		if isFav {
			unfavStr = "un"
		}

		fmt.Printf("-----\n'%s'\n'%s'\n%s\n-----\n[T]ag\t%s[F]avourite\t[B]ump\t[A]rchive\t[D]elete\t[DD]elete with yes\t[O]pen\t[+]Add Article by URL\t[N]ext\t[E]xit\n-----\n",
			article.ResolvedTitle, article.Excerpt, article.ResolvedURL, unfavStr)

		br := bufio.NewReader(os.Stdin)
		input, _ := br.ReadString('\n')

		command := strings.TrimSpace(strings.ToLower(input))

		if command == "" {
			fmt.Println("Please select an option")
		}

		switch command {
		case "t":
			fmt.Println("Enter tags separated by commas (',')")

			br := bufio.NewReader(os.Stdin)
			input, _ = br.ReadString('\n')

			tags := strings.Split(input, ",")

			result, err := pocketActions.MarkArticleWithTag(pocketClient, article.ItemID, tags)
			if result {
				fmt.Println("Tagging success")
			} else {
				if err != nil {
					fmt.Println("Tagging failure: ", err)
				} else {
					fmt.Println("Did not apply any tags")
				}
			}
			break
		case "f":
			if isFav {
				_, err := pocketActions.UnfavouriteArticle(pocketClient, article.ItemID)
				if err != nil {
					fmt.Println("Failure unfavouriting article: ", err)
				} else {
					fmt.Println("Success unfavouriting article")
				}
			} else {
				_, err := pocketActions.FavouriteArticle(pocketClient, article.ItemID)
				if err != nil {
					fmt.Println("Failure favouriting article: ", err)
				} else {
					fmt.Println("Success favouriting article")
				}
			}
			userMarkedFav = true
			isFav = !isFav
			break
		case "b":
			_, err := pocketActions.BumpArticleToTop(pocketClient, article.ItemID)
			if err != nil {
				fmt.Println("Failure bumping article: ", err)
			} else {
				fmt.Println("Success bumping article")
			}
			break
		case "a":
			_, err := pocketActions.ArchiveArticle(pocketClient, article.ItemID)
			if err != nil {
				fmt.Println("Failure archiving article: ", err)
			} else {
				fmt.Println("Success archiving article")
			}
			isNext = true
			break
		case "o":
			utils.OpenBrowser(article.ResolvedURL)
			break
		case "+":
			fmt.Println("Enter the URL to add:")
			url := readUserInput(func(input string) string {
				return strings.TrimSpace(input)
			})
			if !utils.IsURL(url) {
				fmt.Println("Invalid URL")
				break
			}
			result, err := pocketActions.AddArticle(pocketClient, url)
			if err != nil {
				fmt.Println("Error adding article: ", err)
			} else {
				fmt.Printf("Success adding article: '%s' (%s)\n", result["title"], result["resolved_url"])
			}
			break
		case "d":
			fmt.Println("Are you sure you want to delete this article? You won't be able to restore it unless you readd it. [Y/n]")
			br := bufio.NewReader(os.Stdin)
			input, _ := br.ReadString('\n')

			response := strings.TrimSpace(strings.ToLower(input))

			if response != "y" && response != "yes" {
				break
			}
			fallthrough
		case "dd":
			_, err := pocketActions.DeleteArticle(pocketClient, article.ItemID)
			if err != nil {
				fmt.Println("Failure archiving article: ", err)
			} else {
				fmt.Println("Success deleting article")
			}
			fallthrough
		case "n":
			isNext = true
			break
		case "e":
			return
		}
		if isNext {
			i++
			isFav = false
			isNext = false
			userMarkedFav = false
		}
	}
	fmt.Println("End of Pocket list")
}

func main() {
	isVerbose := false
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-v" || os.Args[i] == "--verbose" {
			isVerbose = true
		}
	}

	loadEnvVars(isVerbose)

	pocketClient := initializePocketConnection()

	articles := make([]models.ArticleResult, 0, 10)

	loadPocketArticles(pocketClient, &articles)

	runArticleLoop(pocketClient, &articles)
}

func loadFromJSON(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Could not open file '%s'", path)
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func saveToJSON(path string, v interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Errorf("Could not open file '%s'", path)
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(v)
}

func readUserInput(strFunc func(string) string) string {
	br := bufio.NewReader(os.Stdin)
	input, _ := br.ReadString('\n')

	return strFunc(input)
}
