# cagen

This is a certificate generator. It generates a cert and pem for you instead of having to use openssl.

## Build

`go build .`

## Run
`./cagen -output-path <output path> -common-name <certificate name> -organization <organization> -organizational-unit <unit> -address <address> -locality <city> -province <state or province> -postal-code <postal code> -country <country>`

## Options

| Option | Default | Description |
|--------|----------|-------------|
| prefix | ca | The prefix to the final output (e.g. ca results in ca.pem and ca.crt) |
| output-path | ./ | The path that you want the files saved to |
| common-name |  | What you want the certificate to be named |
| organization |  | The organization that the certificate is for |
| organizational-unit |  | The unit or division that the certificate is for |
| address |  | The address of the organization or unit |
| locality |  | The city, town, etc. of the organization or unit |
| province |  |The state or province of the organization or unit |
| postal-code |  | The postal code of the organization or unit |
| country |  | The country of the organization or unit |
| key-size | 2048 | The key size of the certificate. Must be divisible by 1024 |
| days-to-expire | 365 | When the certificate expires. Must be greater than or equal to 90 |