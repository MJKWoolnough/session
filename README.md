# session
--
    import "github.com/MJKWoolnough/session"


## Usage

#### type Session

```go
type Session struct {
}
```

Session stores data.

The data will be removed from storage after a timeout has expired. The timeout
is reset if the data is read from storage.

#### func  New

```go
func New(timeout time.Duration) *Session
```
New creates a new session with the given timeout

#### func (*Session) Get

```go
func (s *Session) Get(name string) interface{}
```
Get retrieves a value from storage, resetting the timer when it does so

#### func (*Session) Remove

```go
func (s *Session) Remove(name string)
```
Remove removes the named item from storage

#### func (*Session) Set

```go
func (s *Session) Set(name string, data interface{})
```
Set places the data into the session, setting the timer.

If something with the same name already exists it will be removed from the
storage
