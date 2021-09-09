module github.com/ONSdigital/dp-identity-api

go 1.16

// update to etcd 3.3.23 to fix security issues as reported by make audit
// note we can't use this version directly because the vendor mod is broken
// see https://github.com/etcd-io/etcd/issues/11154

replace github.com/coreos/etcd => go.etcd.io/etcd v0.0.0-20200716221548-4873f5516bd9

require (
	github.com/ONSdigital/dp-component-test v0.3.0
	github.com/ONSdigital/dp-healthcheck v1.1.0
	github.com/ONSdigital/dp-net v1.2.0
	github.com/ONSdigital/log.go/v2 v2.0.9
	github.com/aws/aws-sdk-go v1.38.34
	github.com/cucumber/godog v0.11.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-memdb v1.3.2 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/maxcnunes/httpfake v1.2.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.2.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/afero v1.6.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.5.2 // indirect
	golang.org/x/crypto v0.0.0-20210505212654-3497b51f5e64 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
