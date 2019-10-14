package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/config"
	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/models"
	"github.com/Coteh/gyroid/lib/utils"
)

// configSubfolderName is the subfolder containing config files
const configSubfolderName = ".config/gyroid"

// configFileName is the name of the config file
const configFileName = "config.yml"

// authFileName is the name of the auth file
const authFileName = "auth.json"

// consumerKeyFileName is the name of the consumer key file
const consumerKeyFileName = "consumer_key"

// PocketAuth contain the access token TODO remove
type PocketAuth struct {
	AccessToken string `json:"access_token"`
}

func loadConsumerKey() string {
	consumerKeyFilePath := getFilePathFromConfigFolder(consumerKeyFileName)
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
	subfolder := filepath.Join(homeDir, configSubfolderName)
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

	authFilePath := getFilePathFromConfigFolder(authFileName)

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

func loadPocketArticles(pocketClient *connector.PocketClient, articlesList *[]models.ArticleResult, mut *sync.Mutex) {
	// Count should be high to prevent rate limiting for user
	count := 200

	err := actions.GetUntaggedArticles(pocketClient, 0, count, articlesList, mut)
	if err != nil {
		log.Fatal(err)
	}

	i := count
	// TODO what if user has more than 1000 articles? :O
	// TODO fix shifting problem that can occur when user adds article before goroutine gets untagged articles (#40)
	for i < 1000 {
		go func(start int) {
			err := actions.GetUntaggedArticles(pocketClient, start, count, articlesList, mut)
			if err != nil {
				log.Fatal(err)
			}
		}(i)
		i += count
	}
}

func addArticle(pocketClient *connector.PocketClient, config *config.Config, clipboardManager utils.ClipboardManager) (addedArticle *models.AddedArticleResult, err error) {
	gotURLFromClipboard := false
	url := ""
	if config.Clipboard && utils.IsURLInClipboard(clipboardManager) && !utils.CheckURLMostRecentlyAdded(clipboardManager) {
		var err error
		url, err = utils.GetFromClipboard(clipboardManager)
		if err != nil {
			fmt.Println("There was an error getting URL from clipboard. Continuing to manual add...")
		} else {
			fmt.Printf("'%s' was found in your clipboard. Would you like to add it? [Y/n]\n", url)
			if readYesNoFromUser() {
				clipboardManager.SetMostRecentlyAddedURL(url)
				gotURLFromClipboard = true
			} else {
				fmt.Println("Did not add URL from clipboard to Pocket list.")
			}
		}
	}
	if !gotURLFromClipboard {
		fmt.Println("Enter the URL to add:")
		url = readUserInput(func(input string) string {
			return strings.TrimSpace(input)
		})
		if url == "" {
			return nil, errors.New("Did not enter a URL")
		}
		if !utils.IsURL(url) {
			return nil, errors.New("Invalid URL")
		}
	}
	result, err := actions.AddArticle(pocketClient, url)
	if err != nil {
		return nil, err
	} else {
		result.ArticleResult.ResolvedTitle = result.Title
		return result, nil
	}
}

func printArticle(article models.ArticleResult, minutes int) {
	var minutesStr string
	if minutes > 0 {
		var plural rune
		if minutes != 1 {
			plural = 's'
		}
		minutesStr = fmt.Sprintf("Expected Read Time: %d minute%c\n", minutes, plural)
	}
	fmt.Printf("-----\n'%s'\n'%s'\n%s%s\n-----\n", article.ResolvedTitle, article.Excerpt, minutesStr, article.ResolvedURL)
}

func printArticleActions(isFav bool, config *config.Config) {
	unfavStr := ""
	if isFav {
		unfavStr = "un"
	}

	fmt.Printf("[T]ag\t%s[F]avourite\t[B]ump\t[A]rchive\t[D]elete\t[DD]elete with yes\t[O]pen",
		unfavStr)
	if config.Clipboard {
		fmt.Printf("\t[C]opy URL")
	}
	fmt.Printf("\t[+]Add Article by URL\t[N]ext\t[E]xit\n-----\n")
}

func printSuccessfullyAddedArticle(article *models.AddedArticleResult) {
	fmt.Printf("Success adding article: '%s' (%s)\n", article.Title, article.ResolvedURL)
}

func runArticleLoop(pocketClient *connector.PocketClient, articlesList *[]models.ArticleResult, config *config.Config, clipboardManager utils.ClipboardManager, mut *sync.Mutex) {
	i := 0
	isFav := false
	isNext := false
	userMarkedFav := false
	refiningNew := false
	var article models.ArticleResult
	for i < len(*articlesList) {
		if !refiningNew {
			article = (*articlesList)[i]
		}

		if !userMarkedFav && article.Favorite == 1 {
			isFav = true
		}

		minutes := utils.CalculateExpectedReadTime(article.WordCount)
		printArticle(article, minutes)
		printArticleActions(isFav, config)

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
		case "b":
			_, err := actions.BumpArticleToTop(pocketClient, article.ItemID)
			if err != nil {
				fmt.Println("Failure bumping article: ", err)
			} else {
				fmt.Println("Success bumping article")
			}
		case "a":
			_, err := actions.ArchiveArticle(pocketClient, article.ItemID)
			if err != nil {
				fmt.Println("Failure archiving article: ", err)
			} else {
				fmt.Println("Success archiving article")
			}
			isNext = true
		case "o":
			utils.OpenBrowser(article.ResolvedURL)
		case "+":
			result, err := addArticle(pocketClient, config, clipboardManager)
			if err != nil {
				fmt.Printf("There was an error adding article: %s\n", err)
				break
			}
			printSuccessfullyAddedArticle(result)
			if config.RefineNew == "prompt" {
				fmt.Printf("Would you like to refine it now? [Y/n]")
			}
			refiningNew = config.RefineNew == "yes" || config.RefineNew == "prompt" && readYesNoFromUser()
			if !refiningNew {
				mut.Lock()
				*articlesList = append(*articlesList, *result.ArticleResult)
				mut.Unlock()
			} else {
				article = *result.ArticleResult
				isNext = false
				isFav = false
				userMarkedFav = false
			}
		case "c":
			if !config.Clipboard {
				fmt.Println("Clipboard not enabled")
				break
			}
			err := clipboardManager.CopyToClipboard(article.ResolvedURL)
			if err != nil {
				fmt.Println("Failure copying article URL: ", err)
			} else {
				fmt.Println("Success copying article URL")
			}
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
		case "e":
			fallthrough
		case "q":
			return
		}
		if isNext {
			// If user added article then said yes they would like to refine it right away,
			// then incrementing i is blocked when saying "next" so that they are brought back
			// to article they were looking at before
			if !refiningNew {
				i++
			}
			isFav = false
			isNext = false
			userMarkedFav = false
			refiningNew = false
		}
	}
	if len(*articlesList) == 0 {
		fmt.Println("You have no articles in your list")
	} else {
		fmt.Println("End of Pocket list reached")
	}
	fmt.Println("Would you like to add more? [Y/n]")
	if readYesNoFromUser() {
		result, err := addArticle(pocketClient, config, clipboardManager)
		if err != nil {
			fmt.Printf("There was an error adding article: %s\n", err)
		} else {
			newArticleList := []models.ArticleResult{*result.ArticleResult}
			printSuccessfullyAddedArticle(result)
			runArticleLoop(pocketClient, &newArticleList, config, clipboardManager, mut)
		}
	} else {
		fmt.Println("Bye")
	}
}

func main() {
	consumerKey := loadConsumerKey()
	configObj := loadConfig(getFilePathFromConfigFolder(configFileName))
	clipboardManager := new(utils.ClipboardManagerImpl)

	pocketClient := initializePocketConnection(consumerKey)

	articles := make([]models.ArticleResult, 0, 10)

	var mut sync.Mutex

	loadPocketArticles(pocketClient, &articles, &mut)

	runArticleLoop(pocketClient, &articles, configObj, clipboardManager, &mut)
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
