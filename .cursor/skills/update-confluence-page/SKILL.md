---
name: update-confluence-page
description: Updates Confluence pages with ADF content produced by mdadf, through the Atlassian MCP (or the REST API v2 as a fallback). Use when updating Confluence pages, creating pages with ADF content, or when the user asks to update or create Confluence pages with markdown/ADF content.
---

# Update Confluence Pages

Updates Confluence pages with Atlassian Document Format (ADF) content that `mdadf` converts from Markdown. Use the Atlassian MCP first. Use the REST API only if the MCP isn't available.

## Step 1: Convert Markdown to ADF

Convert markdown to ADF JSON using the `mdadf` utility:

```bash
./mdadf -c input.md > adf_content.json
```

## Step 2: Update the page with the Atlassian MCP

The Atlassian MCP's `updateConfluencePage` and `createConfluencePage` tools take ADF directly: set `contentFormat: "adf"` and pass the ADF document as a JSON string in `body`.

1. Call `getConfluencePage` with `cloudId: "equisolve.atlassian.net"` and the `pageId` to confirm the page and read its current title.
2. Call `updateConfluencePage` with:
   - `cloudId`: `"equisolve.atlassian.net"`
   - `pageId`: the page ID
   - `title`: the page title (keep the current title unless you're renaming the page)
   - `contentFormat`: `"adf"`
   - `body`: the contents of `adf_content.json`, as a string
   - `versionMessage`: a short description of the change (optional)

The MCP handles the version number, so you don't send one.

To create a page instead, call `createConfluencePage` with `cloudId`, `spaceId`, `title`, `contentFormat: "adf"`, and `body` set the same way. Pass `parentId` to nest it under another page.

If the `mcp__atlassian__*` tools aren't available, follow `mcp-setup.md` in the `jira` skill (equisolve/eq-agent-skills, `plugins/eq-dev/skills/jira/`), or use the REST API below.

## REST API fallback

Call the Confluence REST API v2 with the `jira-api` helper. It comes from the `jira` skill in equisolve/eq-agent-skills (`plugins/eq-dev/skills/jira/scripts/jira-api`). It reads the email and API token from `~/.config/jira/email` and `~/.config/jira/api-token` and keeps them out of the command line and output, and Confluence on the same site accepts the same token. Never read the token file yourself.

The helper takes any path on the site: `jira-api [METHOD] <path> [json-body]`, where `json-body` is a JSON string or `@file`.

### Create Update Payload

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

### Make API Call

```bash
jira-api PUT /wiki/api/v2/pages/{PAGE_ID} @update_payload.json
```

### Create New Page

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
jira-api POST /wiki/api/v2/pages @create_payload.json
```

### Common Errors

- **409 CONFLICT**: Version number not incremented. Get current version and increment it.
- **400 INVALID_MESSAGE**: `value` must be a JSON string, not an object. Use `json.dumps()`.
- **401 UNAUTHORIZED**: The email or API token in `~/.config/jira/` is wrong or expired. Run `jira-api --help` for where to create a new token.

### Get Current Page Version

To get the current version before updating:

```bash
jira-api /wiki/api/v2/pages/{PAGE_ID} | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['version']['number'])"
```

To read the current body as ADF, add `?body-format=atlas_doc_format` to the path.

### Complete Example

```bash
# 1. Convert markdown to ADF
./mdadf -c plan.md > adf.json

# 2. Get current version
CURRENT_VERSION=$(jira-api /wiki/api/v2/pages/2254307333 | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['version']['number'])")

# 3. Create update payload
CURRENT_VERSION=$CURRENT_VERSION python3 << 'PYEOF'
import json, os
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
    "version": {"number": int(os.environ["CURRENT_VERSION"]) + 1}
}
with open('update.json', 'w') as f:
    json.dump(update_payload, f, indent=2)
PYEOF

# 4. Update page
jira-api PUT /wiki/api/v2/pages/2254307333 @update.json
```

## Notes

- The Atlassian MCP's `updateConfluencePage` and `createConfluencePage` accept `contentFormat: "adf"`, so you don't need the REST API to send ADF
- Confluence accepts ADF format via `atlas_doc_format` representation in the REST API
- Confluence converts ADF to storage format internally (this is normal)
- API endpoint: `/wiki/api/v2/pages/{id}` for updates, `/wiki/api/v2/pages` for creates
