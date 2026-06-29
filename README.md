# Details

Until now, the server is generic like `https://example.com` and for local I use my Python http module to run a server listening on port 8000 like this:

```bash
uv run python -m http.server 8000
```

But going forward I am using a Go server that can handle mutiple connections at a time