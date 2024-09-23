# Webhook Samples Directory

The `docs/webhook-samples` directory contains sample JSON payloads for various webhook events. These samples are useful for testing and understanding the structure of webhook events that your application might receive.

## Purpose

- **Testing**: Use these sample payloads to simulate webhook events during development and testing.
- **Documentation**: Understand the structure and content of different webhook events.
- **Debugging**: Compare actual webhook payloads with these samples to identify discrepancies or issues.

## Contents

The directory includes JSON files representing different types of webhook events. Each file is named according to the event it represents. For example:

- `pr-opened.json`: Payload for a "Pull Request Opened" event.
- `push.json`: Payload for a "Push" event.

## Usage

1. **Development**: Use these samples to mock webhook events in your local development environment.
2. **Testing**: Integrate these samples into your automated tests to ensure your application handles webhook events correctly.
3. **Reference**: Refer to these samples to understand the data structure and fields included in different webhook events.

## Example

Here is an example of how you might use a sample payload in your tests:

```bash
#!/bin/bash

# Define the URL to send the webhook to
WEBHOOK_URL="https://example.com/webhook"

# Define the path to the JSON file
JSON_FILE="docs/webhook-samples/github/pr-opened.json"

# Send the JSON payload using curl
curl -X POST -H "Content-Type: application/json" -d @"$JSON_FILE" "$WEBHOOK_URL"
```