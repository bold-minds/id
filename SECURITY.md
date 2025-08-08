# Security Policy

## Supported Versions

We actively support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please follow these steps:

### 1. **Do Not** Create a Public Issue

Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.

### 2. Report Privately

Send an email to **security@bold-minds.com** with the following information:

- **Subject**: Security Vulnerability in bold-minds/id
- **Description**: Detailed description of the vulnerability
- **Steps to Reproduce**: Clear steps to reproduce the issue
- **Impact**: Potential impact and severity assessment
- **Suggested Fix**: If you have ideas for a fix (optional)

### 3. Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Varies based on complexity, typically within 30 days

### 4. Disclosure Process

1. We will acknowledge receipt of your vulnerability report
2. We will investigate and validate the vulnerability
3. We will develop and test a fix
4. We will coordinate disclosure timing with you
5. We will release a security update
6. We will publicly acknowledge your responsible disclosure (if desired)

## Security Considerations

### Entropy Sources

This library provides multiple entropy source options:

- **Default**: Uses `math/rand` with time-based seeding (suitable for most applications)
- **Secure**: Uses `crypto/rand` for cryptographically secure randomness
- **Custom**: Allows you to provide your own entropy source

For security-sensitive applications, always use `NewSecureGenerator()`:

```go
// For security-sensitive applications
gen := id.NewSecureGenerator()
```

### ULID Properties

ULIDs have the following security-relevant properties:

- **Predictable Timestamp**: The first 48 bits encode timestamp in milliseconds
- **Random Component**: The remaining 80 bits are random (when using appropriate entropy)
- **Not Cryptographically Secure**: ULIDs are not designed to be cryptographically secure identifiers

### Best Practices

1. **Use Secure Generation**: For sensitive applications, use `NewSecureGenerator()`
2. **Validate Input**: Always validate ULIDs from external sources using `IsIdValid()`
3. **Handle Errors**: Properly handle all error returns from library functions
4. **Avoid Timing Attacks**: Be aware that timestamp extraction reveals creation time
5. **Rate Limiting**: Consider rate limiting ULID generation in public APIs

### Known Limitations

- ULIDs reveal approximate creation time
- Monotonic ordering within the same millisecond depends on entropy source
- Not suitable as cryptographic tokens or passwords
- Should not be used for security-critical random number generation

## Security Updates

Security updates will be:

- Released as patch versions (e.g., 1.0.1)
- Documented in the CHANGELOG.md
- Announced through GitHub releases
- Tagged with security labels

## Acknowledgments

We appreciate responsible disclosure and will acknowledge security researchers who help improve the security of this project.

Thank you for helping keep our project and users safe!
