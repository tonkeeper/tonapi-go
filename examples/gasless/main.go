package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tonkeeper/tonapi-go"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/contract/jetton"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	tongoWallet "github.com/tonkeeper/tongo/wallet"
)

func printConfigAndReturnRelayAddress(tonapiCli *tonapi.Client) (ton.AccountID, error) {
	cfg, err := tonapiCli.GaslessConfig(context.Background())
	if err != nil {
		return ton.AccountID{}, fmt.Errorf("failed to get gasless config: %w", err)
	}
	fmt.Printf("Available gas jettons:\n")
	for _, gasJetton := range cfg.GasJettons {
		fmt.Printf("Gas jetton master: %s\n", gasJetton.MasterID)
	}
	fmt.Printf("Relay address to send fees to: %v\n", cfg.RelayAddress)
	relayer := ton.MustParseAccountID(cfg.RelayAddress)
	return relayer, nil
}

func main() {

	// this is a simple example of how to send a gasless transfer.
	// you only need to specify your seed and a destination address.

	// the seed is not sent to the network, it is used to sign messages locally.

	seed := "..!!! REPLACE THIS WITH YOUR SEED !!! .."
	destination := ton.MustParseAccountID("... !!! REPLACE THIS WITH A CORRECT DESTINATION !!! ....")

	// we send 1 USDt to the destination.
	usdtMaster := ton.MustParseAccountID("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs") // USDt jetton master.
	jettonAmount := int64(1_000_000)                                                         // amount in the smallest jetton units. This is 1 USDt.

	// if you need to send lots of requests in parallel,
	// make sure you use a tonapi token.
	tonapiCli, err := tonapi.NewClient(tonapi.TonApiURL, &tonapi.Security{})
	if err != nil {
		panic(err)
	}
	cli, err := liteapi.NewClient(liteapi.Mainnet())
	if err != nil {
		panic(err)
	}
	// we use USDt in this example,
	// so we just print all supported gas jettons and get the relay address.
	// we have to send excess to the relay address in order to make a transfer cheaper.
	relay, err := printConfigAndReturnRelayAddress(tonapiCli)
	if err != nil {
		panic(err)
	}
	params := tonapi.GaslessEstimateParams{
		MasterID: usdtMaster.ToRaw(),
	}
	j := jetton.New(usdtMaster, cli)
	walletPrivateKey, err := tongoWallet.SeedToPrivateKey(seed)
	if err != nil {
		panic(err)
	}
	networkID, err := cli.GetNetworkGlobalID(context.Background())
	if err != nil {
		panic(err)
	}
	opts := tongoWallet.Options{
		NetworkGlobalID: &networkID,
	}
	w5 := tongoWallet.NewWalletV5R1(walletPrivateKey.Public().(ed25519.PublicKey), opts)

	msgCh := make(chan tlb.Message, 1)

	// this is a trick with proxy. we don't want to send the original transaction to the network.
	// we pass the proxy to Wallet's New function and intercepts the outgoing message.
	proxy := &proxy{
		msgCh: msgCh,
		cli:   cli,
	}
	wallet, err := tongoWallet.New(walletPrivateKey, tongoWallet.V5R1, proxy, tongoWallet.WithNetworkGlobalID(networkID))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Wallet address: %v\n", wallet.GetAddress())

	walletAddress := wallet.GetAddress()
	msg := jetton.TransferMessage{
		Jetton:              j,
		Sender:              walletAddress,
		JettonAmount:        big.NewInt(jettonAmount),
		Destination:         destination,
		ResponseDestination: &relay, // excess, because some TONs will be sent back to the relay address, commission will be lowered.
		AttachedTon:         50_000_000,
		ForwardTonAmount:    1,
	}
	err = wallet.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	m := <-msgCh

	cell := boc.NewCell()
	if err := tlb.Marshal(cell, m); err != nil {
		panic(err)
	}
	rawMessages, err := tongoWallet.ExtractRawMessages(tongoWallet.V5R1, cell)
	if err != nil {
		panic(err)
	}
	if len(rawMessages) != 1 {
		panic("invalid rawMessages")
	}
	msgBoc, err := rawMessages[0].Message.ToBocString()
	if err != nil {
		panic(err)
	}

	// msgBoc is our transfer message.

	publicKey := walletPrivateKey.Public().(ed25519.PublicKey)
	estimateReq := tonapi.GaslessEstimateReq{
		WalletAddress:   walletAddress.ToRaw(),
		WalletPublicKey: hex.EncodeToString(publicKey),
		Messages: []tonapi.GaslessEstimateReqMessagesItem{
			{Boc: msgBoc},
		},
	}
	// we send a single message containing a transfer from our wallet to a desired destination.
	// as a result of estimation, TonAPI returns a list of messages that we need to sign.
	// the first message is a fee transfer to the relay address, the second message is our original transfer.
	signRawParams, err := tonapiCli.GaslessEstimate(context.Background(), &estimateReq, params)
	if err != nil {
		panic(err)
	}

	// signRawParams is the same structure as signRawParams in tonconnect.
	var msgs []tongoWallet.Sendable
	for _, msg := range signRawParams.Messages {
		cells, err := boc.DeserializeBocHex(msg.Payload.Value)
		if err != nil {
			panic(err)
		}
		if len(cells) != 1 {
			panic("invalid cells")
		}
		dest := tongo.MustParseAccountID(msg.Address)
		amount := decimal.RequireFromString(msg.Amount)
		rawMessage := RawMessage{
			Dest:    dest,
			Amount:  amount.BigInt().Int64(),
			Payload: cells[0],
		}
		msgs = append(msgs, rawMessage)
	}

	// OK, at this point, we have everything we need to send a gasless transfer.
	state, err := cli.GetAccountState(context.Background(), wallet.GetAddress())
	if err != nil {
		panic(err)
	}
	nextMsgParams, err := w5.NextMessageParams(state)
	if err != nil {
		panic(err)
	}
	// the message has to be V5MsgTypeSignedInternal.
	conf := tongoWallet.MessageConfig{
		Seqno:      nextMsgParams.Seqno,
		ValidUntil: time.Now().UTC().Add(tongoWallet.DefaultMessageLifetime),
		V5MsgType:  tongoWallet.V5MsgTypeSignedInternal,
	}
	body, err := wallet.CreateMessageBody(conf, msgs...)
	if err != nil {
		panic(err)
	}
	m, err = ton.CreateExternalMessage(wallet.GetAddress(), body, nextMsgParams.Init, tlb.VarUInteger16{})
	if err != nil {
		panic(err)
	}

	cell = boc.NewCell()
	if err := tlb.Marshal(cell, m); err != nil {
		panic(err)
	}
	msgBoc, err = cell.ToBocBase64()
	if err != nil {
		panic(err)
	}
	sendReq := tonapi.GaslessSendReq{
		WalletPublicKey: hex.EncodeToString(publicKey),
		Boc:             msgBoc,
	}
	if _, err = tonapiCli.GaslessSend(context.Background(), &sendReq); err != nil {
		panic(err)
	}
	fmt.Printf("A gasless transfer sent\n")
}

