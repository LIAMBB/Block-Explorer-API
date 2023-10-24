package main

import (
	"time"

	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// NMC chain parameters
var nmcMainnetChainParams = chaincfg.Params{
	Name:        "mainnet",
	Net:         0xfeb4bef9,
	DefaultPort: "18336",
	DNSSeeds: []chaincfg.DNSSeed{
		{Host: "seed.bitcoin.sipa.be", HasFiltering: true},
		{Host: "dnsseed.bluematt.me", HasFiltering: true},
		{Host: "dnsseed.bitcoin.dashjr.org", HasFiltering: false},
		{Host: "seed.bitcoinstats.com", HasFiltering: true},
		{Host: "seed.bitnodes.io", HasFiltering: false},
		{Host: "seed.bitcoin.jonasschnelli.ch", HasFiltering: true},
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x1d00ffff,
	BIP0034Height:            227931, // 000000000000024b89b42a942fe0d9fea3bb44ab7bd1b19115dd6a759c0808b8
	BIP0065Height:            388381, // 000000000000000004c2b624ed5d7756c508d90fd0da2c7c679febfa6c4735f0
	BIP0066Height:            363725, // 00000000000000000379eaa19dce8c9b722d46ae6a57c2f1a988119488b50931
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []chaincfg.Checkpoint{
		{Height: 11111, Hash: newHashFromStr("0000000069e244f73d78e8fd29ba2fd2ed618bd6fa2ee92559f542fdb26e7c1d")},
		{Height: 33333, Hash: newHashFromStr("000000002dd5588a74784eaa7ab0507a18ad16a236e7b1ce69f00d7ddfb5d0a6")},
		{Height: 74000, Hash: newHashFromStr("0000000000573993a3c9e41ce34471c079dcf5f52a0e824a81e7f953b8661a20")},
		{Height: 105000, Hash: newHashFromStr("00000000000291ce28027faea320c8d2b054b2e0fe44a773f3eefb151d6bdc97")},
		{Height: 134444, Hash: newHashFromStr("00000000000005b12ffd4cd315cd34ffd4a594f430ac814c91184a0d42d2b0fe")},
		{Height: 168000, Hash: newHashFromStr("000000000000099e61ea72015e79632f216fe6cb33d7899acb35b75c8303b763")},
		{Height: 193000, Hash: newHashFromStr("000000000000059f452a5f7340de6682a977387c17010ff6e6c3bd83ca8b1317")},
		{Height: 210000, Hash: newHashFromStr("000000000000048b95347e83192f69cf0366076336c639f9b7228e9ba171342e")},
		{Height: 216116, Hash: newHashFromStr("00000000000001b4f4b433e81ee46494af945cf96014816a4e2370f11b23df4e")},
		{Height: 225430, Hash: newHashFromStr("00000000000001c108384350f74090433e7fcf79a606b8e797f065b130575932")},
		{Height: 250000, Hash: newHashFromStr("000000000000003887df1f29024b06fc2200b55f8af8f35453d7be294df2d214")},
		{Height: 267300, Hash: newHashFromStr("000000000000000a83fbd660e918f218bf37edd92b748ad940483c7c116179ac")},
		{Height: 279000, Hash: newHashFromStr("0000000000000001ae8c72a0b0c301f67e3afca10e819efa9041e458e9bd7e40")},
		{Height: 300255, Hash: newHashFromStr("0000000000000000162804527c6e9b9f0563a280525f9d08c12041def0a0f3b2")},
		{Height: 319400, Hash: newHashFromStr("000000000000000021c6052e9becade189495d1c539aa37c58917305fd15f13b")},
		{Height: 343185, Hash: newHashFromStr("0000000000000000072b8bf361d01a6ba7d445dd024203fafc78768ed4368554")},
		{Height: 352940, Hash: newHashFromStr("000000000000000010755df42dba556bb72be6a32f3ce0b6941ce4430152c9ff")},
		{Height: 382320, Hash: newHashFromStr("00000000000000000a8dc6ed5b133d0eb2fd6af56203e4159789b092defd8ab2")},
		{Height: 400000, Hash: newHashFromStr("000000000000000004ec466ce4732fe6f1ed1cddc2ed4b328fff5224276e3f6f")},
		{Height: 430000, Hash: newHashFromStr("000000000000000001868b2bb3a285f3cc6b33ea234eb70facf4dcdf22186b87")},
		{Height: 460000, Hash: newHashFromStr("000000000000000000ef751bbce8e744ad303c47ece06c8d863e4d417efc258c")},
		{Height: 490000, Hash: newHashFromStr("000000000000000000de069137b17b8d5a3dfbd5b145b2dcfb203f15d0c4de90")},
		{Height: 520000, Hash: newHashFromStr("0000000000000000000d26984c0229c9f6962dc74db0a6d525f2f1640396f69c")},
		{Height: 550000, Hash: newHashFromStr("000000000000000000223b7a2298fb1c6c75fb0efc28a4c56853ff4112ec6bc9")},
		{Height: 560000, Hash: newHashFromStr("0000000000000000002c7b276daf6efb2b6aa68e2ce3be67ef925b3264ae7122")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	Deployments: [chaincfg.DefinedDeployments]chaincfg.ConsensusDeployment{
		chaincfg.DeploymentTestDummy: {
			BitNumber: 28,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Unix(11991456010, 0), // January 1, 2008 UTC
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Unix(1230767999, 0), // December 31, 2008 UTC
			),
		},
		chaincfg.DeploymentTestDummyMinActivation: {
			BitNumber:                 22,
			CustomActivationThreshold: 1815,    // Only needs 90% hash rate.
			MinActivationHeight:       10_0000, // Can only activate after height 10k.
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Time{}, // Always available for vote
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Time{}, // Never expires
			),
		},
		chaincfg.DeploymentCSV: {
			BitNumber: 0,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Unix(1462060800, 0), // May 1st, 2016
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Unix(1493596800, 0), // May 1st, 2017
			),
		},
		chaincfg.DeploymentSegwit: {
			BitNumber: 1,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Unix(1479168000, 0), // November 15, 2016 UTC
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Unix(1510704000, 0), // November 15, 2017 UTC.
			),
		},
		chaincfg.DeploymentTaproot: {
			BitNumber: 2,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Unix(1619222400, 0), // April 24th, 2021 UTC.
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Unix(1628640000, 0), // August 11th, 2021 UTC.
			),
			CustomActivationThreshold: 1815, // 90%
			MinActivationHeight:       709_632,
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "nc", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x34, // starts with 1
	ScriptHashAddrID:        0x0d, // starts with 3
	PrivateKeyID:            0xe4, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0xd7, 0xdc, 0x6e, 0x31}, // starts with xprv
	HDPublicKeyID:  [4]byte{0xd7, 0xdd, 0x63, 0x70}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0,
}

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(0x495fab29, 0), // 2009-01-03 18:15:05 +0000 UTC
		Bits:       0x1d00ffff,               // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x7c2bac1d,               // 2083236893
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
	0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
	0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
	0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a,
})

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, /* |.......E| */
				0x54, 0x68, 0x65, 0x20, 0x54, 0x69, 0x6d, 0x65, /* |The Time| */
				0x73, 0x20, 0x30, 0x33, 0x2f, 0x4a, 0x61, 0x6e, /* |s 03/Jan| */
				0x2f, 0x32, 0x30, 0x30, 0x39, 0x20, 0x43, 0x68, /* |/2009 Ch| */
				0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x72, /* |ancellor| */
				0x20, 0x6f, 0x6e, 0x20, 0x62, 0x72, 0x69, 0x6e, /* | on brin| */
				0x6b, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x65, 0x63, /* |k of sec|*/
				0x6f, 0x6e, 0x64, 0x20, 0x62, 0x61, 0x69, 0x6c, /* |ond bail| */
				0x6f, 0x75, 0x74, 0x20, 0x66, 0x6f, 0x72, 0x20, /* |out for |*/
				0x62, 0x61, 0x6e, 0x6b, 0x73, /* |banks| */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
				0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
				0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
				0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
				0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
				0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
				0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
				0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
				0x1d, 0x5f, 0xac, /* |._.| */
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
	0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
	0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
	0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00,
})
var bigOne = big.NewInt(1)

var mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)

// newHashFromStr converts the passed big-endian hex string into a
// chainhash.Hash.  It only differs from the one available in chainhash in that
// it panics on an error since it will only (and must only) be called with
// hard-coded, and therefore known good, hashes.
func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it
		// can take applications down without them having a chance to
		// recover which is extremely annoying, however an exception is
		// being made in this case because the only way this can panic
		// is if there is an error in the hard-coded hashes.  Thus it
		// will only ever potentially panic on init and therefore is
		// 100% predictable.
		panic(err)
	}
	return hash
}
