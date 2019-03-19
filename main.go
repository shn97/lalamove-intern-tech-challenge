package main

import (
	"context"
	"fmt"
	"bufio"
	"strings"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

type RepoInfo struct {
	owner string
	repo string
	minVersion string
}

func SortByDesc(releases []*semver.Version) {
	semver.Sort(releases)

	for i, j := 0, len(releases)-1; i < j; i, j = i+1, j-1 {
		releases[i], releases[j] = releases[j], releases[i]
	}
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
		
	SortByDesc(releases)
	
	currMajor := int64(-1)
	currMinor := int64(-1)
	for _, release := range releases {
		if release.Compare(*minVersion) > 0 &&
			(release.Major != currMajor || release.Minor != currMinor) {
			
			versionSlice = append(versionSlice, release)
			currMajor = release.Major
			currMinor = release.Minor
		}
	}
	
	return versionSlice
}

func ParseFile(fileName string) []*RepoInfo {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)	
	scanner.Scan()
	
	var repoInfos []*RepoInfo
	for scanner.Scan() {
		col := strings.Split(scanner.Text(), ",")
		repositoryInput := col[0]
		repository := strings.Split(repositoryInput, "/")
		owner := repository[0]
		repo :=repository[1]	
		minVersion := col[1]
		
		info :=&RepoInfo{owner, repo, minVersion}
		repoInfos = append(repoInfos, info)	
	}
	return repoInfos
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	fileNameArg := os.Args[1]
	repoInfos := ParseFile(fileNameArg)
	
	// GitHub
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	
	for _, repoInfo := range repoInfos {
		releases, _, err := client.Repositories.ListReleases(ctx, repoInfo.owner, repoInfo.repo, opt)
		if err != nil {
			panic(err) // is this really a good way?
		}
		minVersion := semver.New(repoInfo.minVersion)
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)
	
		fmt.Printf("latest versions of %s/%s: %s\n", repoInfo.owner, repoInfo.repo, versionSlice)
	}
}
