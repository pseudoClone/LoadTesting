# Details

Until now, the server is generic like `https://example.com` and for local I use my Python http module to run a server listening on port 8000 like this:

```bash
uv run python -m http.server 8000
```

But going forward I am using a Go server that can handle mutiple connections at a time.


> Note: So, I made a customServer. And the request body is just the string sent over the connection to the client. I still don't know how it is 250 bytes when the length of body is 24. I guess it is padding or other headers
> Also, Python's http.server is not handling requests good. The docs say that there is an internal buffer for request. I guess it is just creating some sort of race condition in the queue/buffer.