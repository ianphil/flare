# Flare
Is a command line utility that runs newman tests against an API. It does this by createing multiple containers in Azure Container Instances (ACI). One executes the requests contained in a Postman collection file. The other serves the report generated.

This can help you test your web app/api from any Azure Data Center where ACI instances will run.

## To use Flare:
Clone the repo

`git clone https://github.com/iphilpot/flare.git`

Build the project

`go build -o flare main.go`

Run flare

`./flare --help`

This project is in the very early stages of development. If you have questions please post an issue.

## Contributing
Please see our [CONTRIBUTING](CONTRIBUTING.md) guide.

## Code of Conduct
Please see our [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md).
