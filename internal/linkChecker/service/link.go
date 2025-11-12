package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/Negat1v9/link-checker/internal/linkChecker/linkstore"
	linkmodel "github.com/Negat1v9/link-checker/internal/linkChecker/model"
)

type LinkCheckerService struct {
	httpClient *http.Client

	// store for links group
	linkStore *linkstore.LinkStore
}

func NewLinkService(linkStore *linkstore.LinkStore) *LinkCheckerService {
	return &LinkCheckerService{
		httpClient: http.DefaultClient,
		linkStore:  linkStore,
	}
}

// links is array of strings like google.com, ya.ru, apple.com
func (s *LinkCheckerService) CheckLinks(ctx context.Context, links []string) *linkmodel.CheckLinkResponse {
	linksResult := make(map[string]linkmodel.LinkStatus, len(links))
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			isAvailable := s.checkUrl(ctx, link)
			mu.Lock()
			if isAvailable {
				linksResult[link] = linkmodel.LinkStatusAvailable
			} else {
				linksResult[link] = linkmodel.LinkStatusNotAvailable
			}
			mu.Unlock()
			wg.Done()
		}(link)
	}

	wg.Wait()

	return &linkmodel.CheckLinkResponse{
		Links:    linksResult,
		LinksNum: s.linkStore.CreateLinksGroup(links), // save link group
	}

}

func (s *LinkCheckerService) CreateLinksGroupPdfReport(ctx context.Context, groups []int) (PdfFileWriter, error) {
	links := make([]string, 0)
	// recieve all links from all groups
	for _, groupID := range groups {
		links = append(links, s.linkStore.GetLinksByGroup(groupID)...)
	}

	linksResult := make(map[string]linkmodel.LinkStatus, len(links))
	wg := sync.WaitGroup{}
	mu := sync.RWMutex{}
	if len(links) == 0 {
		return nil, fmt.Errorf("not found links group")
	}
	pdfFile := s.createPDFReport()
	for i, link := range links {
		wg.Add(1)
		go func(link string) {
			// check if has dublicate in links array
			mu.RLock()
			if _, ok := linksResult[link]; ok {
				// dublicate -> stop check status
				mu.RUnlock()
				return
			}
			mu.RUnlock()
			// recieve state of link
			isAvailable := s.checkUrl(ctx, link)
			mu.Lock()
			// set status
			if isAvailable {
				linksResult[link] = linkmodel.LinkStatusAvailable
			} else {
				linksResult[link] = linkmodel.LinkStatusNotAvailable
			}
			// writer in report row with link state with mutex lock
			s.addReportRowLinkState(pdfFile, i+1, link, linksResult[link])
			mu.Unlock()
			wg.Done()
		}(link)
	}

	// wait get all links state
	wg.Wait()

	return pdfFile, nil
}

// return true only if request on url return status < 400
func (s *LinkCheckerService) checkUrl(ctx context.Context, link string) bool {
	url, err := url.Parse(link)
	if err != nil {
		return false
	}
	// set schema if was not set
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	// url is not available by default if NewRequestWithContext returned error
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return false
	}
	// url is not available by default if request returned error
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false
	}

	// url is not available
	if resp.StatusCode >= 400 {
		return false
	}

	return true
}
