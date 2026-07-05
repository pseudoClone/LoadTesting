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


# Development Notes

This started out as a way to learn Go using only *The Go Programming Language* by Brian Kernighan & Alan Donovan and [Go documentation](https://pkg.go.dev/net/http) while completely blocking out AI and chatbots, even for concepts.

I am leaving this note mostly for myself. Up until this point, no AI has been used and if I add anything other than Requests Per Seconds, it might probably involve AI assistance unless I explicitly say otherwise. Maybe because I have tried and failed in making a Token Buckeet implementation.

Nothing really against AI but really it's getting out of hand now. Why would we commoditize our thinking in tokens? 