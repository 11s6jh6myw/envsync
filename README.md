# envsync

> Diff and sync `.env` files across environments with redaction support.

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git
cd envsync && go build -o envsync .
```

---

## Usage

**Diff two `.env` files:**

```bash
envsync diff .env.staging .env.production
```

**Sync missing keys from one file to another:**

```bash
envsync sync .env.example .env.local
```

**Redact sensitive values before displaying output:**

```bash
envsync diff .env.staging .env.production --redact
```

Output example:

```
+ DB_HOST=db.production.internal
~ API_URL: https://staging.api.com → https://api.com
- DEBUG=true
```

Keys marked with `SECRET`, `PASSWORD`, `TOKEN`, or `KEY` are automatically masked when `--redact` is enabled.

---

## Flags

| Flag | Description |
|-----------|-------------------------------|
| `--redact` | Mask sensitive values in output |
| `--dry-run` | Preview sync changes without writing |
| `--output` | Output format: `text`, `json` |

---

## License

[MIT](LICENSE) © 2024 yourusername