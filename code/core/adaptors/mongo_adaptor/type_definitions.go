package mongo_adaptor

type DBSpace int
const (
	SHARED_DB DBSpace = iota
	APP_DB
)

