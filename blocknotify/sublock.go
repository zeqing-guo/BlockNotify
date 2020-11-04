package blocknotify

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func SubBlock(
	ctx context.Context,
	webSocket string,
) (chan *types.Header, ethereum.Subscription) {
	wsclient, err := ethclient.DialContext(ctx, webSocket)
	if err != nil {
		log.WithError(err).Fatalln("fail to dial")
	}

	headers := make(chan *types.Header)
	sub, err := wsclient.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.WithError(err).Fatalln("fail to sub")
	}
	return headers, sub
}
