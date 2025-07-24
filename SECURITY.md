# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, I encourage responsible disclosure.

- **Email**: [security@rafayhingoro.me](mailto:security@rafayhingoro.me)
- **Please include**:
  - A clear description of the issue
  - Steps to reproduce (if possible)
  - Potential impact and suggested fixes (if known)

Please **do not open GitHub issues for security concerns**. I aim to respond within **5 working days**.

## Supported Versions

| Version       | Status            |
| ------------- | ----------------- |
| `main` branch | ✅ Supported      |
| older tags    | ❌ Not maintained |

Only the latest release or `main` branch is supported for security updates.

## Disclosure Process

I follow a coordinated disclosure process:
1. Vulnerability is reported privately.
2. I investigate and confirm the issue.
3. A fix is developed and released.
4. (Optional) Acknowledgment is given in release notes or this file.

## Notes for Users

- Keep your dependencies up to date.
- Use tools like [`govulncheck`](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) or [`gosec`](https://github.com/securego/gosec) to monitor for known vulnerabilities.
- This project does not have a formal bug bounty program, but I appreciate reports that help improve security.
