package getignore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const GitignoreRepoContentsURL = "https://api.github.com/repos/GitHub/gitignore/contents"

type GitignoreService struct {
	client *http.Client
}

type GetIgnoresResponse struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

func (g *GitignoreService) getRawIgnores() ([]GetIgnoresResponse, error) {

	if g.client == nil {
		g.client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := g.client.Get(GitignoreRepoContentsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var contents []GetIgnoresResponse

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&contents)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return contents, nil
}

func (g *GitignoreService) GetIgnores() ([]string, error) {

	ignores, err := g.getRawIgnores()
	if err != nil {
		return nil, err
	}

	ignoreNames := make([]string, 0, len(ignores))

	for _, ignore := range ignores {
		if ignore.Type != "file" {
			continue
		}

		name, found := strings.CutSuffix(ignore.Name, ".gitignore")
		if !found { // doesn't end with .gitignore
			continue
		}

		ignoreNames = append(ignoreNames, name)
	}

	return ignoreNames, nil
}

func (g *GitignoreService) GetGitignoreContents(lang string) (*string, error) {
	availableIgnoreFiles, err := g.getRawIgnores()
	if err != nil {
		return nil, err
	}

	var requestedGitignore *GetIgnoresResponse
	for _, ignore := range availableIgnoreFiles {
		ignoreName, found := strings.CutSuffix(ignore.Name, ".gitignore")
		if !found {
			continue
		}

		if strings.ToLower(ignoreName) == strings.ToLower(lang) {
			requestedGitignore = &ignore
		}
	}

	if requestedGitignore == nil {
		return nil, fmt.Errorf("no gitignore found for language %s, see --list", lang)
	}

	if g.client == nil {
		g.client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := g.client.Get(requestedGitignore.DownloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download gitignore: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code when downloading: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	content := string(body)
	return &content, nil
}