type RawMessage struct {
	Dest    ton.AccountID
	Amount  int64
	Payload *boc.Cell
}

func (m RawMessage) ToInternal() (message tlb.Message, mode uint8, err error) {
	info := tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
	}

	info.IntMsgInfo = &struct {
		IhrDisabled bool
		Bounce      bool
		Bounced     bool
		Src         tlb.MsgAddress
		Dest        tlb.MsgAddress
		Value       tlb.CurrencyCollection
		IhrFee      tlb.Grams
		FwdFee      tlb.Grams
		CreatedLt   uint64
		CreatedAt   uint32
	}{
		IhrDisabled: true,
		Bounce:      false,
		Src:         (*ton.AccountID)(nil).ToMsgAddress(),
		Dest:        m.Dest.ToMsgAddress(),
	}
	info.IntMsgInfo.Value.Grams = tlb.Grams(m.Amount)

	intMsg := tlb.Message{
		Info: info,
	}

	intMsg.Body.IsRight = true //todo: check length and
	intMsg.Body.Value = tlb.Any(*m.Payload)
	return intMsg, tongoWallet.DefaultMessageMode, nil
}

type proxy struct {
	msgCh chan tlb.Message
	cli   *liteapi.Client
}

func (p *proxy) GetSeqno(ctx context.Context, account tongo.AccountID) (uint32, error) {
	return p.cli.GetSeqno(ctx, account)
}

func (p *proxy) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	cells, err := boc.DeserializeBoc(payload)
	if err != nil {
		panic(err)
	}
	var msg tlb.Message
	if err := tlb.Unmarshal(cells[0], &msg); err != nil {
		panic(err)
	}
	p.msgCh <- msg
	return 0, nil
}

func (p *proxy) GetAccountState(ctx context.Context, accountID tongo.AccountID) (tlb.ShardAccount, error) {
	return p.cli.GetAccountState(ctx, accountID)
}
