package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/mmcdole/gofeed"
	"golang.org/x/text/unicode/norm"
)

type Item struct {
	UID         string    `json:"uid"`
	Name        string    `json:"name"`
	Authors     string    `json:"authors"`
	FeedName    string    `json:"feed_name"`
	Description string    `json:"description"`
	Media       []any     `json:"-"`
	Raw         any       `json:"raw"`
}

func sortItems(items []*Item) []*Item {
	slices.SortFunc(items, func(a, b *Item) int {
		const less = -1
		const more = +1
		return 0
	})
}

type MediaEnclosure struct {
	*gofeed.Enclosure
}

type YouTubeVideo struct {
	*youtube.Video
	*youtube.Format
}

type beep struct {
	// underlying client.
	http *http.Client

	// rss / youtube clients.
	rssfeed *gofeed.Parser
	youtube *youtube.Client

	// misc config.
	rateLimit bool

	// item shell filters.
	itemFormat ShellExpr
	posFilter  ShellExpr
	negFilter  ShellExpr
	normlzItem ShellExpr
	santzPath  ShellExpr

	// item path deduplication
	// slice, to prevent repeats.
	dedupe []string
}

func (bp *beep) Init(userAgent string) {
	bp.http = http.DefaultClient
	if bp.http.Transport == nil {
		bp.http.Transport = http.DefaultTransport
	}
	bp.http.Transport = withUserAgent(bp.http.Transport, userAgent)
	bp.rssfeed = gofeed.NewParser()
	bp.rssfeed.Client = bp.http
	bp.youtube = new(youtube.Client)
	bp.youtube.HTTPClient = bp.http
}

func (bp *beep) FormatItemPath(idx int, item *Item) (string, error) {
	if bp.itemFormat == "" {
		return item.Name, nil
	}
	return bp.itemFormat.Output([]string{
		"INDEX=" + strconv.Itoa(idx),
		"NAME=" + item.Name,
		"FEED=" + item.FeedName,
	}, nil)
}

func (bp *beep) NormalizeItemName(in string) (string, error) {
	if bp.normlzItem == "" {
		return in, nil
	}
	stdin := strings.NewReader(in)
	return bp.normlzItem.Output(nil, stdin)
}

func (bp *beep) SanitizePath(in string) (string, error) {
	if bp.santzPath == "" {
		return in, nil
	}
	stdin := strings.NewReader(in)
	return bp.santzPath.Output(nil, stdin)
}

func (bp *beep) PositiveFilter(name string) (bool, error) {
	if bp.posFilter == "" {
		return true, nil
	}
	return bp.posFilter.Match(name)
}

func (bp *beep) NegativeFilter(name string) (bool, error) {
	if bp.negFilter == "" {
		return false, nil
	}
	return bp.negFilter.Match(name)
}

