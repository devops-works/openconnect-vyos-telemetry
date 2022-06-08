package openconnectmetrics

// Entries is a list of Entry from occtl output
type Entries []Entry

// Entry contains a single entry from occtl show users output
type Entry struct {
	ID              int      `json:"ID"`
	Username        string   `json:"Username"`
	Groupname       string   `json:"Groupname"`
	State           string   `json:"State"`
	Vhost           string   `json:"vhost"`
	Device          string   `json:"Device"`
	Mtu             string   `json:"MTU"`
	RemoteIP        string   `json:"Remote IP"`
	Location        string   `json:"Location"`
	LocalDeviceIP   string   `json:"Local Device IP"`
	IPv4            string   `json:"IPv4"`
	PTPIPv4         string   `json:"P-t-P IPv4"`
	UserAgent       string   `json:"User-Agent"`
	Rx              string   `json:"RX"`
	Tx              string   `json:"TX"`
	AverageRX       string   `json:"Average RX"`
	AverageTX       string   `json:"Average TX"`
	Dpd             string   `json:"DPD"`
	KeepAlive       string   `json:"KeepAlive"`
	Hostname        string   `json:"Hostname"`
	ConnectedAt     string   `json:"Connected at"`
	FullSession     string   `json:"Full session"`
	Session         string   `json:"Session"`
	TLSCiphersuite  string   `json:"TLS ciphersuite"`
	DTLSCipher      string   `json:"DTLS cipher"`
	DNS             []string `json:"DNS"`
	Nbns            []string `json:"NBNS"`
	SplitDNSDomains []string `json:"Split-DNS-Domains"`
	// Routes             []string `json:"Routes"`
	// NoRoutes           []string `json:"No-routes"`
	// IRoutes            []string `json:"iRoutes"`
	// RestrictedToRoutes string   `json:"Restricted to routes"`
	// RestrictedToPorts  []string `json:"Restricted to ports"`
}
