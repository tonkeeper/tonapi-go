// Code generated by ogen, DO NOT EDIT.

package tonapi

// AccountDnsBackResolveParams is parameters of accountDnsBackResolve operation.
type AccountDnsBackResolveParams struct {
	// Account ID.
	AccountID string
}

// AddressParseParams is parameters of addressParse operation.
type AddressParseParams struct {
	// Account ID.
	AccountID string
}

// BlockchainAccountInspectParams is parameters of blockchainAccountInspect operation.
type BlockchainAccountInspectParams struct {
	// Account ID.
	AccountID string
}

// DnsResolveParams is parameters of dnsResolve operation.
type DnsResolveParams struct {
	// Domain name with .ton or .t.me.
	DomainName string
}

// EmulateMessageToAccountEventParams is parameters of emulateMessageToAccountEvent operation.
type EmulateMessageToAccountEventParams struct {
	AcceptLanguage OptString
	// Account ID.
	AccountID            string
	IgnoreSignatureCheck OptBool
}

// EmulateMessageToEventParams is parameters of emulateMessageToEvent operation.
type EmulateMessageToEventParams struct {
	AcceptLanguage       OptString
	IgnoreSignatureCheck OptBool
}

// EmulateMessageToTraceParams is parameters of emulateMessageToTrace operation.
type EmulateMessageToTraceParams struct {
	IgnoreSignatureCheck OptBool
}

// EmulateMessageToWalletParams is parameters of emulateMessageToWallet operation.
type EmulateMessageToWalletParams struct {
	AcceptLanguage OptString
}

// ExecGetMethodForBlockchainAccountParams is parameters of execGetMethodForBlockchainAccount operation.
type ExecGetMethodForBlockchainAccountParams struct {
	// Account ID.
	AccountID string
	// Contract get method name.
	MethodName string
	Args       []string
	FixOrder   OptBool
}

// GaslessEstimateParams is parameters of gaslessEstimate operation.
type GaslessEstimateParams struct {
	// Jetton to pay commission.
	MasterID string
}

// GetAccountParams is parameters of getAccount operation.
type GetAccountParams struct {
	// Account ID.
	AccountID string
}

// GetAccountDiffParams is parameters of getAccountDiff operation.
type GetAccountDiffParams struct {
	// Account ID.
	AccountID string
	StartDate int64
	EndDate   int64
}

// GetAccountDnsExpiringParams is parameters of getAccountDnsExpiring operation.
type GetAccountDnsExpiringParams struct {
	// Account ID.
	AccountID string
	// Number of days before expiration.
	Period OptInt
}

// GetAccountEventParams is parameters of getAccountEvent operation.
type GetAccountEventParams struct {
	// Account ID.
	AccountID string
	// Event ID or transaction hash in hex (without 0x) or base64url format.
	EventID        string
	AcceptLanguage OptString
	// Filter actions where requested account is not real subject (for example sender or receiver jettons).
	SubjectOnly OptBool
}

// GetAccountEventsParams is parameters of getAccountEvents operation.
type GetAccountEventsParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
	// Show only events that are initiated by this account.
	Initiator OptBool
	// Filter actions where requested account is not real subject (for example sender or receiver jettons).
	SubjectOnly OptBool
	// Omit this parameter to get last events.
	BeforeLt  OptInt64
	Limit     int
	StartDate OptInt64
	EndDate   OptInt64
}

// GetAccountExtraCurrencyHistoryByIDParams is parameters of getAccountExtraCurrencyHistoryByID operation.
type GetAccountExtraCurrencyHistoryByIDParams struct {
	// Account ID.
	AccountID string
	// Extra currency id.
	ID             int32
	AcceptLanguage OptString
	// Omit this parameter to get last events.
	BeforeLt  OptInt64
	Limit     int
	StartDate OptInt64
	EndDate   OptInt64
}

// GetAccountInscriptionsParams is parameters of getAccountInscriptions operation.
type GetAccountInscriptionsParams struct {
	// Account ID.
	AccountID string
	Limit     OptInt
	Offset    OptInt
}

// GetAccountInscriptionsHistoryParams is parameters of getAccountInscriptionsHistory operation.
type GetAccountInscriptionsHistoryParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
	// Omit this parameter to get last events.
	BeforeLt OptInt64
	Limit    OptInt
}

// GetAccountInscriptionsHistoryByTickerParams is parameters of getAccountInscriptionsHistoryByTicker operation.
type GetAccountInscriptionsHistoryByTickerParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
	Ticker         string
	// Omit this parameter to get last events.
	BeforeLt OptInt64
	Limit    OptInt
}

