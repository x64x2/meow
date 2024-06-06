package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/x64x2/go-fflag"
	"github.com/x64x2/go-mimetypes"

	"golang.org/x/sys/unix"
)

var (
	// std file descriptors.
	stdout = os.Stdout
	stderr = os.Stderr

	// dryrun prevents download
	// of feed media (for tests).
	dryrun bool
)

func init() {
	// Add the often-used audio/mp3 mimetype too.
	mimetypes.PreferredExts["audio/mp3"] = "mp3"
	mimetypes.PreferredExts["video/mp4"] = "mp4"
}

func main() {
	var code int

	// Run main application
	if err := run(); err != nil {
		fatalf("%v", err)
		code = 1
	}

	// Exit with code
	unix.Exit(code)
}

func run() error {
	const (
		defaultUserAgent = `feedloader/0.1`
		defaultItemPath  = `echo $INDEX - $NAME`
		defaultPosFilter = ``
		defaultNegFilter = ``
		defaultNormalize = `tr -s "[:space:]" " "`
		defaultSanitize  = `tr -d "/|" | iconv --from-code=UTF-8 -c`
	)

	var (
		feedURLs  []string
		userAgent string
		outdir    string
	)

	// Declare runtime flags to global fflag.FlagSet
	fflag.StringSliceVar(&feedURLs, "u", "url", nil, "Feed URL address(es). Can be RSS, atom, YouTube Playlist")
	fflag.StringVar(&userAgent, "", "user-agent", defaultUserAgent, "HTTP client User-Agent")
	fflag.BoolVar(&dryrun, "", "dry-run", false, "Only perform dry-run (no media downloads)")
	fflag.StringVar(&outdir, "o", "output-dir", ".", "Output directory")
	fflag.Help()

	// Parse flags, don't allow unrecognized.
	if xtra, err := fflag.Parse(); err != nil {
		return err
	} else if len(xtra) > 0 {
		return fmt.Errorf("unexpected args: %v", xtra)
	}

	// Ensure output dir exists.
	err := mkdirAll(outdir, 0755)
	if err != nil {
		return fmt.Errorf("error making output dir: %w", err)
	}

	// Create new program context.
	ctx := context.Background()
	ctx, cncl := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cncl()

	// Initialize loader
	// underlying clients.
	loader.Init(userAgent)

	