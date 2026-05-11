package pagination

import (
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PaginationDTO is a reusable, generic pagination response compatible with
// enterprise APIs (mirrors / modern REST pagination structure).
type PaginationDTO[T any] struct {
    CurrentPage  int     `json:"current_page"`
    Data         []T     `json:"data"`
    FirstPageURL *string `json:"first_page_url,omitempty"`
    From         int     `json:"from,omitempty"`
    LastPage     int     `json:"last_page"`
    LastPageURL  *string `json:"last_page_url,omitempty"`
    NextPageURL  *string `json:"next_page_url,omitempty"`
    Path         string  `json:"path"`
    PerPage      int     `json:"per_page"`
    PrevPageURL  *string `json:"prev_page_url,omitempty"`
    To           int     `json:"to,omitempty"`
    Total        int64   `json:"total"`
}

// BuildPaginationFromContext builds a PaginationDTO using request context and
// standard page/limit/total values. It preserves other query parameters (e.g., search, sort)
// and generates first/last/next/prev URLs. Uses BASE_URL env var if present, with a sensible default.
func BuildPaginationFromContext[T any](c *gin.Context, items []T, total int64, page, limit int) PaginationDTO[T] {
    if page <= 0 {
        page = 1
    }
    if limit <= 0 {
        limit = 20
    }

    // -------Base URL from env, fallback to default----
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "https://erp.vidatech.com.bd"
    }
    baseURL = strings.TrimRight(baseURL, "/")

    // Path and original query params
    path := c.Request.URL.Path
    rawQuery := c.Request.URL.RawQuery

    // Parse existing query and keep all params except page & limit/page_size
    q, _ := url.ParseQuery(rawQuery)
    // detect which page-size param the client used so we preserve it in generated URLs
    pageSizeParam := "limit"
    if _, ok := q["page_size"]; ok {
        pageSizeParam = "page_size"
    } else if _, ok := q["limit"]; ok {
        pageSizeParam = "limit"
    }
    // remove pagination keys from the base query so we don't duplicate them
    q.Del("page")
    q.Del("limit")
    q.Del("page_size")

    // Utility to build url string with replaced page/limit
    buildURL := func(p int) string {
        // copy query
        qq := url.Values{}
        for k, vals := range q {
            for _, v := range vals {
                qq.Add(k, v)
            }
        }
    qq.Set("page", strconv.Itoa(p))
    qq.Set(pageSizeParam, strconv.Itoa(limit))
        u := baseURL + path
        if encoded := qq.Encode(); encoded != "" {
            u = u + "?" + encoded
        }
        return u
    }

    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    if totalPages < 1 {
        totalPages = 1
    }

    // -------- compute from/to------------
    var from int
    if total == 0 {
        from = 0
    } else {
        from = (page-1)*limit + 1
    }
    to := page * limit
    if int64(to) > total {
        to = int(total)
    }

    var firstURL, lastURL, nextURL, prevURL *string

    //------------ Always provide first and last--------
    f := buildURL(1)
    firstURL = &f
    l := buildURL(totalPages)
    lastURL = &l

    if page < totalPages {
        n := buildURL(page + 1)
        nextURL = &n
    }
    if page > 1 && page <= totalPages {
        p := buildURL(page - 1)
        prevURL = &p
    }

    dto := PaginationDTO[T]{
        CurrentPage:  page,
        Data:         items,
        FirstPageURL: firstURL,
        From:         from,
        LastPage:     totalPages,
        LastPageURL:  lastURL,
        NextPageURL:  nextURL,
        Path:         baseURL + path,
        PerPage:      limit,
        PrevPageURL:  prevURL,
        To:           to,
        Total:        total,
    }

    return dto
}

// BuildPaginationMap builds a pagination response as a map[string]interface{}
// This is useful for non-generic callers (other packages) that have concrete
// slices and cannot easily use generics. The "data" field will contain the
// provided data value as-is. All URL fields (first, last, next, prev) are always
// included in the map; prev/next will be nil when not applicable.
func BuildPaginationMap(c *gin.Context, data interface{}, total int64, page, limit int) map[string]interface{} {
    dto := BuildPaginationFromContext(c, []interface{}{}, total, page, limit)

    // Convert dto to generic map but keep provided data in `data` field
    // Always include all fields (including nil URLs) for consistent API contract
    m := map[string]interface{}{
        "current_page":   dto.CurrentPage,
        "data":           data,
        "first_page_url": dto.FirstPageURL,
        "from":           dto.From,
        "last_page":      dto.LastPage,
        "last_page_url":  dto.LastPageURL,
        "next_page_url":  dto.NextPageURL,
        "path":           dto.Path,
        "per_page":       dto.PerPage,
        "prev_page_url":  dto.PrevPageURL,
        "to":             dto.To,
        "total":          dto.Total,
    }

    return m
}