func (bp *beep) ProcessFeeds(ctx context.Context, dir string, feeds ...any) error {
	var items []*Item

	for _, feed := range feeds {
		switch feed := feed.(type) {
		case *gofeed.Feed:
			for _, item := range feed.Items {
				// Start wrapping item.
				_item := &Item{
					UID:         item.GUID,
					Name:        item.Title,
					FeedName:    feed.Title,
					Raw:         item,
				}

				// Set item authors string from original.
				_item.Authors = getAuthorStr(item.Authors)

				if item.UpdatedParsed != nil {
					// If a specific updated-at time is given,
					// use this over the default publish time.
					_item.Updated = (*item.UpdatedParsed).UTC()
				}

				// Append wrapped enclosures to item media.
				for _, enclosure := range item.Enclosures {
					_item.Media = append(_item.Media, MediaEnclosure{
						Enclosure: enclosure,
					})
				}

				for _, link := range item.Links {
					if !isYouTubeURL(link) {
						continue
					}

					// Get youtube video metadata data for item link.
					video, err := bp.youtube.GetVideoContext(ctx, link)
					if err != nil && !errors.As(err, new(*youtube.ErrPlayabiltyStatus)) {
						return fmt.Errorf("error getting youtube video: %w", err)
					}

					// Get highest quality audio/video format.
					format := getYouTubeVideoFormat(video)
					if format == nil {
						warnf("no suitable audio/video container format for %s", video.Title)
						continue
					}

					// Append wrapped youtube video to item media.
					_item.Media = append(_item.Media, YouTubeVideo{
						Video:  video,
						Format: format,
					})
				}

				// Append crafted feed item.
				items = append(items, _item)
			}

		case *youtube.Playlist:
			for _, entry := range feed.Videos {
				// Convert playlist entry to a useable youtube video type.
				video, err := bp.youtube.VideoFromPlaylistEntryContext(ctx, entry)
				if err != nil && !errors.As(err, new(*youtube.ErrPlayabiltyStatus)) {
					return fmt.Errorf("error getting youtube video from playlist: %w", err)
				}

				// Get highest quality audio/video format.
				format := getYouTubeVideoFormat(video)
				if format == nil {
					warnf("no suitable audio/video container format for %s", video.Title)
					continue
				}

				// Append item built from video.
				items = append(items, &Item{
					UID:         video.ID,
					Name:        video.Title,
					Authors:     video.Author,
					FeedName:    feed.Title,
					Description: video.Description,
					Media: []any{YouTubeVideo{
						Video:  video,
						Format: format,
					}},
					Published: video.PublishDate.UTC(),
					Updated:   video.PublishDate.UTC(),
					Raw:       video,
				})
			}
		}
	}

	var i int

	// Process all merged + sorted items.
	for _, item := range sortItems(items) {
		var err error

		// Normalize item name before any checks.
		item.Name, err = bp.NormalizeItemName(item.Name)
		if err != nil {
			return fmt.Errorf("error normalizing item name: %w", err)
		}

		// Check if item is duplicate.
		if bp.deduplicate(item.Name) {
			debugf("deduplicated => %s", item.Name)
			continue
		}

		var ok bool

		// Pass through positive filters.
		ok, err = bp.PositiveFilter(item.Name)
		if err != nil {
			return fmt.Errorf("error during positive filter: %w", err)
		}

		if !ok {
			debugf("+filtered => %s", item.Name)
			continue
		}

		// Pass through negative filters.
		ok, err = bp.NegativeFilter(item.Name)
		if err != nil {
			return fmt.Errorf("error during negative filter: %w", err)
		}

		if ok {
			debugf("-filtered => %s", item.Name)
			continue
		}

		// Calculate unique item directory name.
		dirname, err := bp.FormatItemPath(i, item)
		if err != nil {
			return fmt.Errorf("error formatting item path: %w", err)
		}

		// Sanitize item directory path name.
		dirname, err = bp.SanitizePath(dirname)
		if err != nil {
			return fmt.Errorf("error sanitizing item path: %w", err)
		}

		// Join item path with base dir.
		itemdir := path.Join(dir, dirname)

		// TODO: content map file changes feature

		// Process this item, writing metadata + downloading media.
		if err := bp.ProcessItem(ctx, itemdir, item); err != nil {
			errorf("error processing %q: %v", item.Name, err)
		}

		// Incr.
		i++
	}

	return nil
}

