package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/models"
	"github.com/Coteh/gyroid/lib/utils"
	"github.com/Coteh/gyroid/lib/config"
)

const CONFIG_SUBFOLDER_NAME = ".config/gyroid"
const CONFIG_FILE_NAME = "config.yml"
const AUTH_FILE_NAME = "auth.json"
const CONSUMER_KEY_FILE_NAME = "consumer_key"

type PocketAuth struct {
	AccessToken string `json:"access_token"`
}

func loadConsumerKey() string {
	consumerKeyFilePath := getFilePathFromConfigFolder(CONSUMER_KEY_FILE_NAME)
	file, err := os.Open(consumerKeyFilePath)
	if err != nil {
		log.Fatalf("Could not find consumer key file. Please provide consumer key to %s", consumerKeyFilePath)
	}
	defer file.Close()
	consumerKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading consumer key from file")
	}
	return string(consumerKeyBytes)
}

func getConfigSubfolderPath() string {
	homeDir := utils.UserHomeDir()
	subfolder := filepath.Join(homeDir, CONFIG_SUBFOLDER_NAME)
	err := os.MkdirAll(subfolder, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create directory for config dir")
	}
	return subfolder
}

func getFilePathFromConfigFolder(filename string) string {
	subfolder := getConfigSubfolderPath()
	filePath := filepath.Join(subfolder, filename)
	return filePath
}

func loadConfig(filePath string) *config.Config {
	reader, err := os.Open(filePath)
	var configObj *config.Config
	if err != nil {
		fmt.Println("Could not load config file, generating new one")
		configObj = handleConfigFileError(filePath)
	} else {
		defer reader.Close()
		configObj, err = config.ReadConfig(reader)
		if err != nil {
			fmt.Println("Error parsing config file, generating new one")
			configObj = handleConfigFileError(filePath)
		}
	}
	return configObj
}

func handleConfigFileError(filePath string) *config.Config {
	configObj := &config.Config{}
	writer, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating config file, could not open file handle")
		return configObj
	}
	err = config.WriteConfig(configObj, writer)
	if err != nil {
		fmt.Println("Error writing config file: " + err.Error())
	}
	return configObj
}

func initializePocketConnection(consumerKey string) *connector.PocketClient {
	pocketAuth := &PocketAuth{}

	var accessToken string

	authFilePath := getFilePathFromConfigFolder(AUTH_FILE_NAME)

	err := loadFromJSON(authFilePath, &pocketAuth)
	if err != nil {
		fmt.Println("Could not read auth.json file.")
	} else {
		accessToken = pocketAuth.AccessToken
	}

	pocketClient := connector.CreatePocketClient(consumerKey, accessToken)

	if accessToken == "" {
		fmt.Println("Creating new access token")
		accessToken, err = actions.PerformAuth(pocketClient, 3000, utils.OpenBrowser)
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

func loadPocketArticles(pocketClient *connector.PocketClient, articlesList *[]models.ArticleResult) {
	var mut sync.Mutex

	// Count should be high to prevent rate limiting for user
	count := 200

	err := actions.GetUntaggedArticles(pocketClient, 0, count, articlesList, &mut)
	if err != nil {
		log.Fatal(err)
	}

	i := count
	for i < 1000 {
		go func(start int) {
			err := actions.GetUntaggedArticles(pocketClient, start, count, articlesList, &mut)
			if err != nil {
				log.Fatal(err)
			}
		}(i)
		i += count
	}
}

func runArticleLoop(pocketClient *connector.PocketClient, articlesList *[]models.ArticleResult, config *config.Config) {
	articles := *articlesList

	if len(articles) == 0 {
		fmt.Println("No untagged articles")
		return
	}

	clipboardManager := new(utils.ClipboardManagerImpl)

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

		command := readUserInput(func(input string) string {
			return strings.TrimSpace(strings.ToLower(input))
		})

		if command == "" {
			fmt.Println("Please select an option")
		}

		switch command {
		case "t":
			fmt.Println("Enter tags separated by commas (',')")

			tags := readUserInputAsArray(func(input string) []string {
				return strings.Split(input, ",")
			})

			result, err := actions.MarkArticleWithTag(pocketClient, article.ItemID, tags)
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
				_, err := actions.UnfavouriteArticle(pocketClient, article.ItemID)
				if err != nil {
					fmt.Println("Failure unfavouriting article: ", err)
				} else {
					fmt.Println("Success unfavouriting article")
				}
			} else {
				_, err := actions.FavouriteArticle(pocketClient, article.ItemID)
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
			_, err := actions.BumpArticleToTop(pocketClient, article.ItemID)
			if err != nil {
				fmt.Println("Failure bumping article: ", err)
			} else {
				fmt.Println("Success bumping article")
			}
			break
		case "a":
			_, err := actions.ArchiveArticle(pocketClient, article.ItemID)
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
			gotURLFromClipboard := false
			url := ""
			if config.Clipboard && utils.IsURLInClipboard(clipboardManager) {
				var err error
				url, err = utils.GetFromClipboard(clipboardManager)
				if err != nil {
					fmt.Println("There was an error getting URL from clipboard. Continuing to manual add...")
				} else {
					fmt.Printf("'%s' was found in your clipboard. Would you like to add it? [Y/n]\n", url)
					if readYesNoFromUser() {
						gotURLFromClipboard = true
					} else {
						fmt.Println("Did not add URL to Pocket list.")
					}
				}
			}
			if !gotURLFromClipboard {
				fmt.Println("Enter the URL to add:")
				url = readUserInput(func(input string) string {
					return strings.TrimSpace(input)
				})
				if !utils.IsURL(url) {
					fmt.Println("Invalid URL")
					break
				}
			}
			result, err := actions.AddArticle(pocketClient, url)
			if err != nil {
				fmt.Println("Error adding article: ", err)
			} else {
				fmt.Printf("Success adding article: '%s' (%s)\n", result["title"], result["resolved_url"])
			}
			break
		case "d":
			fmt.Println("Are you sure you want to delete this article? You won't be able to restore it unless you readd it. [Y/n]")

			if !readYesNoFromUser() {
				fmt.Println("Article was not deleted")
				break
			}
			fallthrough
		case "dd":
			_, err := actions.DeleteArticle(pocketClient, article.ItemID)
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
			fallthrough
		case "q":
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
	consumerKey := loadConsumerKey()
	configObj := loadConfig(getFilePathFromConfigFolder(CONFIG_FILE_NAME))

	pocketClient := initializePocketConnection(consumerKey)

	articles := make([]models.ArticleResult, 0, 10)

	loadPocketArticles(pocketClient, &articles)

	runArticleLoop(pocketClient, &articles, configObj)
}

func loadFromJSON(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open file '%s'", path)
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func saveToJSON(path string, v interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("Could not open file '%s'", path)
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(v)
}

func readUserInputRaw() string {
	br := bufio.NewReader(os.Stdin)
	input, _ := br.ReadString('\n')

	return input
}

func readUserInput(strFunc func(string) string) string {
	return strFunc(readUserInputRaw())
}

func readUserInputAsArray(strFunc func(string) []string) []string {
	return strFunc(readUserInputRaw())
}

func readYesNoFromUser() bool {
	response := readUserInput(func(input string) string {
		return strings.TrimSpace(strings.ToLower(input))
	})

	return response == "y" || response == "yes"
}
