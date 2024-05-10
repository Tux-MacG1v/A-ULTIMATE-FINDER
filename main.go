package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func importModule(moduleName string) error {
	cmd := exec.Command("go", "get", moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadGitHubRelease(owner, repo string) error {
	tagNameURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/version.db", owner, repo)
	resp, err := http.Get(tagNameURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tagName := ""
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		tagName = strings.TrimSpace(string(bodyBytes))
	}

	if tagName == "" {
		return fmt.Errorf("failed to get tag name")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tagName)
	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get release info")
	}

	var releaseInfo struct {
		Assets []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, asset := range releaseInfo.Assets {
		assetURL := asset.BrowserDownloadURL
		assetName := asset.Name
		downloadPath := filepath.Join(currentDir, assetName)

		fmt.Println("Downloading ULTIMATE-FINDER...")
		resp, err := http.Get(assetURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(downloadPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}
	}

	fmt.Println("LATEST TOOL DOWNLOAD COMPLETE...")
	fmt.Println()

	executablePath := filepath.Join(currentDir, releaseInfo.Assets[0].Name)
	fmt.Println("TOOL IS RUNNING...")
	cmd := exec.Command(executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	owner := "Tux-MacG1v"
	repo := "A-ULTIMATE-FINDER"

	fmt.Println("SEARCHING FOR LATEST...")
	if err := downloadGitHubRelease(owner, repo); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
