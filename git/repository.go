package git

import (
	"context"
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/misc"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"regexp"
)

var repoCache = map[string]map[SHA1]Repository{}

type Repository interface {
	fmt.Stringer
	Owner() string
	Name() string
	URL() *url.URL
	SHA1() SHA1
	Find(pattern string, entryType EntryType) ([]Entry, error)
}

func MakeRepoUrlStr(owner string, repo string) string {
	return fmt.Sprintf("http://github.com/%v/%v", owner, repo)
}

func GetRepository(owner, repoName string, sha1 SHA1, token string) (Repository, error) {
	if owner == "" {
		return nil, fmt.Errorf("repo must have an owner")
	}
	if repoName == "" {
		return nil, fmt.Errorf("repo must have a name")
	}
	if sha1 == "" {
		return nil, fmt.Errorf("repo must have a revision")
	}
	urlStr := MakeRepoUrlStr(owner, repoName)

	var ok bool
	var reposBySHA1 map[SHA1]Repository
	if reposBySHA1, ok = repoCache[urlStr]; !ok {
		reposBySHA1 = map[SHA1]Repository{}
		repoCache[urlStr] = reposBySHA1
	}
	var repo Repository
	if repo, ok = reposBySHA1[sha1]; !ok {
		repo = &repository{owner: owner, repo: repoName, sha1: sha1, token: token}
		reposBySHA1[sha1] = repo
	}
	return repo, nil
}

type repository struct {
	owner string
	repo  string
	sha1  SHA1
	token string

	ctx          context.Context
	httpClient   *http.Client
	githubClient *github.Client
}

func (r *repository) SHA1() SHA1 {
	return r.sha1
}

func (r *repository) String() string {
	return fmt.Sprintf("%v/%v", r.owner, r.repo)
}

func (r *repository) Owner() string {
	return r.owner
}

func (r *repository) Name() string {
	return r.repo
}

func (r *repository) URL() *url.URL {
	urlStr := MakeRepoUrlStr(r.owner, r.repo)
	return misc.UrlMustParse(urlStr)
}

func (r *repository) Find(pattern string, entryType EntryType) ([]Entry, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	client := r.GitHubClient()
	tree, _, err := client.Git.GetTree(r.Context(), r.owner, r.repo, r.sha1.Full(), true)
	if err != nil {
		return nil, err
	}
	if tree.GetTruncated() {
		return nil, fmt.Errorf("repository %v/%v has too many files to return as a flat list", r.owner, r.repo)
	}

	var entries []Entry
	for _, e := range tree.Entries {
		eType := GetEntryType(e)
		if eType != entryType {
			continue
		}
		path := e.GetPath()
		if !re.MatchString(path) {
			continue
		}
		u, err := url.Parse(e.GetURL())
		if err != nil {
			return entries, err
		}
		entries = append(entries, r.NewEntry(path, e.GetSHA(), eType, e.GetSize(), u))
	}

	return entries, nil
}

func (r *repository) Context() context.Context {
	if r.ctx == nil {
		r.ctx = context.Background()
	}
	return r.ctx
}

func (r *repository) HttpClient() *http.Client {
	if r.httpClient == nil {
		if r.token == "" {
			r.httpClient = http.DefaultClient
		} else {
			token := oauth2.Token{AccessToken: r.token}
			tokenSource := oauth2.StaticTokenSource(&token)
			r.httpClient = oauth2.NewClient(r.Context(), tokenSource)
		}
	}
	return r.httpClient
}

func (r *repository) GitHubClient() *github.Client {
	if r.githubClient == nil {
		httpClient := r.HttpClient()
		r.githubClient = github.NewClient(httpClient)
	}
	return r.githubClient
}
