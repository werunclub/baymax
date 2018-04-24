package helpers

import (
	"context"
	"fmt"
	"hash/fnv"
	"net"

	"github.com/satori/go.uuid"
	"github.com/smallnest/rpcx/share"
)

const (
	RPCPath = "/rpc"
)

var (
	privateBlocks []*net.IPNet
)

func init() {
	for _, b := range []string{"10.0.0.0/8", "100.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

// ExtractAddress 解析网络地址
func ExtractAddress(addr string) (string, error) {
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]") {
		return addr, nil
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("Failed to get interface addresses! Err: %v", err)
	}

	var ipAddr []byte

	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}

		if !IsPrivateIP(ip.String()) {
			continue
		}

		ipAddr = ip
		break
	}

	if ipAddr == nil {
		return "", fmt.Errorf("No private IP address found, and explicit IP not provided")
	}

	return net.IP(ipAddr).String(), nil
}

// IsPrivateIP 是否为私有IP
func IsPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

// Hash consistently chooses a hash bucket number in the range [0, numBuckets) for the given key. numBuckets must be >= 1.
func Hash(key uint64, buckets int32) int32 {
	if buckets <= 0 {
		buckets = 1
	}

	var b, j int64

	for j < int64(buckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return int32(b)
}

// HashString get a hash value of a string
func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// HashServiceAndArgs define a hash function
type HashServiceAndArgs func(len int, options ...interface{}) int

// JumpConsistentHash selects a server by serviceMethod and args
func JumpConsistentHash(len int, options ...interface{}) int {
	keyString := ""
	for _, opt := range options {
		keyString = keyString + "/" + ToString(opt)
	}
	key := HashString(keyString)
	return int(Hash(key, int32(len)))
}

// ToString toString
func ToString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}

// NewMetaDataContext gen context for metadata
func NewMetaDataContext(req map[string]string) context.Context {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, req)
	return context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))
}

// MetaData 元数据
type MetaData struct {
	ctx context.Context
	req map[string]string
	res map[string]string
}

// NewMetaDataFormContext 上下文生成元数据
func NewMetaDataFormContext(ctx context.Context) MetaData {
	return MetaData{
		ctx: ctx,
		req: ctx.Value(share.ReqMetaDataKey).(map[string]string),
		res: ctx.Value(share.ResMetaDataKey).(map[string]string),
	}
}

// Request 请求元数据
func (m MetaData) Request() map[string]string {
	return m.req
}

// Response 响应元数据
func (m MetaData) Response() map[string]string {
	return m.res
}

// Get 获取元数据
func (m MetaData) Get(key string) string {
	return m.req[key]
}

// Set 添加元数据
func (m MetaData) Set(key, val string) error {
	m.res[key] = val
	return nil
}

// GetRequestID 获取 requestID
func GetRequestID(ctx context.Context) string {
	meta := NewMetaDataFormContext(ctx)
	requestID := meta.Get("request_id")

	if requestID == "" {
		requestID = uuid.NewV4().String()
	}

	return requestID
}
