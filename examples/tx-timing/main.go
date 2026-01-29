package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	sse "github.com/r3labs/sse/v2"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

const (
	numRuns = 20
)

type RunResult struct {
	SSELatency  time.Duration
	LiteLatency time.Duration
	SSESuccess  bool
	LiteSuccess bool
}

func main() {
	seed := ""
	destinationStr := ""
	apiKey := ""
	amount := uint64(10_000_000) // 0.01 TON

	destination := ton.MustParseAccountID(destinationStr)

	cli, err := liteapi.NewClient(liteapi.Testnet())
	if err != nil {
		fmt.Printf("failed to create liteapi client: %v\n", err)
		os.Exit(1)
	}

	privateKey, err := wallet.SeedToPrivateKey(seed)
	if err != nil {
		fmt.Printf("failed to convert seed to private key: %v\n", err)
		os.Exit(1)
	}

	w, err := wallet.New(privateKey, wallet.V3R2, cli)
	if err != nil {
		fmt.Printf("failed to create wallet: %v\n", err)
		os.Exit(1)
	}

	walletAddress := w.GetAddress()
	fmt.Printf("Wallet address: %v\n", walletAddress.ToHuman(true, false))
	fmt.Printf("Running %d tests...\n\n", numRuns)

	var results []RunResult
	var sseLatencies, liteLatencies []time.Duration

	for i := 1; i <= numRuns; i++ {
		fmt.Printf("=== RUN %d/%d ===\n", i, numRuns)

		seqnoBefore, err := cli.GetSeqno(context.Background(), walletAddress)
		if err != nil {
			fmt.Printf("  failed to get seqno: %v\n", err)
			continue
		}

		result := runTest(cli, w, walletAddress, destination, amount, apiKey)
		results = append(results, result)

		if result.SSESuccess {
			sseLatencies = append(sseLatencies, result.SSELatency)
			fmt.Printf("  SSE:        %v\n", result.SSELatency)
		} else {
			fmt.Printf("  SSE:        ERROR\n")
		}

		if result.LiteSuccess {
			liteLatencies = append(liteLatencies, result.LiteLatency)
			fmt.Printf("  Liteserver: %v\n", result.LiteLatency)
		} else {
			fmt.Printf("  Liteserver: ERROR\n")
		}

		fmt.Println()

		if i < numRuns {
			fmt.Printf("  Waiting for seqno to increment...")
			waitForSeqnoIncrement(cli, walletAddress, seqnoBefore)
			fmt.Printf(" done\n")
		}
	}

	fmt.Printf("\n========== FINAL RESULTS ==========\n\n")

	if len(sseLatencies) > 0 {
		fmt.Printf("SSE Statistics (%d successful runs):\n", len(sseLatencies))
		fmt.Printf("  Average: %v\n", average(sseLatencies))
		fmt.Printf("  Min:     %v\n", min(sseLatencies))
		fmt.Printf("  Max:     %v\n", max(sseLatencies))
		fmt.Printf("  Median:  %v\n", median(sseLatencies))
	} else {
		fmt.Printf("SSE: No successful runs\n")
	}

	fmt.Println()

	if len(liteLatencies) > 0 {
		fmt.Printf("Liteserver Statistics (%d successful runs):\n", len(liteLatencies))
		fmt.Printf("  Average: %v\n", average(liteLatencies))
		fmt.Printf("  Min:     %v\n", min(liteLatencies))
		fmt.Printf("  Max:     %v\n", max(liteLatencies))
		fmt.Printf("  Median:  %v\n", median(liteLatencies))
	} else {
		fmt.Printf("Liteserver: No successful runs\n")
	}
}

func runTest(cli *liteapi.Client, w wallet.Wallet, walletAddress ton.AccountID, destination ton.AccountID, amount uint64, apiKey string) RunResult {
	result := RunResult{}

	initialState, err := cli.GetAccountState(context.Background(), walletAddress)
	if err != nil {
		fmt.Printf("failed to get initial account state: %v\n", err)
		return result
	}
	initialLT := initialState.Account.Account.Storage.LastTransLt

	sseFoundCh := make(chan time.Time, 1)
	liteFoundCh := make(chan time.Time, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		listenSSE(ctx, walletAddress.ToRaw(), apiKey, sseFoundCh)
	}()

	time.Sleep(500 * time.Millisecond)

	sendUtime := time.Now().UTC().Unix()
	msg := wallet.SimpleTransfer{
		Amount:  tlb.Grams(amount),
		Address: destination,
		Comment: fmt.Sprintf("timing test %d", sendUtime),
	}

	sendBeforeTime := time.Now()
	err = w.Send(context.Background(), msg)
	if err != nil {
		fmt.Printf("failed to send transaction: %v\n", err)
		cancel()
		return result
	}

	sendTime := time.Now()
	fmt.Printf("  Send time:    %v\n", sendTime.Sub(sendBeforeTime))

	wg.Add(1)
	go func() {
		defer wg.Done()
		pollLiteserver(ctx, cli, walletAddress, initialLT, liteFoundCh)
	}()

	timeout := time.After(60 * time.Second)

	for !result.SSESuccess || !result.LiteSuccess {
		select {
		case t := <-sseFoundCh:
			if !result.SSESuccess {
				result.SSELatency = t.Sub(sendTime)
				result.SSESuccess = true
			}
		case t := <-liteFoundCh:
			if !result.LiteSuccess {
				result.LiteLatency = t.Sub(sendTime)
				result.LiteSuccess = true
			}
		case <-timeout:
			cancel()
			return result
		}
	}

	cancel()
	return result
}

func listenSSE(ctx context.Context, accountRaw string, apiKey string, foundCh chan<- time.Time) {
	url := fmt.Sprintf("https://rt-testnet.tonapi.io/sse/transactions?account=%s", accountRaw)

	client := sse.NewClient(url)
	client.Headers = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
	}

	err := client.SubscribeWithContext(ctx, "", func(msg *sse.Event) {
		switch string(msg.Event) {
		case "heartbeat":
			return
		case "message":
			select {
			case foundCh <- time.Now():
			default:
			}
		}
	})

	if err != nil && ctx.Err() == nil {
		fmt.Printf("SSE error: %v\n", err)
	}
}

func waitForSeqnoIncrement(cli *liteapi.Client, account ton.AccountID, seqnoBefore uint32) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-timeout:
			return
		case <-ticker.C:
			seqno, err := cli.GetSeqno(context.Background(), account)
			if err != nil {
				continue
			}
			if seqno > seqnoBefore {
				return
			}
		}
	}
}

func pollLiteserver(ctx context.Context, cli *liteapi.Client, account ton.AccountID, initialLT uint64, foundCh chan<- time.Time) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			state, err := cli.GetAccountState(ctx, account)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}

			currentLT := state.Account.Account.Storage.LastTransLt
			if currentLT > initialLT {
				select {
				case foundCh <- time.Now():
				default:
				}
				return
			}
		}
	}
}

func average(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}

func min(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	m := durations[0]
	for _, d := range durations[1:] {
		if d < m {
			m = d
		}
	}
	return m
}

func max(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	m := durations[0]
	for _, d := range durations[1:] {
		if d > m {
			m = d
		}
	}
	return m
}

func median(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}
