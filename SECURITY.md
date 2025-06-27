# Security Policy

## Supported Versions

We actively support and maintain the latest version of LEA. Please use the latest release for all production deployments.

| Version | Supported          |
| ------- | ----------------- |
| Latest  | :white_check_mark: |
| Older   | :x:               |

If you require long-term support or backports for older versions, please contact the maintainers.

---

## Reporting a Vulnerability

If you discover a security vulnerability, **do not open a public issue**.  
Please use one of the following confidential channels:

- Email: [contact@thekrew.app]
- GitHub ["Report a security vulnerability"](../../security/advisories/new) feature

Please provide us with:

- A description of the vulnerability and its impact
- Steps to reproduce or proof-of-concept code, if possible
- Any proposed mitigations or fixes

We will acknowledge your report within **3 business days** and work with you on a coordinated disclosure and fix.

---

## Project Security Posture

- **Cryptographic Review:**  
  This project implements the Lightweight Encryption Algorithm (LEA) as specified in publicly available standards. While care has been taken in the implementation, it has not undergone a formal third-party cryptographic audit at this time.  
  **If you are a cryptography professional and are interested in reviewing or auditing this project, please get in touch.**

- **Production Readiness:**  
  We strive to ensure that this repository is safe and robust for production use. However, cryptographic software is inherently complex, and new vulnerabilities may be discovered over time.  
  Users are encouraged to review the source code, conduct independent audits, and follow best practices for deployment.

- **Security Updates:**  
  Security issues will be fixed promptly and new releases will be published as needed.

---

## Responsible Disclosure

We value and encourage responsible disclosure.  
We will keep all reports confidential, coordinate fixes before public disclosure, and credit reporters (unless anonymity is requested).

---

## Security Best Practices for Users

- Always use the latest version of LEA.
- Review [RELEASE NOTES](./CHANGELOG.md) for security updates.
- Deploy LEA with secure key management and in accordance with best cryptographic practices.
- Do not share secrets or keys in public issues or pull requests.

---

## References

- [GitHub Security Advisories](https://docs.github.com/en/code-security/security-advisories)
- [Responsible Disclosure Guidelines](https://responsibledisclosure.com/)
- [LEA Algorithm Specification](https://en.wikipedia.org/wiki/LEA_(cipher))

Thank you for helping us keep LEA secure!
