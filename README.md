# feedloader

Suckless podcast downloader. Now with support for podcasts contained in YouTube playlists!


Usage:

```
Usage: ./beep...
    --dry-run
    	Only perform dry-run (no media downloads)
 -h --help
    	Print usage information
    --item-format shell-expr (default: echo $INDEX - $NAME)
    	Feed item path formatting (input = env $INDEX,$NAME,$FEED, result = stdout)
    --item-normalize shell-expr (default: tr -s "[:space:]" " ")
    	Feed item name normalizer (input = stdin, result = stdout)
    --negative-filter shell-expr
    	Feed item name negative filter (input = stdin, result = exit code)
 -o --output-dir string (default: .)
    	Output directory
    --positive-filter shell-expr
    	Feed item name positive filter (input = stdin, result = exit code)
    --rate-limit
    	Rate limit (by rand interval) between downloads
    --sanitize-path shell-expr (default: tr -d "/|" | iconv --from-code=UTF-8 -c)
    	Filesystem path sanitizer (input = stdin, result = stdout)
 -u --url []string
    	Feed URL address(es). Can be RSS, atom, YouTube Playlist
    --user-agent string (default: feedloader/0.1)
    	HTTP client User-Agent
```

Example use cases:

- Downloading Fall of Civilizations podcast episodes, which are already numbered:
```
./beep--url 'https://feeds.soundcloud.com/users/soundcloud:users:572119410/sounds.rss' \
             --item-format "echo ${TITLE}" \
             --output-dir 'Fall of Civilizations'
```