// GetAccountJettonBalanceParams is parameters of getAccountJettonBalance operation.
type GetAccountJettonBalanceParams struct {
	// Account ID.
	AccountID string
	// Jetton ID.
	JettonID string
	// Accept ton and all possible fiat currencies, separated by commas.
	Currencies []string
	// Comma separated list supported extensions.
	SupportedExtensions []string
}

// GetAccountJettonHistoryByIDParams is parameters of getAccountJettonHistoryByID operation.
type GetAccountJettonHistoryByIDParams struct {
	// Account ID.
	AccountID string
	// Jetton ID.
	JettonID       string
	AcceptLanguage OptString
	// Omit this parameter to get last events.
	BeforeLt  OptInt64
	Limit     int
	StartDate OptInt64
	EndDate   OptInt64
}

// GetAccountJettonsBalancesParams is parameters of getAccountJettonsBalances operation.
type GetAccountJettonsBalancesParams struct {
	// Account ID.
	AccountID string
	// Accept ton and all possible fiat currencies, separated by commas.
	Currencies []string
	// Comma separated list supported extensions.
	SupportedExtensions []string
}

// GetAccountJettonsHistoryParams is parameters of getAccountJettonsHistory operation.
type GetAccountJettonsHistoryParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
	// Omit this parameter to get last events.
	BeforeLt  OptInt64
	Limit     int
	StartDate OptInt64
	EndDate   OptInt64
}

// GetAccountMultisigsParams is parameters of getAccountMultisigs operation.
type GetAccountMultisigsParams struct {
	// Account ID.
	AccountID string
}

// GetAccountNftHistoryParams is parameters of getAccountNftHistory operation.
type GetAccountNftHistoryParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
	// Omit this parameter to get last events.
	BeforeLt  OptInt64
	Limit     int
	StartDate OptInt64
	EndDate   OptInt64
}

// GetAccountNftItemsParams is parameters of getAccountNftItems operation.
type GetAccountNftItemsParams struct {
	// Account ID.
	AccountID string
	// Nft collection.
	Collection OptString
	Limit      OptInt
	Offset     OptInt
	// Selling nft items in ton implemented usually via transfer items to special selling account. This
	// option enables including items which owned not directly.
	IndirectOwnership OptBool
}

// GetAccountNominatorsPoolsParams is parameters of getAccountNominatorsPools operation.
type GetAccountNominatorsPoolsParams struct {
	// Account ID.
	AccountID string
}

// GetAccountPublicKeyParams is parameters of getAccountPublicKey operation.
type GetAccountPublicKeyParams struct {
	// Account ID.
	AccountID string
}

// GetAccountSeqnoParams is parameters of getAccountSeqno operation.
type GetAccountSeqnoParams struct {
	// Account ID.
	AccountID string
}

// GetAccountSubscriptionsParams is parameters of getAccountSubscriptions operation.
type GetAccountSubscriptionsParams struct {
	// Account ID.
	AccountID string
}

// GetAccountTracesParams is parameters of getAccountTraces operation.
type GetAccountTracesParams struct {
	// Account ID.
	AccountID string
	// Omit this parameter to get last events.
	BeforeLt OptInt64
	Limit    OptInt
}

// GetAccountsParams is parameters of getAccounts operation.
type GetAccountsParams struct {
	Currency OptString
}

// GetAllAuctionsParams is parameters of getAllAuctions operation.
type GetAllAuctionsParams struct {
	// Domain filter for current auctions "ton" or "t.me".
	Tld OptString
}

// GetAllRawShardsInfoParams is parameters of getAllRawShardsInfo operation.
type GetAllRawShardsInfoParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
}

// GetBlockchainAccountTransactionsParams is parameters of getBlockchainAccountTransactions operation.
type GetBlockchainAccountTransactionsParams struct {
	// Account ID.
	AccountID string
	// Omit this parameter to get last transactions.
	AfterLt OptInt64
	// Omit this parameter to get last transactions.
	BeforeLt  OptInt64
	Limit     OptInt32
	SortOrder OptGetBlockchainAccountTransactionsSortOrder
}

// GetBlockchainBlockParams is parameters of getBlockchainBlock operation.
type GetBlockchainBlockParams struct {
	// Block ID.
	BlockID string
}

// GetBlockchainBlockTransactionsParams is parameters of getBlockchainBlockTransactions operation.
type GetBlockchainBlockTransactionsParams struct {
	// Block ID.
	BlockID string
}

// GetBlockchainConfigFromBlockParams is parameters of getBlockchainConfigFromBlock operation.
type GetBlockchainConfigFromBlockParams struct {
	// Masterchain block seqno.
	MasterchainSeqno int32
}

