# Argo Server SSO

> v2.9 and after

It is possible to use [Dex](https://github.com/dexidp/dex) for authentication. [This document](argo-server-sso-argocd.md) describes how to set up Argo Workflows and Argo CD so that Argo Workflows uses Argo CD's Dex server for authentication.

## To start Argo Server with SSO

Firstly, configure the settings [workflow-controller-configmap.yaml](workflow-controller-configmap.yaml) with the correct OAuth 2 values. If working towards an OIDC configuration the Argo CD project has [guides](https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/) on its similar (though different) process for setting up OIDC providers. It also includes examples for specific providers. The main difference is that the Argo CD docs mention that their callback address endpoint is `/auth/callback`.  For Argo Workflows, the default format is `/oauth2/callback` as shown in [this comment](https://github.com/argoproj/argo-workflows/blob/93c11a24ff06049c2197149acd787f702e5c1f9b/docs/workflow-controller-configmap.yaml#L329) in the default `values.yaml` file in the helm chart.

Next, create the Kubernetes secrets for holding the OAuth2 `client-id` and `client-secret`. You may refer to the kubernetes documentation on [Managing secrets](https://kubernetes.io/docs/tasks/configmap-secret/). For example by using `kubectl` with literals:

```bash
kubectl create secret -n argo generic client-id-secret \
  --from-literal=client-id-key=foo

kubectl create secret -n argo generic client-secret-secret \
  --from-literal=client-secret-key=bar
```

Then, start the Argo Server using the SSO [auth mode](argo-server-auth-mode.md):

```bash
argo server --auth-mode sso --auth-mode ...
```

## Token Revocation

> v2.12 and after

As of v2.12 we issue a JWE token for users rather than give them the ID token from your OAuth2 provider. This token is opaque and has a longer expiry time (10h by default).

The token encryption key is automatically generated by the Argo Server and stored in a Kubernetes secret name `sso`.

You can revoke all tokens by deleting the encryption key and restarting the Argo Server (so it generates a new key).

```bash
kubectl delete secret sso
```

!!! Warning
    The old key will be in the memory the any running Argo Server, and they will therefore accept and user with token encrypted using the old key. Every Argo Server MUST be restarted.

All users will need to log in again. Sorry.

## SSO RBAC

> v2.12 and after

You can optionally add RBAC to SSO. This allows you to give different users different access levels. Except for `client` auth mode, all users of the Argo Server must ultimately use a service account. So we allow you to define rules that map a user (maybe using their OIDC groups) to a service account in the same namespace as argo server by annotating the service account.

To allow service accounts to manage resources in other namespaces create a role and role binding in the target namespace.

RBAC config is installation-level, so any changes will need to be made by the team that installed Argo. Many complex rules will be burdensome on that team.

Firstly, enable the `rbac:` setting in [workflow-controller-configmap.yaml](workflow-controller-configmap.yaml). You likely want to configure RBAC using groups, so add `scopes:` to the SSO settings:

```yaml
sso:
  # ...
  scopes:
   - groups
  rbac:
    enabled: true
```

!!! Note
    Not all OIDC providers support the `groups` scope. Please speak to your provider about their options.

To configure a service account to be used, annotate it:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  annotations:
    # The rule is an expression used to determine if this service account
    # should be used.
    # * `groups` - an array of the OIDC groups
    # * `iss` - the issuer ("argo-server")
    # * `sub` - the subject (typically the username)
    # Must evaluate to a boolean.
    # If you want an account to be the default to use, this rule can be "true".
    # Details of the expression language are available in
    # https://expr-lang.org/docs/language-definition.
    workflows.argoproj.io/rbac-rule: "'admin' in groups"
    # The precedence is used to determine which service account to use when
    # Precedence is an integer. It may be negative. If omitted, it defaults to "0".
    # Numerically higher values have higher precedence (not lower, which maybe
    # counter-intuitive to you).
    # If two rules match and have the same precedence, then which one used will
    # be arbitrary.
    workflows.argoproj.io/rbac-rule-precedence: "1"
```

If no rule matches, we deny the user access.

Tip: You'll probably want to configure a default account to use if no other rule matches, e.g. a read-only account, you can do this as follows:

```yaml
metadata:
  name: read-only
  annotations:
    workflows.argoproj.io/rbac-rule: "true"
    workflows.argoproj.io/rbac-rule-precedence: "0"
```

The precedence must be the lowest of all your service accounts.

As of Kubernetes v1.24, secrets for a service account token are no longer automatically created.
Therefore, service account secrets for SSO RBAC must be created manually.
See [Service Account Secrets](service-account-secrets.md) for detailed instructions.

## SSO RBAC Namespace Delegation

> v3.3 and after

You can optionally configure RBAC SSO per namespace.
Typically, an organization has a Kubernetes cluster and a central team (the owner of the cluster) manages the cluster. Along with this, there are multiple namespaces which are owned by individual teams. This feature would help namespace owners to define RBAC for their own namespace.

The feature is currently in beta.
To enable the feature, set env variable `SSO_DELEGATE_RBAC_TO_NAMESPACE=true` in your argo-server deployment.

### Recommended usage

Configure a default account in the installation namespace that allows access to all users of your organization. This service account allows a user to login to the cluster. You could optionally add a workflow read-only role and role-binding.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: user-default-login
  annotations:
    workflows.argoproj.io/rbac-rule: "true"
    workflows.argoproj.io/rbac-rule-precedence: "0"
```

!!! Note
    All users MUST map to a cluster service account (such as the one above) before a namespace service account can apply.

Now, for the namespace that you own, configure a service account that allows members of your team to perform operations in your namespace.
Make sure that the precedence of the namespace service account is higher than the precedence of the login service account. Create an appropriate role for this service account and bind it with a role-binding.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-namespace-read-write-user
  namespace: my-namespace
  annotations:
    workflows.argoproj.io/rbac-rule: "'my-team' in groups"
    workflows.argoproj.io/rbac-rule-precedence: "1"
```

With this configuration, when a user is logged in via SSO, makes a request in `my-namespace`, and the `rbac-rule` matches, this service account allows the user to perform that operation. If no service account matches in the namespace, the first service account (`user-default-login`) and its associated role will be used to perform the operation.

## SSO Login Time

> v2.12 and after

By default, your SSO session will expire after 10 hours. You can change this by adding a `sessionExpiry` to your [workflow-controller-configmap.yaml](workflow-controller-configmap.yaml) under the SSO heading.

```yaml
sso:
  # Expiry defines how long your login is valid for in hours. (optional)
  sessionExpiry: 240h
```

## Custom claims

> v3.1.4 and after

If your OIDC provider provides groups information with a claim name other than `groups`, you could configure config-map to specify custom claim name for groups. Argo now arbitrary custom claims and any claim can be used for `expr eval`. However, since group information is displayed in UI, it still needs to be an array of strings with group names as elements.

The `customClaim` in this case will be mapped to `groups` key and we can use the same key `groups` for evaluating our expressions

```yaml
sso:
  # Specify custom claim name for OIDC groups.
  customGroupClaimName: argo_groups
```

If your OIDC provider provides groups information only using the user-info endpoint (e.g. Okta), you could configure `userInfoPath` to specify the user info endpoint that contains the groups claim.

```yaml
sso:
  userInfoPath: /oauth2/v1/userinfo
```

### Example Expression

```bash
# assuming customClaimGroupName: argo_groups
workflows.argoproj.io/rbac-rule: "'argo_admins' in groups"
```

## Filtering groups

> v3.5 and after

You can configure `filterGroupsRegex` to filter the groups returned by the OIDC provider. Some use-cases for this include:

- You have multiple applications using the same OIDC provider, and you only want to use groups that are relevant to Argo Workflows.
- You have many groups and exceed the [4KB cookie size limit](https://chromestatus.com/feature/4946713618939904) (cookies are used to store authentication tokens). If this occurs, login will fail.

```yaml
sso:
    # Specify a list of regular expressions to filter the groups returned by the OIDC provider.
    # A logical "OR" is used between each regex in the list
    filterGroupsRegex:
    - ".*argo-wf.*"
    - ".*argo-workflow.*"
```