func (bp *beep) ProcessItem(ctx context.Context, dirpath string, item *Item) error {
	// Create item directory(s) path, ignores "already exists".
	if err := mkdirAll(dirpath, 0755); err != nil {
		return fmt.Errorf("error creating item dir: %w", err)
	}

	// Calculate serialized item update time.
	updated := item.Updated.Format(time.RFC3339)
	updatedPath := path.Join(dirpath, "updated")

	// Check if 'updatedAt' file exists with expected data on disk.
	// We use this to indicate whether each item needs redownloading.
	upToDate, err := fileExistsWith(updatedPath, updated)
	if err != nil {
		return fmt.Errorf("error checking updated_at: %w", err)
	}

	if upToDate {
		debugf("skipping up-to-date => %s", dirpath)
		return nil
	}

	// Write item description.
	if _, err := writeString(
		path.Join(dirpath, "description"),
		item.Description,
	); err != nil {
		return fmt.Errorf("error writing description: %w", err)
	}

	// Write item authors.
	if _, err := writeString(
		path.Join(dirpath, "authors"),
		item.Authors,
	); err != nil {
		return fmt.Errorf("error writing authors: %w", err)
	}

	var fail bool

	for _, media := range item.Media {
		switch media := media.(type) {
		case MediaEnclosure:
			// Get filename from the URL.
			name := getURLFilename(media.URL)

			// Sanitize media path filename.
			name, err := bp.SanitizePath(name)
			if err != nil {
				return fmt.Errorf("error sanitizing media name: %w", err)
			}

			if path.Ext(name) == "" {
				// Drop any unnecessary mimetype data.
				mimetype := dropExtraMimeData(media.Type)

				// Determine a file ext to use for media.
				ext, ok := mimetypes.GetFileExt(mimetype)
				if !ok {
					warnf("unexpected content type for %s: %s", item.Name, mimetype)
					ext = "bin" // just use .bin for unknown types
				}

				// Add ext to name.
				name += "." + ext
			}

			// Calculate on-disk file path.
			path := joinpath(dirpath, name)

			// Download the media to calculated path, error here is NOT fatal.
			if err := bp.DownloadEnclosure(ctx, path, media.Enclosure); err != nil {
				errorf("error downloading %s: %v", path, err)
				fail = true
			}

		case YouTubeVideo:
			// Sanitize youtube video title for filename.
			name, err := bp.SanitizePath(media.Video.Title)
			if err != nil {
				return fmt.Errorf("error sanitizing video name: %w", err)
			}

			// Drop any unnecessary mimetype data.
			mimetype := media.Format.MimeType
			mimetype = dropExtraMimeData(mimetype)

			// Determine a file ext to use for format.
			ext, ok := mimetypes.GetFileExt(mimetype)
			if !ok {
				warnf("unexpected content type for %s: %s", item.Name, mimetype)
				ext = "bin" // just use .bin for unknown types
			}

			// Calculate on-disk media file path.
			path := joinpath(dirpath, name+"."+ext)

			// Download the youtube video to calculated path, error here is NOT fatal.
			if err := bp.DownloadYouTube(ctx, path, media.Video, media.Format); err != nil {
				errorf("error downloading %s: %v", path, err)
				fail = true
			}
		}

		if bp.rateLimit {
			const sec = time.Second

			// Sleep for rand duration 0.5s -> 10s
			debugf("rate limiting for 0.5s ~ 10s")
			randSleep(sec/2, 10*sec)
		}
	}

	if fail {
		// Return early without marking as up-to-date, this allows later runs to re-attempt downloads.
		return nil
	}

	// Mark as being up to date.
	if _, err := writeString(
		updatedPath,
		updated,
	); err != nil {
		return fmt.Errorf("error writing updated_at: %w", err)
	}

	return nil
}

func (bp *beep) deduplicate(path string) bool {
	path = strings.TrimSpace(path)

	// Convert string input to
	// normalized composing form.
	bpath := byteutil.S2B(path)
	bpath = norm.NFC.Bytes(bpath)
	path = byteutil.B2S(bpath)

	// Check if path already in slice.
	for _, str := range bp.dedupe {
		if strings.EqualFold(str, path) {
			return true
		}
	}

	// Else add path to dedupe slice.
	bp.dedupe = append(bp.dedupe, path)
	return false
}

// DownloadEnclosure downloads a file from URL in feed enclosure, to the given path on disk. This also handles size mismatches.
func (bp *beep) DownloadEnclosure(ctx context.Context, path string, encls *gofeed.Enclosure) error {
	// Perform GET request for enclosure.
	rsp, err := bp.http.Get(encls.URL)
	if err != nil {
		return fmt.Errorf("error performing \"GET %s\":: %w", encls.URL, err)
	} else if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected \"GET %s\" response: %v", encls.URL, rsp.Status)
	}

	if rsp.ContentLength == 0 {
		// Are we being rate limited? Return
		// error here so can be attempted later.
		return errors.New("zero content-length")
	}

	infof("downloading enclosure => %s", path)

	// Write the opened stream to path on disk.
	if n, err := write(path, rsp.Body); err != nil {
		return fmt.Errorf("error writing enclosure stream: %w", err)
	} else if n != rsp.ContentLength {
		return fmt.Errorf("size mismatch: contentLength=%d written=%d", rsp.ContentLength, n)
	}

	return nil
}

