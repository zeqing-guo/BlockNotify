package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zeqing-guo/BlockNotify/blocknotify"
)

func main() {
	app := cli.NewApp()
	app.Name = "blocknotify"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "redis_url", Value: "localhost:6379", Usage: "redis url", EnvVar: "REDIS_URL"},
		cli.StringFlag{Name: "endpoint_url", Value: "localhost:8545", Usage: "ethereum endpoint url", EnvVar: "ENDPOINT_URL"},
		cli.StringFlag{Name: "log_level", Value: "debug", Usage: "log level", EnvVar: "LOG_LEVEL"},
		cli.StringFlag{Name: "key", Value: "block.new", Usage: "the key for subscribe", EnvVar: "KEY"},
	}
	app.Action = run
	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	lvl, err := log.ParseLevel(ctx.String("log_level"))
	if err != nil {
		log.WithError(err).Fatalln("invalid log level")
	}

	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	c, cancel := context.WithCancel(context.Background())

	rClient := redis.NewClient(&redis.Options{
		Addr:     ctx.String("redis_url"),
		Password: "",
		DB:       0,
	})
	if _, err := rClient.Ping(c).Result(); err != nil {
		log.WithError(err).Fatal("redis connection error")
	}
	log.Infof(
		"redis connect successful, %s",
		rClient.Options().Addr,
	)

	headCh, sub := blocknotify.SubBlock(c, ctx.String("endpoint_url"))

	go func() {
		for {
			select {
			case <-c.Done():
				log.Info("quit blocknotify")
			case err := <-sub.Err():
				log.WithError(err).Fatal("fail to get new block")
			case header := <-headCh:
				log.Infof("get head of block#%d", header.Number.Uint64())
				err := rClient.Publish(c, ctx.String("key"), header.Number.Uint64()).Err()
				if err != nil {
					log.WithError(err).Fatal("fail to publish new block")
				}
			}
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	log.Infoln("wait for program quiting")
	cancel()

	return nil
}
