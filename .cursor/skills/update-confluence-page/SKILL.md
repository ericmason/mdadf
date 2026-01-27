---
name: update-confluence-page
description: Updates Confluence pages using the REST API v2 with ADF format. Use when updating Confluence pages, creating pages with ADF content, or when the user asks to update or create Confluence pages with markdown/ADF content.
---

# Update Confluence Pages

Updates Confluence pages using the REST API v2 with Atlassian Document Format (ADF) content.

## Authentication

Get the API token from macOS Keychain:

```bash
TOKEN=$(security find-generic-password -l "acli" -w 2>&1 | sed 's/^go-keyring-base64://' | base64 -d)
```

The token is stored by `acli` and can be extracted this way.

## Update Page with ADF Content

### Step 1: Prepare ADF Content

Convert markdown to ADF JSON using the `mdadf` utility:

```bash
./mdadf -c input.md > adf_content.json
```

### Step 2: Create Update Payload

Create a JSON payload with the ADF content:

```python
import json

# Read ADF content
with open('adf_content.json', 'r') as f:
    adf_doc = json.load(f)

# Create update payload
update_payload = {
    "id": "PAGE_ID",
    "status": "current",
    "title": "Page Title",
    "body": {
        "representation": "atlas_doc_format",
        "value": json.dumps(adf_doc)  # ADF as JSON string, not object
    },
    "version": {
        "number": CURRENT_VERSION + 1  # Must increment version
    }
}

# Save payload
with open('update_payload.json', 'w') as f:
    json.dump(update_payload, f, indent=2)
```

**Important**: 
- `value` must be a JSON **string** (use `json.dumps()`), not a JSON object
- Version number must be incremented from current version
- Use `atlas_doc_format` representation for ADF content

### Step 3: Make API Call

```bash
TOKEN=$(security find-generic-password -l "acli" -w 2>&1 | sed 's/^go-keyring-base64://' | base64 -d)

curl -X PUT "https://equisolve.atlassian.net/wiki/api/v2/pages/{PAGE_ID}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -u "eric@equisolve.com:$TOKEN" \
  -d @update_payload.json
```

## Create New Page

For creating a new page, use POST instead:

```python
create_payload = {
    "spaceId": "SPACE_ID",
    "status": "current",
    "title": "Page Title",
    "body": {
        "representation": "atlas_doc_format",
        "value": json.dumps(adf_doc)
    }
}
```

```bash
curl -X POST "https://equisolve.atlassian.net/wiki/api/v2/pages" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -u "eric@equisolve.com:$TOKEN" \
  -d @create_payload.json
```

## Common Errors

- **409 CONFLICT**: Version number not incremented. Get current version and increment it.
- **400 INVALID_MESSAGE**: `value` must be a JSON string, not an object. Use `json.dumps()`.
- **401 UNAUTHORIZED**: Token extraction failed. Verify acli is authenticated.

## Get Current Page Version

To get the current version before updating:

```bash
curl -s "https://equisolve.atlassian.net/wiki/api/v2/pages/{PAGE_ID}" \
  -H "Accept: application/json" \
  -u "eric@equisolve.com:$TOKEN" | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['version']['number'])"
```

## Complete Example

```bash
# 1. Convert markdown to ADF
./mdadf -c plan.md > adf.json

# 2. Get current version
CURRENT_VERSION=$(curl -s "https://equisolve.atlassian.net/wiki/api/v2/pages/2254307333" \
  -H "Accept: application/json" \
  -u "eric@equisolve.com:$TOKEN" | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['version']['number'])")

# 3. Create update payload
python3 << 'PYEOF'
import json
with open('adf.json', 'r') as f:
    adf_doc = json.load(f)
update_payload = {
    "id": "2254307333",
    "status": "current",
    "title": "Page Title",
    "body": {
        "representation": "atlas_doc_format",
        "value": json.dumps(adf_doc)
    },
    "version": {"number": CURRENT_VERSION + 1}
}
with open('update.json', 'w') as f:
    json.dump(update_payload, f, indent=2)
PYEOF

# 4. Update page
TOKEN=$(security find-generic-password -l "acli" -w 2>&1 | sed 's/^go-keyring-base64://' | base64 -d)
curl -X PUT "https://equisolve.atlassian.net/wiki/api/v2/pages/2254307333" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -u "eric@equisolve.com:$TOKEN" \
  -d @update.json
```

## Notes

- Confluence accepts ADF format via `atlas_doc_format` representation
- The MCP tool only supports `storage` format, so use REST API directly for ADF
- Confluence converts ADF to storage format internally (this is normal)
- API endpoint: `/wiki/api/v2/pages/{id}` for updates, `/wiki/api/v2/pages` for creates
