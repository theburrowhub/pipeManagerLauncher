# (WARNING: First version of the README, still needs to be updated)

# PipeManager
**PipeManager** is a Tekton-based pipeline management tool designed to streamline the process of transforming webhooks into Kubernetes jobs that trigger Tekton PipelineRuns. The software enables management of pipelines running in your system, such as stopping, rerunning, and canceling them. It supports integration with GitHub, Bitbucket, GitLab, and other version control systems through webhooks, using a configuration file (YAML) to define the pipeline behavior.

## Key Features
* **Webhook Listener**: Listens for incoming webhooks from GitHub, Bitbucket, GitLab, etc.
* **JSON Transformation**: Converts incoming webhook JSON payloads into Kubernetes jobs using a YAML configuration with CEL (Common Expression Language) queries.
* **Kubernetes Job Creation**: Dynamically creates Kubernetes jobs to trigger containerized workflows.
* **Tekton Integration**: Transforms the Kubernetes job output into a Tekton PipelineRun using Go templates.
* **Pipeline Management**: Provides a management interface to stop, rerun, or cancel Tekton pipelines.
* **Local Development Environment**: Runs on K3s or Minikube with local Docker image registry and ngrok for webhook testing.

## Architecture Overview
 
1. **Webhook Listener**:

* The main application listens for webhooks (e.g., from GitHub, Bitbucket, GitLab).
* Based on the webhook payload, the application transforms the data into a Kubernetes job.

2. **YAML Configuration**:

* Uses a YAML file that defines CEL queries to parse and transform the webhook payload into actionable data.
* This data is used to create the Kubernetes job.

3. **Kubernetes Job**:

* The transformed webhook data is used to spin up a Kubernetes job that will run a container.
* This container is dynamically generated using Go templates.

4. **PipelineRun Creation**:

* The Kubernetes job executes the container, which will translate a bunch of files in the folder `pipelines` of the repository into a Tekton PipelineRun.
* These files will contain the definition of the pipeline to be executed by Tekton.
* The PipelineRun will trigger the pipeline defined in the repository's `pipelines` folder.

5. **Pipeline Management**:

* A management interface allows users to stop, rerun, and cancel pipelines running in the system.
* Offers real-time feedback on pipeline status and execution. 

## Development Environment Setup 
PipeManager can be developed and tested in a local environment using **K3s** or **Minikube**, along with **Docker** and **ngrok**.

### Prerequisites
* [K3s](https://k3s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/start/) (local Kubernetes clusters)
* [Docker](https://docs.docker.com/get-docker/) (containerization platform)
* [ngrok](https://ngrok.com/) (secure introspectable tunnels to localhost)
* [Tekton](https://tekton.dev/) (Kubernetes-native CI/CD pipelines)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) (Kubernetes command-line tool)
* [Go](https://golang.org/doc/install) (Go programming language)

### Step 1: Install K3s or Minikube
* Install K3s: Follow (K3s installation guide)[https://k3s.io/].
* Install Minikube: Follow the (Minikube installation guide)[https://minikube.sigs.k8s.io/docs/start/].

Ensure you can interact with your local cluster using `kubectl`:

```bash
kubectl get nodes
```

### Step 2: Set Up Docker Image Registry

PipeManager requires a local Docker image registry to store container images.

For Minikube:

```bash
minikube addons enable registry
```

For K3s, you can set up a local registry using [this guide](https://rancher.com/docs/k3s/latest/en/installation/private-registry/).

### Step 3: Install Tekton

Install Tekton Pipelines into your local Kubernetes cluster:

```bash
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```

### Step 4: Expose Webhook Listener Using ngrok
Use **ngrok** to expose the local webhook listener to the web.

Start ngrok:

```bash
ngrok http 8080
```

You’ll receive a forwarding URL (e.g., https://12345.ngrok.io). Use this URL to configure your GitHub, GitLab, or Bitbucket webhook.

### Step 5: Run the Application Locally

Clone the repository:

```bash
git clone https://github.com/your-repo/pipe-manager.git
cd pipe-manager
```

Build the application:

```bash
go build -o pipe-manager
```

Start the application:

```bash
./pipe-manager
```

The application will listen for incoming webhooks on port 8080 by default. You can change the port in the configuration.

### Step 6: Configure Your YAML Files

Configure the transformation of the webhook payload in the pipeline-config.yaml file:

```yaml
query:
  - name: "extract-repo"
    expression: "json.payload.repository.full_name"
  - name: "extract-branch"
    expression: "json.payload.ref"
  - name: "extract-commit"
    expression: "json.payload.head_commit.id"

jobTemplate: |
  apiVersion: batch/v1
  kind: Job
  metadata:
    name: {{ .Name }}
  spec:
    template:
      spec:
        containers:
        - name: pipeline-runner
          image: my-repo/pipeline-image:latest
          env:
            - name: REPO
              value: "{{ .Repo }}"
            - name: BRANCH
              value: "{{ .Branch }}"
            - name: COMMIT
              value: "{{ .Commit }}"
```

### Step 7: Test the Webhook Integration

Once the webhook listener is running, trigger a test webhook from your version control system. Ensure that the JSON payload is correctly transformed into a Kubernetes job and that the job spawns a `PipelineRun` in Tekton.

### Step 8: Managing Pipelines
The management interface allows you to control the state of the pipelines running on the system:

* **Stop a Pipeline**: Stop a currently running pipeline.
* **Rerun a Pipeline**: Rerun a failed or completed pipeline.
* **Cancel a Pipelin**e: Cancel an ongoing pipeline execution.

### Troubleshooting
* **Webhook not triggering**: Check if ngrok is correctly forwarding the requests to the local system.
* **PipelineRun not starting**: Verify the Kubernetes job configuration and ensure the correct Tekton pipelines are installed.
* **Image pull issues**: Ensure the local Docker image registry is properly set up and configured.

## Contributing
Feel free to submit issues or pull requests. We welcome contributions!