# LEA — Full Changelog

> For the complete list of changes, see [all releases](https://github.com/KopyTKG/LEA/releases).

---

## [v3.2.0](https://github.com/KopyTKG/LEA/releases/tag/v3.2.0) — 2025-07-01

- **[Security] Added CMAC authentication support** ([#28](https://github.com/KopyTKG/LEA/issues/28))
  - Integrated CMAC for message authentication, ensuring integrity and authenticity of encrypted data.
  - CLI option to enable/disable CMAC verification.
  - Documentation and examples updated.

---

## [v3.1.0](https://github.com/KopyTKG/LEA/releases/tag/v3.1.0) — 2025-06-30

- Official Windows-compatible build ([#30](https://github.com/KopyTKG/LEA/issues/30), [#31](https://github.com/KopyTKG/LEA/issues/31))
- Automatic IV generation for each encryption ([#29](https://github.com/KopyTKG/LEA/issues/29))
- Expanded metadata struct (version, mode, IV fields)
- Improved decryption output (files match original size exactly)
- Robust handling for padding and binary streams
- Better error handling and logging
- Codebase refactoring for safety and maintainability

---

## [v3.0.0](https://github.com/KopyTKG/LEA/releases/tag/v3.0.0) — 2025-06-23

### 🚨 Breaking Changes

- **New Key Management System**
  - Introduced armored key file format with metadata and SHA-256 digests for integrity.
  - New `lea keygen` command for secure keyfile generation.
  - Encryption/decryption now require `--key=<file>`. Old `-ek` and `-es` flags deprecated.
  - Legacy seed/key input no longer supported.
- **CLI Redesign**
  - Overhauled CLI arguments and help output for clarity.
  - Flags are now more consistent and descriptive.
- **New Encryption Modes**
  - Added CTR mode (alongside CBC).
  - Mode selection via CLI.
- **File Wiping**
  - Added secure multi-pass file wiping following DoD 5220.22-M standard.
- **Progress UI**
  - Improved progress and status reporting with hooks for UI integration.

### 🐞 Bug Fixes & Enhancements

- Deprecated legacy key/seed flags and logic.
- Refactored key schedule and digesting for modularity and safety.
- Improved modularity and codebase maintainability.
- Updated README with new usage, key management, and legal disclaimer.

---

## [v2.0.1](https://github.com/KopyTKG/LEA/releases/tag/v2.0.1) — 2025-06-22

- Updated documentation for modernized CLI and usage.
- Improved error handling during key validation and file access.
- Minor CLI help improvements.
- Removed legacy GitHub workflow files.
- Synchronized branch with `main` for v2 support.

---

## [v2.0.0](https://github.com/KopyTKG/LEA/releases/tag/v2.0.0) — 2025-06-17

- **Major rewrite** for codebase modularity.
- Split encryption logic into separate modules for algorithm, I/O, and CLI.
- Initial support for new key file format (transition period).
- Improved command-line parsing and argument validation.
- Preparation for upcoming v3.0.0 changes.

---

## [v1.6.1](https://github.com/KopyTKG/LEA/releases/tag/v1.6.1) — 2024-08-15

- Added file fingerprinting for output verification.
- Improved logging for encryption/decryption operations.
- Updated README and documentation.
- Fixed merge conflicts and streamlined release workflow.
- Added initial GitHub Actions CI.

---

## [v1.4.0 (CFB,OFB)](https://github.com/KopyTKG/LEA/releases/tag/v1.4.0) — 2024-06-26

- Added support for CFB and OFB block cipher modes.
- CLI flags for selecting mode (`--mode=cfb`, `--mode=ofb`).
- Improved internal test coverage for multiple modes.
- Updated help and mode documentation.

---

## [v1.2.0](https://github.com/KopyTKG/LEA/releases/tag/v1.2.0) — 2024-06-25

- Refactored block encryption logic for performance.
- Improved error reporting for invalid key/seed.
- Added CLI argument validation.
- Updated Readme with usage examples.

---

## [v1.1.0](https://github.com/KopyTKG/LEA/releases/tag/v1.1.0) — 2024-06-07

- Initial public release of LEA command line encryption tool.
- Supports LEA block cipher (128/192/256-bit keys).
- Basic CBC mode encryption/decryption.
- Command-line interface for file-based operations.
- Basic key/seed management (legacy).

---

## [v1.0.0](https://github.com/KopyTKG/LEA/releases/tag/v1.0.0) — 2024-06-06

- First production-ready version.
- Core LEA algorithm implementation.
- Basic CLI for encryption/decryption.

---
