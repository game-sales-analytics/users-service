package passhash

type Err error

var (
	ErrInvalidHash         Err
	ErrIncompatibleVersion Err
	ErrGenerateRandom      Err
)
