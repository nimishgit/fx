package main

import (
	"sync"

	"github.com/uber-go/uberfx/core"
	"github.com/uber-go/uberfx/examples/keyvalue/kv"
	"github.com/yarpc/yarpc-go/encoding/thrift"
)

type YarpcHandler struct {
	sync.RWMutex

	items map[string]string
}

func NewYarpcThriftHandler(svc *core.Service) (thrift.Service, error) {
	return kv.New(&YarpcHandler{items: map[string]string{}}), nil
}

func (h *YarpcHandler) GetValue(req *thrift.ReqMeta, key *string) (string, *thrift.ResMeta, error) {
	h.RLock()
	defer h.RUnlock()

	if value, ok := h.items[*key]; ok {
		return value, nil, nil
	}

	return "", nil, &kv.ResourceDoesNotExist{Key: *key}
}

func (h *YarpcHandler) SetValue(req *thrift.ReqMeta, key *string, value *string) (*thrift.ResMeta, error) {
	h.Lock()

	h.items[*key] = *value
	h.Unlock()
	return nil, nil
}