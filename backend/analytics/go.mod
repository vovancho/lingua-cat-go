module github.com/vovancho/lingua-cat-go/analytics

go 1.24.2

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.34.0
	github.com/ThreeDotsLabs/watermill v1.4.6
	github.com/ThreeDotsLabs/watermill-kafka/v3 v3.0.6
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.26.0
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	github.com/vovancho/lingua-cat-go/pkg/auth v0.0.0
	github.com/vovancho/lingua-cat-go/pkg/db v0.0.0-00010101000000-000000000000
	github.com/vovancho/lingua-cat-go/pkg/error v0.0.0
	github.com/vovancho/lingua-cat-go/pkg/request v0.0.0-00010101000000-000000000000
	github.com/vovancho/lingua-cat-go/pkg/response v0.0.0
	github.com/vovancho/lingua-cat-go/pkg/tracing v0.0.0-00010101000000-000000000000
	github.com/vovancho/lingua-cat-go/pkg/translator v0.0.0-00010101000000-000000000000
	github.com/vovancho/lingua-cat-go/pkg/validator v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/sdk v1.35.0
)

replace (
	github.com/vovancho/lingua-cat-go/pkg/auth => ../pkg/auth
	github.com/vovancho/lingua-cat-go/pkg/db => ../pkg/db
	github.com/vovancho/lingua-cat-go/pkg/error => ../pkg/error
	github.com/vovancho/lingua-cat-go/pkg/request => ../pkg/request
	github.com/vovancho/lingua-cat-go/pkg/response => ../pkg/response
	github.com/vovancho/lingua-cat-go/pkg/tracing => ../pkg/tracing
	github.com/vovancho/lingua-cat-go/pkg/translator => ../pkg/translator
	github.com/vovancho/lingua-cat-go/pkg/validator => ../pkg/validator
)

require (
	github.com/ClickHouse/ch-go v0.65.1 // indirect
	github.com/IBM/sarama v1.43.3 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dnwe/otelsarama v0.0.0-20240308230250-9388d9d40bc0 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.31 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lithammer/shortuuid/v3 v3.0.7 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.72.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
