package events

const websubDeprecationNotice = "Halt! It appears you are trying to use WebSub, which has been deprecated. For more information, see: https://discuss.dev.twitch.tv/t/deprecation-of-websub-based-webhooks/32152"

var (
	isAnonymous         bool
	forwardAddress      string
	transport           string
	noConfig            bool
	fromUser            string
	toUser              string
	giftUser            string
	eventID             string
	eventMessageID      string
	secret              string
	eventStatus         string
	subscriptionStatus  string
	itemID              string
	itemName            string
	cost                int64
	count               int
	description         string
	gameID              string
	tier                string
	timestamp           string
	charityCurrentValue int
	charityTargetValue  int
	clientId            string
	version             string
	websocketClient     string
	banStart            string
	banEnd              string
)
