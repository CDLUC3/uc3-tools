package git

import (
	"context"
	"fmt"
	"github.com/CDLUC3/uc3-tools/uc3-build-info/shared"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var repoCache = map[string]map[SHA1]Repository{}

type Repository interface {
	fmt.Stringer
	Owner() string
	Name() string
	URL() *url.URL
	SHA1() SHA1
	Find(pattern string, entryType EntryType) (entries []Entry, errors []error)

	// Unexported symbols
	GetEntry(path string, sha1 SHA1, eType EntryType, size int, u *url.URL) Entry
}

func MakeRepoUrlStr(owner string, repo string) string {
	return fmt.Sprintf("http://github.com/%v/%v", owner, repo)
}

func GetRepository(owner, repoName string, sha1 SHA1) (Repository, error) {
	if owner == "" {
		return nil, fmt.Errorf("repo must have an owner")
	}
	if repoName == "" {
		return nil, fmt.Errorf("repo must have a name")
	}
	if sha1 == "" {
		return nil, fmt.Errorf("repo must have a revision")
	}
	urlStrLowerCase := strings.ToLower(MakeRepoUrlStr(owner, repoName))

	var ok bool
	var reposBySHA1 map[SHA1]Repository
	if reposBySHA1, ok = repoCache[urlStrLowerCase]; !ok {
		reposBySHA1 = map[SHA1]Repository{}
		repoCache[urlStrLowerCase] = reposBySHA1
	}
	var repo Repository
	if repo, ok = reposBySHA1[sha1]; !ok {
		repo = &repository{owner: owner, repo: repoName, sha1: sha1}
		reposBySHA1[sha1] = repo
	}
	return repo, nil
}

type repository struct {
	owner   string
	repo    string
	sha1    SHA1
	entries map[string]map[SHA1]Entry

	ctx          context.Context
	httpClient   *http.Client
	githubClient *github.Client

	tree *github.Tree
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
	return shared.UrlMustParse(urlStr)
}

func (r *repository) GetEntry(path string, sha1 SHA1, eType EntryType, size int, u *url.URL) Entry {
	if r.entries == nil {
		r.entries = map[string]map[SHA1]Entry{}
	}
	var ok bool
	var entriesBySHA1 map[SHA1]Entry
	if entriesBySHA1, ok = r.entries[path]; !ok {
		entriesBySHA1 = map[SHA1]Entry{}
		r.entries[path] = entriesBySHA1
	}
	var e Entry
	if e, ok = entriesBySHA1[sha1]; !ok {
		e = &entry{path: path, sha1: sha1, eType: eType, size: size, url: u, repository: r}
		entriesBySHA1[sha1] = e
	}
	return e
}

func (r *repository) Tree() (*github.Tree, error) {
	if r.tree == nil {
		client := r.GitHubClient()
		tree, _, err := client.Git.GetTree(r.Context(), r.owner, r.repo, r.sha1.Full(), true)
		if err != nil {
			return nil, err
		}
		r.tree = tree
	}
	return r.tree, nil
}

func (r *repository) Find(pattern string, entryType EntryType) ([]Entry, []error) {

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, []error{err}
	}

	tree, err := r.Tree()
	if err != nil {
		return nil, []error{err}
	}
	if tree.GetTruncated() {
		err := fmt.Errorf("repository %v/%v has too many files to return as a flat list", r.owner, r.repo)
		return nil, []error{err}
	}

	var entries []Entry
	var errors []error
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
			errors = append(errors, err)
		} else {
			entry := r.GetEntry(path, SHA1(e.GetSHA()), eType, e.GetSize(), u)
			if !r.isThisRepo(entry.URL()) {
				err := fmt.Errorf("entry URL %v does not appear to belong to repository %v/%v (repo moved?)", entry.URL(), r.owner, r.repo)
				errors = append(errors, err)
			}
			entries = append(entries, entry)
		}
	}

	return entries, errors
}

func (r *repository) Context() context.Context {
	if r.ctx == nil {
		r.ctx = context.Background()
	}
	return r.ctx
}

func (r *repository) HttpClient() *http.Client {
	if r.httpClient == nil {
		token := strings.TrimSpace(Token)
		if token == "" {
			_, _ = fmt.Fprintln(os.Stderr, tokenNotProvided)
			os.Exit(-1)
		} else {
			oauth2Token := oauth2.Token{AccessToken: token}
			tokenSource := oauth2.StaticTokenSource(&oauth2Token)
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

func (r *repository) isThisRepo(entryUrl *url.URL) bool {
	urlPath := strings.ToLower(entryUrl.Path)
	prefix := strings.ToLower(fmt.Sprintf("/repos/%v/%v/", r.owner, r.repo))
	if !strings.HasPrefix(urlPath, prefix) {
		return false
	}
	return true
}

