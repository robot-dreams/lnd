package main

import (
	"context"
	"errors"

	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/urfave/cli"
)

var setChannelStatusCommand = cli.Command{
	Name:     "setchannelstatus",
	Category: "Channels",
	Usage:    "Sets the status of an existing channel.",
	Description: `
	Manually sets the status of an existing channel.`,
	ArgsUsage: "funding_txid [output_index] action",
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
		cli.StringFlag{
			Name:  "action",
			Usage: "the action to take, must be one of enable|disable|auto",
		},
	},
	Action: actionDecorator(setChannelStatus),
}

func setChannelStatus(ctx *cli.Context) error {
	conn := getClientConn(ctx, false)
	defer conn.Close()

	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "setchannelstatus")
		return nil
	}

	channelPoint, err := parseChannelPoint(ctx)
	if err != nil {
		return err
	}

	var action routerrpc.SetChannelStatusAction
	switch ctx.String("action") {
	case "enable":
		action = routerrpc.SetChannelStatusAction_ENABLE
	case "disable":
		action = routerrpc.SetChannelStatusAction_DISABLE
	case "auto":
		action = routerrpc.SetChannelStatusAction_AUTO
	default:
		return errors.New("action must be one of enable|disable|auto")
	}
	req := &routerrpc.SetChannelStatusRequest{
		ChanPoint: channelPoint,
		Action:    action,
	}

	client := routerrpc.NewRouterClient(conn)
	ctxb := context.Background()
	resp, err := client.SetChannelStatus(ctxb, req)
	if err != nil {
		return err
	}

	printRespJSON(resp)

	return nil
}
