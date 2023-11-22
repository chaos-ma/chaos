package direct

/**
* created by mengqi on 2023/11/21
 */

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

func init() {
	resolver.Register(NewBuilder())
}

type DirectBuilder struct{}

// NewBuilder creates a directBuilder which is used to factory direct resolvers.
// example:
//
//	direct://<authority>/127.0.0.1:9000,
func NewBuilder() *DirectBuilder {
	return &DirectBuilder{}
}

func (d *DirectBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}

	//grpc建立连接的逻辑都是这里UpdateState
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

func (d *DirectBuilder) Scheme() string {
	return "direct"
}

var _ resolver.Builder = &DirectBuilder{}
