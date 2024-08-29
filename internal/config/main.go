package config

import (
	"log"

	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
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
	viper  *viper.Viper
}

func New(getter kv.Getter) Config {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return &config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		viper:      v,
	}
}

func (c *config) ContractAddress() string {
	return c.viper.GetString("ethereum.contract_address")
}

func (c *config) RPCURL() string {
	return c.viper.GetString("ethereum.rpc_url")
}