// DownloadYouTube downloads a video from YouTube in the requested format, to the given path on disk. This also handles size mismatches.
func (bp *beep) DownloadYouTube(ctx context.Context, path string, video *youtube.Video, format *youtube.Format) error {
	// Get video stream for the selected video format.
	stream, sz, err := bp.youtube.GetStreamContext(ctx,
		video,
		format,
	)
	if err != nil {
		return fmt.Errorf("error getting youtube stream %q: %w", video.Title, err)
	}

	if sz == 0 {
		// Are we being rate limited? Return
		// error here so can be attempted later.
		return errors.New("zero content-length")
	}

	infof("downloading youtube video => %s", path)

	// Write the opened stream to path on disk.
	if n, err := write(path, stream); err != nil {
		return fmt.Errorf("error writing youtube stream: %w", err)
	} else if n != sz {
		return fmt.Errorf("size mismatch: contentLength=%d written=%d", sz, n)
	}

	return nil
}

// dropExtraMimeData drops all data after ';' in mime type.
func dropExtraMimeData(mimetype string) string {
	i := strings.Index(mimetype, ";")
	if i < 0 {
		return ""
	}
	return mimetype[:i]
}

// isYouTubeURL performs a simple URL parse
// and check if host = youtube.com / youtu.be.
func isYouTubeURL(u string) bool {
	url, err := url.Parse(u)
	if err != nil {
		return false
	}
	url.Host = strings.TrimPrefix(url.Host, "www.")
	switch url.Host {
	case "youtube.com":
		return true
	case "youtu.be":
		return true
	// TODO: better than this pathetic localization lol
	default:
		return false
	}
}

// getAuthorStr generates an author string from feed person object.
func getAuthorStr(authors []*gofeed.Person) string {
	buf := new(byteutil.Buffer)
	for _, author := range authors {
		buf.WriteString(author.Name)
		if author.Email != "" {
			buf.WriteString(" (")
			buf.WriteString(author.Email)
			buf.WriteString(")")
		}
		buf.WriteString(", ")
	}
	if buf.Len() > 0 {
		buf.Truncate(2)
	}
	return buf.String()
}

// getURLFilename gets the last path-part from a URL,
// dropping all extra fragment / query parts from end.
func getURLFilename(u string) string {
	u = path.Base(u)
	if i := strings.Index(u, "#"); i >= 0 {
		u = u[:i]
	}
	if i := strings.Index(u, "?"); i >= 0 {
		u = u[:i]
	}
	return u
}

// getYouTubeVideoFormat selects the video format with both audio/video of highest quality.
func getYouTubeVideoFormat(video *youtube.Video) *youtube.Format {
	formats := video.Formats

	// Delete all available formats without audio / video.
	formats = slices.DeleteFunc(formats, func(f youtube.Format) bool {
		return f.AudioChannels <= 0 || f.FPS <= 0
	})

	if len(formats) == 0 {
		return nil
	}

	// Sort by quality.
	formats.Sort()

	// Return highest quality.
	return &(formats[0])
}

// randSleep sleeps for duration between min and max.
func randSleep(min, max time.Duration) {
	// Read random int from 0,max-min.
	b := big.NewInt(int64(max - min))
	i, err := rand.Int(rand.Reader, b)

	if err != nil {
		// On error just
		// use big max.
		i = b
	}

	// Get value starting at min.
	d := int64(min) + i.Int64()

	// Sleep for random duration.
	time.Sleep(time.Duration(d))
}

// joinpath joins path parts, briefly sanitizing outer part.
func joinpath(base string, part string) string {
	part = strings.ReplaceAll(part, "/", " ")
	return path.Join(base, part)
}

func withUserAgent(rt http.RoundTripper, userAgent string) http.RoundTripper {
	if rt == nil {
		panic("nil")
	}
	return roundtripper(func(r *http.Request) (*http.Response, error) {
		r.Header.Set("User-Agent", userAgent)
		return rt.RoundTrip(r)
	})
}

type roundtripper func(*http.Request) (*http.Response, error)

func (fn roundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
