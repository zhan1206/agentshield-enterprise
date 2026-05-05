module github.com/zhan1206/agentshield-enterprise/backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/redis/go-redis/v9 v9.5.1
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/trace v1.24.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.6.0
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.32.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/gorm v1.25.7
	gorm.io/driver/mysql v1.5.4
	github.com/segmentio/kafka-go v0.4.47
)