package fix

type StorageSide string

const (
	Incoming StorageSide = "incoming"
	Outgoing             = "outgoing"
)

type StorageID struct {
	Sender string
	Target string
	Side   StorageSide
}
