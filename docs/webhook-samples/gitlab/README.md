# Headers (event types)

These are the three most common headers that GitLab sends with its webhooks.

- Push Hook
    ```
    X-Gitlab-Event: Push Hook
    ```

- Push Tag Hook
    ```
    X-Gitlab-Event: Tag Push Hook
    ```

- Merge Request Hook
    ```
    X-Gitlab-Event: Merge Request Hook
    ```

More information about GitLab webhooks can be found [here](https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html).

# Payload

There are several examples of payloads that GitLab sends with its webhooks in this directory. The payloads are in JSON format and are named according to the event type that they represent.