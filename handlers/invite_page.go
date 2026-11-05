package handlers

import (
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Trang thiệp có tên khách mời (index.html?invite-name=...).
//
// Zalo / Messenger / Facebook không chạy JavaScript khi tạo khung xem trước link,
// chúng chỉ đọc thẻ <meta property="og:..."> trong HTML. Vì vậy nginx chuyển các
// request index.html có tham số invite-name sang đây; handler đọc file index.html
// trên đĩa, sửa og:title / og:url theo tên khách rồi trả về.
// Không có tên khách thì trả nguyên file như cũ.
//
// Cấu hình bằng biến môi trường:
//   INVITE_INDEX_PATH = đường dẫn tuyệt đối tới index.html của thiệp trên server
//   INVITE_SITE_URL   = địa chỉ trang thiệp, mặc định https://cuoi-hoa-tuyet.htsolution.top/index.html

const (
	invitePageParam   = "invite-name"
	inviteNameMaxLen  = 60
	defaultInviteSite = "https://cuoi-hoa-tuyet.htsolution.top/index.html"
)

var (
	ogTitleRe = regexp.MustCompile(`(<meta\s+property="og:title"\s+content=")[^"]*(")`)
	ogURLRe   = regexp.MustCompile(`(<meta\s+property="og:url"\s+content=")[^"]*(")`)
	spacesRe  = regexp.MustCompile(`\s+`)

	indexCache struct {
		sync.Mutex
		body    string
		modTime time.Time
	}
)

// readInviteIndex đọc index.html, chỉ đọc lại khi file trên đĩa thay đổi.
func readInviteIndex() (string, error) {
	path := os.Getenv("INVITE_INDEX_PATH")
	if path == "" {
		return "", os.ErrNotExist
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	indexCache.Lock()
	defer indexCache.Unlock()
	if indexCache.body != "" && info.ModTime().Equal(indexCache.modTime) {
		return indexCache.body, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	indexCache.body = string(data)
	indexCache.modTime = info.ModTime()
	return indexCache.body, nil
}

func cleanInviteName(raw string) string {
	name := strings.TrimSpace(spacesRe.ReplaceAllString(raw, " "))
	if r := []rune(name); len(r) > inviteNameMaxLen {
		name = string(r[:inviteNameMaxLen])
	}
	return name
}

// GetInvitePage trả về index.html với og:title / og:url riêng cho từng khách mời.
func GetInvitePage(c *gin.Context) {
	page, err := readInviteIndex()
	if err != nil {
		// nginx sẽ tự trả file index.html tĩnh khi gặp lỗi này
		log.Printf("invite-page: không đọc được INVITE_INDEX_PATH=%q: %v", os.Getenv("INVITE_INDEX_PATH"), err)
		c.String(http.StatusServiceUnavailable, "invite page not configured: check INVITE_INDEX_PATH")
		return
	}

	if name := cleanInviteName(c.Query(invitePageParam)); name != "" {
		title := "Thân mời " + name + " tham dự đám cưới Trần Hòa & Ánh Tuyết"
		// $1/$2 trong ReplaceAllString có thể bị hiểu nhầm nếu tên chứa "$",
		// nên dùng ReplaceAllStringFunc và tự escape HTML.
		page = replaceMetaContent(page, ogTitleRe, html.EscapeString(title))

		// Facebook/Messenger lấy lại dữ liệu theo og:url, nên og:url phải là link
		// có tên khách, nếu không sẽ hiện tiêu đề chung.
		site := os.Getenv("INVITE_SITE_URL")
		if site == "" {
			site = defaultInviteSite
		}
		pageURL := site + "?" + url.Values{invitePageParam: {name}}.Encode()
		page = replaceMetaContent(page, ogURLRe, html.EscapeString(pageURL))
	}

	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
}

func replaceMetaContent(page string, re *regexp.Regexp, value string) string {
	return re.ReplaceAllStringFunc(page, func(tag string) string {
		m := re.FindStringSubmatch(tag)
		return m[1] + value + m[2]
	})
}
