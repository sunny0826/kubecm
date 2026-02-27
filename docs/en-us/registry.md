
# Registry

`kubecm registry` is a Git-backed kubeconfig distribution system for team-based cluster management.

In multi-team environments, each team needs a different subset of Kubernetes clusters with different access levels. Instead of manually running `kubecm cloud add` for every cluster, a central Git repository declares clusters organized by **role**, with optional **user** definitions for shared credentials. At sync time, kubecm calls the cloud provider API to fetch the full kubeconfig — **no secrets are stored in Git**.

## How it works

1. A Git repository contains a `registry.yaml`, roles, clusters, and optional users
2. Each team member runs `kubecm registry add` to clone and sync
3. `kubecm registry sync` pulls the latest changes and updates the kubeconfig
4. Stale contexts are automatically removed when clusters are removed from a role

## Repository structure

```
my-kubeconfig-registry/
  registry.yaml             # metadata + template variables
  roles/
    devops.yaml             # role = list of clusters + optional user bindings
    backend.yaml
    readonly.yaml
  clusters/
    eks-prod-eu.yaml        # AWS EKS cluster
    eks-staging-eu.yaml
    aks-prod.yaml           # Azure AKS cluster
    onprem-dc1.yaml         # static (on-prem) cluster
  users/                    # optional: reusable credential definitions
    admin.yaml
    readonly.yaml
```

## File templates

### registry.yaml

The root configuration file. Defines the registry name and template variables that can be used in clusters via Go templates (`{{ .VariableName }}`).

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Registry
metadata:
  name: my-company
  description: "My company Kubernetes clusters"
variables:
  - name: Username
    description: "Your corporate username (e.g. john.doe)"
    required: true
  - name: Environment
    description: "Target environment"
    required: false
    default: "prod"
```

### Role file (roles/devops.yaml)

A role defines which clusters are available to a team. The optional `contextPrefix` is prepended to context names in the kubeconfig.

#### Simple format

The simplest role lists clusters directly:

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: devops
  description: "Full access to all clusters"
contextPrefix: "mycompany"    # optional: contexts become mycompany-<cluster>
contexts:
  - cluster: eks-prod-eu
  - cluster: eks-staging-eu
  - cluster: aks-prod
  - cluster: onprem-dc1
```

> **Tip**: If your cluster names are already descriptive enough, omit `contextPrefix` to use cluster names directly as context names.

#### With user overrides

When the same cluster needs multiple identities (e.g. admin + readonly), bind a **user** to each entry and set a unique `name`:

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Role
metadata:
  name: devops
contexts:
  # Cluster with user override
  - cluster: eks-prod-eu
    user: admin
    name: eks-prod-admin       # required when same cluster appears twice
  - cluster: eks-prod-eu
    user: readonly
    name: eks-prod-ro
  # Cluster without user (credentials from cluster definition)
  - cluster: onprem-dc1
```

> **Note**: When the same cluster appears more than once, you **must** provide a `name` field to disambiguate context names.

### User file (users/admin.yaml)

A user defines reusable credentials that can be bound to any cluster. The user's provider must match the cluster's provider.

```yaml
apiVersion: kubecm.io/v1alpha1
kind: User
metadata:
  name: admin
provider: aws
aws:
  profile: "my-account/AWSAdministratorAccess"
```

```yaml
apiVersion: kubecm.io/v1alpha1
kind: User
metadata:
  name: azure-admin
provider: azure
azure:
  tenantId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

When a user is bound to a cluster, the user's credentials override the cluster's. This allows separating cluster topology (region, cluster name) from authentication (profile, tenant).

Template variables are supported in user fields:

```yaml
apiVersion: kubecm.io/v1alpha1
kind: User
metadata:
  name: team-admin
provider: aws
aws:
  profile: "{{ .AWSProfile }}"
```

### Cluster: AWS EKS (clusters/eks-prod-eu.yaml)

For AWS EKS clusters. At sync, kubecm calls `aws eks describe-cluster` using the specified profile and region. When a user is bound, the user's profile takes precedence.

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Cluster
metadata:
  name: eks-prod-eu
provider: aws
aws:
  region: eu-central-1
  cluster: my-eks-cluster
  profile: "my-account/AWSAdministratorAccess"   # can be omitted if using a user
```

Template variables are supported in all fields:

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Cluster
metadata:
  name: eks-dev
provider: aws
aws:
  region: "{{ .Region }}"
  cluster: "eks-{{ .Environment }}"
  profile: "{{ .AWSProfile }}"
```

### Cluster: Azure AKS (clusters/aks-prod.yaml)

For Azure AKS clusters. At sync, kubecm fetches credentials via the Azure SDK.

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Cluster
metadata:
  name: aks-prod
provider: azure
azure:
  subscriptionId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  resourceGroup: "rg-prod"
  cluster: "aks-prod"
  tenantId: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"    # can be omitted if using a user
```

### Cluster: Static / On-prem (clusters/onprem-dc1.yaml)

For clusters without a supported cloud provider. The full kubeconfig is embedded. Go template variables are supported.

```yaml
apiVersion: kubecm.io/v1alpha1
kind: Cluster
metadata:
  name: onprem-dc1
provider: static
kubeconfig: |
  apiVersion: v1
  kind: Config
  clusters:
    - cluster:
        server: https://k8s.internal:6443
        certificate-authority-data: LS0tLS1C...
      name: onprem-dc1
  contexts:
    - context:
        cluster: onprem-dc1
        user: onprem-dc1
      name: onprem-dc1
  users:
    - user:
        token: "{{ .Username }}-token"
      name: onprem-dc1
```

## Usage

### Add a registry

```bash
# Add with inline variables
kubecm registry add --name mycompany \
  --url git@github.com:myorg/kubeconfig-registry.git \
  --role devops \
  --var Username=john.doe

# Add interactively (prompts for required variables)
kubecm registry add --name mycompany \
  --url git@github.com:myorg/kubeconfig-registry.git \
  --role devops
```

### Sync

```bash
# Sync a specific registry
kubecm registry sync mycompany

# Sync all registries
kubecm registry sync --all

# Preview changes without applying
kubecm registry sync mycompany --dry-run
```

### List registries

```bash
kubecm registry list
```

### Update settings

```bash
# Change role
kubecm registry update mycompany --role backend

# Update a variable
kubecm registry update mycompany --var Username=jane.doe

# Change branch
kubecm registry update mycompany --ref develop
```

### Remove a registry

```bash
# Remove registry and its managed contexts
kubecm registry remove mycompany

# Remove registry but keep kubeconfig contexts
kubecm registry remove mycompany --keep-contexts
```

## Local state

kubecm stores registry state in `~/.kubecm/`:

```
~/.kubecm/
  config.yaml                  # registry entries, variables, managed contexts
  registries/
    mycompany/                 # cloned Git repo
```

## Sync behavior

| Scenario | Action |
|----------|--------|
| New cluster in role | Context **added** to kubeconfig |
| Existing managed context | Context **updated** (registry takes authority) |
| Cluster removed from role | Context **removed** from kubeconfig |
| Context exists but not managed | **Skipped** with warning (never overwrites) |

