# space-dl

Utility to download audio recording of a twitter space that's ended

Usage:

```sh
$ go build
$ youtube-dl $(./space-dl -id 1eaKbNYzMzjKX)
```

## TODO

- [ ] Native go downloader
- [ ] Use http2 for faster downloads
- [ ] Better error messages
- [ ] Get space url from tweet url