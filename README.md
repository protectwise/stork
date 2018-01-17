# Stork
*Answering the question, where do tokens come from?*

Stork is a small utility designed to retrieve tokens from [Hashicorp Vault](https://www.vaultproject.io) for instances running on AWS EC2. If you have an EC2 instance with an IAM profile, you can use Stork to retrieve a token from Vault for you. (via the [Vault EC2 authentication method](https://www.vaultproject.io/docs/auth/aws.html))

## Authentication Workflow

More detailed documentation is available in the [Vault docs](https://www.vaultproject.io/docs/auth/aws.html) but the extremely short version is that Vault is capable of using AWS' EC2 metadata service and IAM profiles to authenticate EC2 clients. This allows you to get secrets to instances without storing any secret data in the AMI, user-data or elsewhere.

Stork is a simple program intended to run on EC2 clients. Stork completes all the steps of the authentication workflow and writes both a permanent nonce (which is meant to be accessible only to Stork) and a temporary token (accessible to whatever application needs to interact with Vault) to files on disk.

You can periodically run Stork to replace expiring or expired tokens, perhaps as a cron job.

## Getting Started

On the Vault server side, you will need to [enable the AWS authentication backend](https://www.vaultproject.io/docs/auth/aws.html#enable-aws-ec2-authentication-in-vault-), configure your Vault servers to have permission to [query EC2](https://www.vaultproject.io/docs/auth/aws.html#recommended-vault-iam-policy), and you will also need to set up a [policy for the IAM role of your EC2 instance](https://www.vaultproject.io/docs/auth/aws.html#configure-the-policies-on-the-role-).

Install Stork on the EC2 client instance (the instance you want to have a token created on):

```go get github.com/protectwise/stork```

Ensure your client instance has an IAM instance profile. Your instance does not need any permissions to any AWS resources, but it does need an instance profile as this is how the Vault server authorizes your client with the right Vault policies.

Once everything is set up in IAM, EC2 and the Vault server, you can run Stork to get a token from your Vault server:

```vault-stork login --server https://vault.internal.yourcompany.com --token /etc/stork/token --nonce /etc/stork/nonce```

(Note this example assumes you've already created the directory `/etc/stork`)

If everything works, Stork will exit with status code 0 and `/etc/stork/token` will be the token that the Vault server gave us.

## Going Further
Stork only gets tokens for you, retrieving secrets and interacting with Vault is up to your application.

For example, you can use the Vault cli to retrieve secrets from your Vault server with this token:

```
export $VAULT_TOKEN=$(cat /etc/stork/token)
vault read secret/test/super_duper_secret
Key             	Value
---             	-----
refresh_interval	768h0m0s
test				I am sekrit!
```

## Frequently Asked Questions

### What's the deal with the nonce?
For a detailed discussion, refer to the [Vault docs](https://www.vaultproject.io/docs/auth/aws.html#client-nonce). By default anything with access to make network requests on your EC2 client can query the EC2 metadata service, which means that anything that could make a network request could make a request to the EC2 metadata service and give that to the Vault server and receive a token.

The nonce ensures that whoever authenticates to Vault first wins. If an attacker tries to impersonate your instance after Stork has run, it would also need access to the nonce to succeed (so keep the nonce safe!). If an attacker beats Stork and receives a token before Stork runs, Stork will return an error from Vault about a client nonce mismatch. This is your opportunity to sound an alarm! Either way, the nonce provides an additional layer of security.

We suggest running Stork as a different user, and using standard Unix permissions to ensure that this user is the only one on client systems with access to read the `nonce` file. The `token` file that gets written should also only be readable to your application in a similar manner.

### I updated my Vault policy but my token is still being denied!
This was one of the most common misunderstandings when we started implementing Vault. It is not EC2 or Stork specific but it applies to all tokens in Vault. Once created, the policies that a token has are **immutable**. You will need a new token with new policies. [Revoke the old token](https://www.vaultproject.io/intro/getting-started/authentication.html#tokens), and create a new one. (re-run Stork!) You can use `vault token-lookup $TOKEN` to see what policies apply to any token.

