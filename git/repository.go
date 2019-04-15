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

type Repository interface {
	fmt.Stringer
	Owner() string
	Name() string
	URL() *url.URL
	SHA1() string
	Find(pattern string, entryType EntryType) ([]Entry, error)
}

func NewRepository(owner, repo, sha1, token string) Repository {
	return &repository{owner: owner, repo: repo, sha1: sha1, token: token}
}

type repository struct {
	owner string
	repo  string
	sha1  string
	token string

	ctx          context.Context
	httpClient   *http.Client
	githubClient *github.Client
}

func (r *repository) SHA1() string {
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
	urlStr := fmt.Sprintf("http://github.com/%v/%v", r.owner, r.repo)
	return misc.UrlMustParse(urlStr)
}

func (r *repository) Find(pattern string, entryType EntryType) ([]Entry, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	client := r.GitHubClient()
	tree, _, err := client.Git.GetTree(r.Context(), r.owner, r.repo, r.sha1, true)
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