// GetBlockchainMasterchainBlocksParams is parameters of getBlockchainMasterchainBlocks operation.
type GetBlockchainMasterchainBlocksParams struct {
	// Masterchain block seqno.
	MasterchainSeqno int32
}

// GetBlockchainMasterchainShardsParams is parameters of getBlockchainMasterchainShards operation.
type GetBlockchainMasterchainShardsParams struct {
	// Masterchain block seqno.
	MasterchainSeqno int32
}

// GetBlockchainMasterchainTransactionsParams is parameters of getBlockchainMasterchainTransactions operation.
type GetBlockchainMasterchainTransactionsParams struct {
	// Masterchain block seqno.
	MasterchainSeqno int32
}

// GetBlockchainRawAccountParams is parameters of getBlockchainRawAccount operation.
type GetBlockchainRawAccountParams struct {
	// Account ID.
	AccountID string
}

// GetBlockchainTransactionParams is parameters of getBlockchainTransaction operation.
type GetBlockchainTransactionParams struct {
	// Transaction ID.
	TransactionID string
}

// GetBlockchainTransactionByMessageHashParams is parameters of getBlockchainTransactionByMessageHash operation.
type GetBlockchainTransactionByMessageHashParams struct {
	// Message ID.
	MsgID string
}

// GetChartRatesParams is parameters of getChartRates operation.
type GetChartRatesParams struct {
	// Accept jetton master address.
	Token       string
	Currency    OptString
	StartDate   OptInt64
	EndDate     OptInt64
	PointsCount OptInt
}

// GetDnsInfoParams is parameters of getDnsInfo operation.
type GetDnsInfoParams struct {
	// Domain name with .ton or .t.me.
	DomainName string
}

// GetDomainBidsParams is parameters of getDomainBids operation.
type GetDomainBidsParams struct {
	// Domain name with .ton or .t.me.
	DomainName string
}

// GetEventParams is parameters of getEvent operation.
type GetEventParams struct {
	// Event ID or transaction hash in hex (without 0x) or base64url format.
	EventID        string
	AcceptLanguage OptString
}

// GetExtraCurrencyInfoParams is parameters of getExtraCurrencyInfo operation.
type GetExtraCurrencyInfoParams struct {
	// Extra currency id.
	ID int32
}

// GetInscriptionOpTemplateParams is parameters of getInscriptionOpTemplate operation.
type GetInscriptionOpTemplateParams struct {
	Type        GetInscriptionOpTemplateType
	Destination OptString
	Comment     OptString
	Operation   GetInscriptionOpTemplateOperation
	Amount      string
	Ticker      string
	Who         string
}

// GetItemsFromCollectionParams is parameters of getItemsFromCollection operation.
type GetItemsFromCollectionParams struct {
	// Account ID.
	AccountID string
	Limit     OptInt
	Offset    OptInt
}

// GetJettonHoldersParams is parameters of getJettonHolders operation.
type GetJettonHoldersParams struct {
	// Account ID.
	AccountID string
	Limit     OptInt
	Offset    OptInt
}

// GetJettonInfoParams is parameters of getJettonInfo operation.
type GetJettonInfoParams struct {
	// Account ID.
	AccountID string
}

// GetJettonTransferPayloadParams is parameters of getJettonTransferPayload operation.
type GetJettonTransferPayloadParams struct {
	// Account ID.
	AccountID string
	// Jetton ID.
	JettonID string
}

// GetJettonsParams is parameters of getJettons operation.
type GetJettonsParams struct {
	Limit  OptInt32
	Offset OptInt32
}

// GetJettonsEventsParams is parameters of getJettonsEvents operation.
type GetJettonsEventsParams struct {
	// Event ID or transaction hash in hex (without 0x) or base64url format.
	EventID        string
	AcceptLanguage OptString
}

// GetMultisigAccountParams is parameters of getMultisigAccount operation.
type GetMultisigAccountParams struct {
	// Account ID.
	AccountID string
}

// GetNftCollectionParams is parameters of getNftCollection operation.
type GetNftCollectionParams struct {
	// Account ID.
	AccountID string
}

// GetNftCollectionsParams is parameters of getNftCollections operation.
type GetNftCollectionsParams struct {
	Limit  OptInt32
	Offset OptInt32
}

// GetNftHistoryByIDParams is parameters of getNftHistoryByID operation.
type GetNftHistoryByIDParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
	// Omit this parameter to get last events.
	BeforeLt  OptInt64
	Limit     int
	StartDate OptInt64
	EndDate   OptInt64
}

