package linkmodel

type LinkStatus string

const (
	LinkStatusAvailable    LinkStatus = "available"
	LinkStatusNotAvailable LinkStatus = "not available"
)

type CheckLinksRequest struct {
	Links []string `json:"links"`
}

type CheckLinkResponse struct {
	Links    map[string]LinkStatus `json:"links"` // key - link, value - statis
	LinksNum int                   `json:"links_num"`
}

type CheckLinkListRequest struct {
	LinksList []int `json:"links_list"`
}
