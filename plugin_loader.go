package main

import (
	"errors"
	"fmt"
	"plugin"

	"github.com/kffl/gocannon/common"
)

var (
	ErrPluginOpen      = errors.New("could not open plugin")
	ErrPluginLookup    = errors.New("could not lookup plugin interface")
	ErrPluginInterface = errors.New("module symbol doesn't match GocannonPlugin")
)

func loadPlugin(file string, silentOutput bool) (common.GocannonPlugin, error) {

	p, err := plugin.Open(file)
	if err != nil {
		return nil, ErrPluginOpen
	}

	pluginSymbol, err := p.Lookup("GocannonPlugin")
	if err != nil {
		return nil, ErrPluginLookup
	}

	var gocannonPlugin common.GocannonPlugin
	gocannonPlugin, ok := pluginSymbol.(common.GocannonPlugin)
	if !ok {
		return nil, ErrPluginInterface
	}

	if !silentOutput {
		fmt.Printf("Plugin %s loaded.\n", gocannonPlugin.GetName())
	}

	return gocannonPlugin, nil
}
