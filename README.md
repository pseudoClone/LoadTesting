# Details

Until now, the server is generic like `https://example.com` and for local I use my Python http module to run a server listening on port 8000 like this:

```bash
uv run python -m http.server 8000
```

But going forward I am using a Go server that can handle mutiple connections at a time.


> Note: So, I made a customServer. And the request body is just the string sent over the connection to the client. I still don't know how it is 250 bytes when the length of body is 24. I guess it is padding or other headers
> Also, Python's http.server is not handling requests good. The docs say that there is an internal buffer for request. I guess it is just creating some sort of race condition in the queue/buffer.

## Usage
Build tester using
```bash
go build -o ldtst.exe .\cmd\customClient\
```

Run the test server using
```bash
go run .\cmd\customServer
```
Run it 
`ldtst.exe -s https://localhost:8000 -n 2`

# **NOTE**
## Responsible Use

This tool is intended **only for authorized load and performance testing** of systems that you own or have explicit permission to test.

Do **not** use this software to target third-party services, websites, or networks without authorization. Unauthorized stress testing or denial-of-service attacks may violate applicable laws, service terms, or organizational policies.

The author does not endorse or encourage the use of this tool for disruptive or malicious activities. Users are solely responsible for ensuring that their use of this software complies with all applicable laws and that they have permission to test the target systems.

Use this software responsibly and at your own risk.
