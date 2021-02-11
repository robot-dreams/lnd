package main

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/urfave/cli"
)

var disableChannelCommand = cli.Command{
	Name:     "disablechannel",
	Category: "Channels",
	Usage:    "Disables an existing channel.",
	Description: `
	Disables an existing channel.`,
	ArgsUsage: "funding_txid [output_index]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "funding_txid",
			Usage: "the txid of the channel's funding transaction",
		},
		cli.IntFlag{
			Name: "output_index",
			Usage: "the output index for the funding output of the funding " +
				"transaction",
		},
	},
	Action: actionDecorator(disableChannel),
}

func disableChannel(ctx *cli.Context) error {
	conn := getClientConn(ctx, false)
	defer conn.Close()

	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "disablechannel")
		return nil
	}

	client := routerrpc.NewRouterClient(conn)
	ctxb := context.Background()

	channelPoint, err := parseChannelPoint(ctx)
	if err != nil {
		return err
	}

	req := &routerrpc.DisableChannelRequest{
		ChanPoint: channelPoint,
	}

	resp, err := client.DisableChannel(ctxb, req)
	if err != nil {
		return err
	}

	printRespJSON(resp)

	return nil
}
