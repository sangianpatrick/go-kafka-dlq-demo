package eventbus

// MessageHeaders is type of message headers
type MessageHeaders map[string]string

// Add will add the key and value to headers.
func (mh MessageHeaders) Add(key, value string) {
	mh[key] = value
}
