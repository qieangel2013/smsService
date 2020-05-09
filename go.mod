module smsService

go 1.13

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/foolin/goview v0.3.0
	github.com/gin-gonic/gin v1.6.2
	github.com/go-openapi/spec v0.19.7 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/golang/protobuf v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/viper v1.6.3
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	github.com/urfave/cli v1.22.4 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	golang.org/x/tools v0.0.0-20200426102838-f3a5411a4c3b // indirect
	julive v0.0.0-00010101000000-000000000000
)

replace julive => ./lib
