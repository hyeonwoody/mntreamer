module mntreamer/monolithic/cmd

go 1.23.5

replace mntreamer/shared => ../../shared

replace mntreamer/monitor/cmd => ../../monitor/cmd

replace mntreamer/platform/cmd => ../../platform/cmd

replace mntreamer/streamer/cmd => ../../streamer/cmd

replace mntreamer/media/cmd => ../../media/cmd

require (
	github.com/gin-gonic/gin v1.10.0
	mntreamer/media/cmd v0.0.0-00010101000000-000000000000
	mntreamer/monitor/cmd v0.0.0-00010101000000-000000000000
	mntreamer/platform/cmd v0.0.0-00010101000000-000000000000
	mntreamer/shared v0.0.0
	mntreamer/streamer/cmd v0.0.0-00010101000000-000000000000
)

require (
	github.com/PuerkitoBio/goquery v1.10.2 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/neo4j/neo4j-go-driver/v5 v5.28.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/playwright-community/playwright-go v0.5001.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/gorm v1.25.12 // indirect
)
