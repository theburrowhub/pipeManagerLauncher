# webhook-listener

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.2](https://img.shields.io/badge/AppVersion-0.0.2-informational?style=flat-square)

This Helm Chart deploys pipeManager, a Go-based application designed to automate pipeline management
using Tekton CI/CD. The deployment includes:

1. Webhook Service: Captures webhook events that trigger the creation and execution of Kubernetes jobs.
   These jobs prepare a pipeline that will be managed and executed by Tekton.
2. Cleanup Service: Responsible for deleting Kubernetes jobs and other resources generated during the
   event processing, ensuring a clean and efficient environment.

This chart is ideal for automating pipelines in a Kubernetes environment, maintaining control over the creation
and cleanup of related resources.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Configuration for the affinity. |
| args[0] | string | `"-c"` |  |
| args[1] | string | `"/etc/pipe-manager/config.yaml"` |  |
| args[2] | string | `"-l"` |  |
| args[3] | string | `":80"` |  |
| autoscaling | object | `{"enabled":false,"maxReplicas":100,"minReplicas":1,"targetCPUUtilizationPercentage":80}` | Configuration for the auto-scaling. |
| autoscaling.enabled | bool | `false` | Enable auto-scaling |
| autoscaling.maxReplicas | int | `100` | Maximum number of pods to run |
| autoscaling.minReplicas | int | `1` | Minimum number of pods to run |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | Target CPU utilization percentage |
| command[0] | string | `"/app/webhook-listener"` |  |
| config.common | object | `{"log":{"file":"stdout","format":"json","level":"info"}}` | Common configuration for the applications of pipeManager. |
| config.common.log | object | `{"file":"stdout","format":"json","level":"info"}` | Logging configuration for the webhook server. |
| config.common.log.file | string | `"stdout"` | Destination for log output. Options: - "stdout": Logs are printed to the standard output (useful for containerized environments). - "stderr": Logs are printed to the standard error. - File path (e.g., "/var/log/webhook.log"): Logs are written to the specified file.  Example: - Use "stdout" if your deployment environment captures and manages logs. - Specify a file path to persist logs on the server. |
| config.common.log.format | string | `"json"` | Format of the log messages. Options: - "text": Plain text logs, human-readable. - "json": Logs are formatted as JSON objects, useful for structured logging and integration with logging systems.  Example: - Choose "text" for simplicity and ease of reading during development. - Choose "json" if you are using log aggregators or need structured logs for analysis. |
| config.common.log.level | string | `"info"` | Logging level for the webhook server. Common levels: - "debug": Detailed information, typically used for diagnosing problems. - "info": Confirmation that things are working as expected. - "warning": An indication that something unexpected happened, or indicative of some problem. - "error": A more serious problem, preventing some function from performing. - "critical": A severe error, indicating that the program itself may be unable to continue running.  Example: - Set to "debug" for verbose logging during development or troubleshooting. - Set to "info" or "warning" for production environments to reduce log verbosity. |
| config.launcher | object | `{"cloneDepth":0,"imageName":"k3d-registry:5111/pipeline-converter","jobNamePrefix":"pipeline-launcher","namespace":"","pullPolicy":"IfNotPresent","tag":"0.0.2","timeout":600}` | Configuration for Pipeline Launcher |
| config.launcher.cloneDepth | int | `0` | Clone depth for the pipeline job. The number of commits to fetch from the repository. Adjust this value based on the size of the repository and the required history for the pipeline job. A value of 0 fetches the entire history of the repository. |
| config.launcher.imageName | string | `"k3d-registry:5111/pipeline-converter"` | Docker image to use for launching pipeline jobs. |
| config.launcher.jobNamePrefix | string | `"pipeline-launcher"` | Prefix for the job name created for each webhook event. The job name will be appended with a unique identifier. (e.g., pipeline-launcher-<unique-id>). Ensure that the prefix does not exceed 25 characters to avoid Kubernetes naming restrictions of 63 characters. |
| config.launcher.namespace | string | `""` | Namespace where the pipeline will be launched. Leave empty to use the current namespace where the webhook server is deployed. |
| config.launcher.pullPolicy | string | `"IfNotPresent"` | Image pull policy for the pipeline job. Options: - "Always": Always pull the image, even if it exists locally. - "IfNotPresent": Pull the image only if it does not exist locally. - "Never": Never pull the image, only use it if it exists locally.  Example: - Set to "IfNotPresent" to avoid pulling the image repeatedly if it is already available locally. |
| config.launcher.tag | string | `"0.0.2"` | Tag of the Docker image to use for launching pipeline jobs. Leave empty to use the current version of pipeManager. |
| config.launcher.timeout | int | `600` | Timeout in seconds for the pipeline job to complete. If the job exceeds this duration, it will be terminated. Adjust this value based on the expected duration of your pipeline jobs. |
| config.webhook | object | `{"routes":[],"workers":8}` | Configuration for Webhook Server |
| config.webhook.routes | list | `[]` | Routes configuration for the webhook server. This section defines the routes that the webhook server will listen to and how to handle incoming webhook payloads. CEL (Common Expression Language) expressions are used to extract and evaluate data from the JSON payload of webhooks. Ensure that the expressions are correctly formatted based on your webhook payload structure. Example route configuration for GitHub webhooks    - name: github      path: /github  # The endpoint path where GitHub will send webhook payloads. Ensure this path is accessible      # and correctly configured in your GitHub repository settings.      eventType: "data.headers['X-Github-Event'][0]"  # CEL expression to retrieve the event.       events:        # Handler for push events (branches and tags)        - type: push          repository: "data.body.payload.repository.ssh_url"  # (Mandatory) CEL expression to extract the repository name from the payload.          commit: "data.body.payload.head_commit.id"  # (Optional) CEL expression to extract the full commit SHA.          variables: # Additional custom variables to be passed to the pipeline.            ref: "data.body.payload.ref.replace('refs/heads/', '').replace('refs/tags/', '')" # CEL expression to extract the branch or tag name.            tag: "data.body.payload.ref.startsWith('refs/tags/')"  # CEL expression to determine if the push includes a tag (true/false).            shortCommit: "data.body.payload.head_commit.id.substring(0,7)"  # CEL expression to get the first 7 characters of the commit SHA.            email: "data.body.payload.head_commit.author.email"  # CEL expression to extract the author's email from the commit.            author: "data.body.payload.head_commit.author.name"  # CEL expression to extract the author's name from the commit.            user: "data.body.payload.pusher.name"  # CEL expression to extract the pusher's username.            custom: "'MY_CUSTOM_VALUE'"  # Custom value to be passed to the pipeline. Literal values must be enclosed in single quotes. |
| config.webhook.workers | int | `8` | Number of concurrent workers processing incoming webhook requests. Adjust this number based on your server's CPU cores and expected load. |
| fullnameOverride | string | `""` | Full name to override the default resource fullname. |
| image | object | `{"pullPolicy":"IfNotPresent","repository":"k3d-registry:5111/webhook-listener","tag":""}` | Image configuration for the webhook listener. |
| image.pullPolicy | string | `"IfNotPresent"` | Image Pull Policy for the webhook listener image. |
| image.repository | string | `"k3d-registry:5111/webhook-listener"` | Docker image to use for the webhook listener. TODO: change this to the final image name |
| image.tag | string | `""` | Tag of the Docker image to use for the webhook listener. Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` | Image Pull Secrets for the webhook listener. |
| ingress | object | `{"annotations":{},"className":"","enabled":false,"hosts":[{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}],"tls":[]}` | Configuration for the ingress |
| ingress.annotations | object | `{}` | Annotations to add to the ingress |
| ingress.className | string | `""` | Ingress class to use |
| ingress.enabled | bool | `false` | Enable Ingress |
| ingress.hosts | list | `[{"host":"chart-example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}]` | Host configuration for the ingress |
| ingress.tls | list | `[]` | TLS configuration for the ingress |
| livenessProbe | object | `{"httpGet":{"path":"/healthz","port":"http"}}` | Configuration for the liveness probes. |
| nameOverride | string | `""` | Name to override the default resource name. |
| nodeSelector | object | `{}` | Configuration for the node selector. |
| podAnnotations | object | `{}` | Configuration for the pod. Annotations. |
| podLabels | object | `{}` | Configuration for the pod. Labels. |
| podSecurityContext | object | `{}` | Configuration for the pod. Security Context. |
| readinessProbe | object | `{"httpGet":{"path":"/healthz","port":"http"}}` | Configuration for the readiness probes. |
| replicaCount | int | `1` | Configuration for the number of replicas to run. |
| resources | object | `{}` | Configuration for the resources |
| securityContext | object | `{}` | Configuration for the container. Security Context. |
| service | object | `{"port":80,"type":"ClusterIP"}` | Configuration for the service |
| service.port | int | `80` | Port to expose the service on |
| service.type | string | `"ClusterIP"` | Type of service to create |
| serviceAccount | object | `{"annotations":{},"automount":true,"create":true,"name":""}` | Configuration for the service account. |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.automount | bool | `true` | Automatically mount a ServiceAccount's API credentials? |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | Configuration for the tolerations. |
| volumeMounts | list | `[]` | Additional volumeMounts on the output Deployment definition. |
| volumes | list | `[]` | Additional volumes on the output Deployment definition. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
