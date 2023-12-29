package props

const (
	CHANNEL_SECRET string = "CHANNEL_SECRET"
	CHANNEL_TOKEN  string = "CHANNEL_TOKEN"
)

var (
	ChannelSecret string
	ChannelToken  string
)

func loadLineProps() {
	ChannelSecret = readEnv(CHANNEL_SECRET, "channelSecret")
	ChannelToken = readEnv(CHANNEL_TOKEN, "channelToken")
}
