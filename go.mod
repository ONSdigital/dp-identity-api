module github.com/ONSdigital/dp-identity-api/v2

go 1.24

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.261.0
	github.com/ONSdigital/dp-authorisation/v2 v2.31.2
	github.com/ONSdigital/dp-component-test v0.15.0
	github.com/ONSdigital/dp-healthcheck v1.6.3
	github.com/ONSdigital/dp-net/v2 v2.12.0
	github.com/ONSdigital/dp-permissions-api v0.25.0
	github.com/ONSdigital/log.go/v2 v2.4.3
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2/config v1.29.10
	github.com/aws/aws-sdk-go-v2/service/s3 v1.78.2
	github.com/aws/smithy-go v1.22.3
	github.com/cucumber/godog v0.15.0
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.3.1
	github.com/smartystreets/goconvey v1.8.1
	github.com/stretchr/testify v1.9.0
)

require (
	//github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.51.3
)

require (
	github.com/ONSdigital/dp-mongodb-in-memory v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.63 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.29.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.17 // indirect
	github.com/chromedp/cdproto v0.0.0-20241110205750-a72e6703cd9b // indirect
	github.com/chromedp/chromedp v0.11.2 // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/justinas/alice v1.2.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/maxcnunes/httpfake v1.2.4 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.17.1 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
