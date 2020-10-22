package cerror

import "fmt"

const (
	TYPE_UNDEFINED    = 0
	TYPE_NOT_FOUND    = 1
	TYPE_UNAUTHORIZED = 2
)

type Error struct {
	Tag             string
	Err             error
	FriendlyMessage string
	Type            int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) PrintErrorWithTag() {
	fmt.Println("[%s] %s ", e.Tag, e.Err)
}

func (e Error) FriendlyMessageWithTag() string {
	return fmt.Sprintf("%s [%s]", e.FriendlyMessage, e.Tag)
}

func New(tag string, err error, friendly string) error {
	return &Error{tag, err, friendly, TYPE_UNDEFINED}
}

func NewAndPrintWithTag(tag string, err error, friendly string) Error {
	cerr := Error{tag, err, friendly, TYPE_UNDEFINED}
	cerr.PrintErrorWithTag()
	return cerr
}
