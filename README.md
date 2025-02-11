# valiDTr - Git Commit Verification CLI

**valiDTr** is a command-line tool designed to verify Git commit chains using GPG signatures. It ensures commit authenticity and helps maintain a secure Git workflow.

---

## 📌 Installation

### Clone the repository and build:
```sh
git clone https://github.com/yourusername/valiDTr.git
cd valiDTr
go build -o valiDTr
```
## 🛠️ CLI Commands

### Verify a git commit signature

```sh
valiDTr verify <commit-hash>
# verifies a Git commit has a vild GPG signature
# Example:
valiDTr verify a1b2c3d4e5
# Output:
Good signature from "Developer Name <dev@example.com>"
```

### Add a GPG Key

```sh
valiDTr add-key <gpg-key-id>
# Adds a developer's GPG key to the verification database
# Example:
valiDTr add-key ABCD1234
# Output:
Key ABCD1234 added successfully!
```

### Revoke a GPG Key

```sh
valiDTr revoke-key <gpg-key-id>
# Revokes a developer's GPG key without affecting past commits
# Example:
valiDTr revoke-key ABCD1234
# Output:
Key ABCD1234 has been revoked.
```

### List Stored GPG Keys

```sh
valiDTr list-keys
# Lists all stored GPG keys along with their status (active/revoked)
# Example:
valiDTr list-keys
# Output:
Stored GPG Keys:
----------------------------
Key ID: ABCD1234 | Status: active
Key ID: WXYZ5678 | Status: revoked
```

### Initialize the Database

```sh
valiDTr init-db
# Creates the SQLite database and initializes the GPG key storage table
# Example:
valiDTr init-db
# Output:
Database initialized successfully!
```

### Check Installed Version

```sh
valiDTr version
# Displays the installed version of valiDTr
# Example:
valiDTr version
# Output:
valiDTr v1.0.0
```

### 🚀 Usage Example
```sh
# Step 1: Initialize the database
valiDTr init-db

# Step 2: Add a GPG key for verification
valiDTr add-key ABCD1234

# Step 3: Verify a commit
valiDTr verify a1b2c3d4e5

# Step 4: List stored keys
valiDTr list-keys

# Step 5: Revoke a key
valiDTr revoke-key ABCD1234
```

### 🎯 Features Roadmap

- [x] Verify Git commit signatures (initial)
- [x] Store & manage GPG keys (initial)
- [x] Revoke keys while preserving past commit validity (initial)
- [] Git Hook Integration (coming soon)
- [] REST API Support (coming soon)
- [] CI/CD Integration (coming soon)

### 🔧 Troubleshooting

#### Error: "GPG key not found"

Solution: Ensure the key exists in your GPG keyring:
```sh
gpg --list-keys
```

#### Error: "Git commit signature verification failed"

Solution: Check if the commit is signed correctly:
```sh
git log --show-signature -1 <commit-hash>
```

### Contributers

esuarkeN (Eric Krause)
