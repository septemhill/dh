package cache

type ShortURLCache interface {
	StoreShortURL(key, url string) error
	LookupShortURLByKey(key string) (string, error)
	StoreShortURLKeyOffset(offset int64) error
	RetrieveShortURLKeyOffset() (int64, error)
}
