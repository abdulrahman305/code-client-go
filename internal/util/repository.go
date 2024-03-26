/*
 * © 2024 Snyk Limited All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/go-git/go-git/v5"
)

func GetRepositoryUrl(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		fmt.Errorf("open local repository: %w", err)
		return "", err
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		fmt.Errorf("get remote: %w", err)
		return "", err
	}

	if len(remote.Config().URLs) == 0 {
		return "", fmt.Errorf("no repository urls available")
	}

	repoUrl := remote.Config().URLs[0]

	if hasCredentials(repoUrl) {
		repoUrl, err = sanatiseCredentials(repoUrl)
	}

	return repoUrl, nil
}

func hasCredentials(rawUrl string) bool {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return false // Failed to parse URL
	}

	if parsedURL.User != nil {
		_, hasPassword := parsedURL.User.Password()
		return hasPassword
	}

	return false // No user info in URL
}

func sanatiseCredentials(url string) (string, error) {
	re, err := regexp.Compile(`(?<=://)[^@]+`)

	if err != nil {
		//fmt.Errorf("Error compiling regex: %w", err)
		return "", err
	}

	strippedURL := re.ReplaceAllString(url, "")
	return strippedURL, nil
}
