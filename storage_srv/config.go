package main

var Config = struct {
	ServiceName string `default:"storage_srv"`

	Server struct {
		Addr string `default:":8080"`
	}

	// 阿里云 OSS 配置
	Storage struct {
		AccessKey    string
		AccessSecret string
		Endpoint     string
		BucketName   string
		MaxSize      int64 `default: 5368709120`
	}

	Url struct {
		BaseUrl  string
		Suffixes []string
	}

	Registry struct {
		Type    string
		Address string
	}
}{}
