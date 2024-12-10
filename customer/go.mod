module github.com/fyerfyer/trade-refactor/customer

go 1.23.1

require (
	github.com/fyerfyer/trade-dependency/dto v0.0.15
	github.com/fyerfyer/trade-dependency/pkg/cache v0.0.3
	github.com/fyerfyer/trade-dependency/pkg/e v0.0.1
	github.com/fyerfyer/trade-dependency/proto/grpc/customer v0.0.7
	github.com/fyerfyer/trade-dependency/proto/grpc/order v0.0.7
	google.golang.org/grpc v1.68.1
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)