// GetNftItemByAddressParams is parameters of getNftItemByAddress operation.
type GetNftItemByAddressParams struct {
	// Account ID.
	AccountID string
}

// GetRatesParams is parameters of getRates operation.
type GetRatesParams struct {
	// Accept ton and jetton master addresses, separated by commas.
	Tokens []string
	// Accept ton and all possible fiat currencies, separated by commas.
	Currencies []string
}

// GetRawAccountStateParams is parameters of getRawAccountState operation.
type GetRawAccountStateParams struct {
	// Account ID.
	AccountID string
	// Target block: (workchain,shard,seqno,root_hash,file_hash).
	TargetBlock OptString
}

// GetRawBlockProofParams is parameters of getRawBlockProof operation.
type GetRawBlockProofParams struct {
	// Known block: (workchain,shard,seqno,root_hash,file_hash).
	KnownBlock string
	// Target block: (workchain,shard,seqno,root_hash,file_hash).
	TargetBlock OptString
	// Mode.
	Mode int32
}

// GetRawBlockchainBlockParams is parameters of getRawBlockchainBlock operation.
type GetRawBlockchainBlockParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
}

// GetRawBlockchainBlockHeaderParams is parameters of getRawBlockchainBlockHeader operation.
type GetRawBlockchainBlockHeaderParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
	// Mode.
	Mode int32
}

// GetRawBlockchainBlockStateParams is parameters of getRawBlockchainBlockState operation.
type GetRawBlockchainBlockStateParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
}

// GetRawBlockchainConfigFromBlockParams is parameters of getRawBlockchainConfigFromBlock operation.
type GetRawBlockchainConfigFromBlockParams struct {
	// Masterchain block seqno.
	MasterchainSeqno int32
}

// GetRawConfigParams is parameters of getRawConfig operation.
type GetRawConfigParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
	// Mode.
	Mode int32
}

// GetRawListBlockTransactionsParams is parameters of getRawListBlockTransactions operation.
type GetRawListBlockTransactionsParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
	// Mode.
	Mode int32
	// Count.
	Count int32
	// Account ID.
	AccountID OptString
	// Lt.
	Lt OptInt64
}

// GetRawMasterchainInfoExtParams is parameters of getRawMasterchainInfoExt operation.
type GetRawMasterchainInfoExtParams struct {
	// Mode.
	Mode int32
}

// GetRawShardBlockProofParams is parameters of getRawShardBlockProof operation.
type GetRawShardBlockProofParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
}

// GetRawShardInfoParams is parameters of getRawShardInfo operation.
type GetRawShardInfoParams struct {
	// Block ID: (workchain,shard,seqno,root_hash,file_hash).
	BlockID string
	// Workchain.
	Workchain int32
	// Shard.
	Shard int64
	// Exact.
	Exact bool
}

// GetRawTransactionsParams is parameters of getRawTransactions operation.
type GetRawTransactionsParams struct {
	// Account ID.
	AccountID string
	// Count.
	Count int32
	// Lt.
	Lt int64
	// Hash.
	Hash string
}

// GetReducedBlockchainBlocksParams is parameters of getReducedBlockchainBlocks operation.
type GetReducedBlockchainBlocksParams struct {
	From int64
	To   int64
}

// GetStakingPoolHistoryParams is parameters of getStakingPoolHistory operation.
type GetStakingPoolHistoryParams struct {
	// Account ID.
	AccountID string
}

// GetStakingPoolInfoParams is parameters of getStakingPoolInfo operation.
type GetStakingPoolInfoParams struct {
	// Account ID.
	AccountID      string
	AcceptLanguage OptString
}

// GetStakingPoolsParams is parameters of getStakingPools operation.
type GetStakingPoolsParams struct {
	// Account ID.
	AvailableFor OptString
	// Return also pools not from white list - just compatible by interfaces (maybe dangerous!).
	IncludeUnverified OptBool
	AcceptLanguage    OptString
}

// GetTraceParams is parameters of getTrace operation.
type GetTraceParams struct {
	// Trace ID or transaction hash in hex (without 0x) or base64url format.
	TraceID string
}

// GetWalletsByPublicKeyParams is parameters of getWalletsByPublicKey operation.
type GetWalletsByPublicKeyParams struct {
	PublicKey string
}

// ReindexAccountParams is parameters of reindexAccount operation.
type ReindexAccountParams struct {
	// Account ID.
	AccountID string
}

// SearchAccountsParams is parameters of searchAccounts operation.
type SearchAccountsParams struct {
	Name string
}
