# vault-go
Go module to access hashicorp vault secrets by application role

## Installation

### Setting Up Go
To install Go, visit [this link](https://golang.org/dl/).

### Installing Module
`go get -u github.com/kangchengkun/vault-go`

## Usage
Before using this Go module, you will need to [install/setup your own vault server](https://www.vaultproject.io/docs/install).

Get/Setup the [application role](https://www.vaultproject.io/docs/auth/approle) from your vault server

```
import github.com/kangchengkun/vault-go

vault.BaseUrl = "your-vault-api-base-url"
vault.AuthUrl = "your-vault-api-auth-url"
vault.RoleID = "your-vault-application-role-id"
vault.SecretID = "your-vault-application-secret-id"



// Init the vault token first
err := vault.Login()
if err != nil {
    fmt.Println("Login to vault failed")
}

// Read secrets from vault
output, err := vault.ReadData(dataPath)
if err != nil {
    fmt.Printf("Read data from vault failed with error: %v\n", err)
}

```

## Contribution

Follow the [Guide](https://go.dev/blog/publishing-go-modules) to publish new versions

```
...
git add .
git commit -m "new updates"

$ git tag vx.x.x
$ git push origin vx.x.x
```