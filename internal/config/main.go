package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer

	ContractAddress() string
	RPCURL() string
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter kv.Getter

	contractAddress string
	rpcURL          string
	once            comfig.Once
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}

func (c *config) loadConfig() {
	c.once.Do(func() interface{} {
		cfg := struct {
			ContractAddress string `fig:"contract_address,required"`
			RPCURL          string `fig:"rpc_url,required"`
		}{}

		rawConfig := kv.MustGetStringMap(c.getter, "ethereum")

		err := figure.
			Out(&cfg).
			From(rawConfig).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to load ethereum config"))
		}

		c.contractAddress = cfg.ContractAddress
		c.rpcURL = cfg.RPCURL

		return nil
	})
}

func (c *config) ContractAddress() string {
	c.loadConfig()
	return c.contractAddress
}

func (c *config) RPCURL() string {
	c.loadConfig()
	return c.rpcURL
}
