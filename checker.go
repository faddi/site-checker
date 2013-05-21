package checker

import (
    "errors"
    stdlog "log"
    "net/http"
    "net/url"
    "time"
)

var l *stdlog.Logger

func log(format string, v ...interface{}) {
    if l != nil {
        l.Printf(format, v...)
    }
}

func SetLogger(log *stdlog.Logger) {
    l = log
}

type CheckResult struct {
    Resp       *http.Response
    Body       []byte
    Connecting time.Duration
    Receiving  time.Duration
    Timestamp  time.Time
    Url        string
    Error      error
}

type Checker struct {
    sites map[string]*site
    out   chan *CheckResult
}

// Initiates a new Checker
func New() *Checker {

    c := new(Checker)
    c.sites = make(map[string]*site)
    c.out = make(chan *CheckResult)

    log("Checker init")

    return c
}

func (c *Checker) AddUrl(rawUrl string, delay time.Duration) error {

    log("Adding url: %s\n", rawUrl)

    u, err := url.Parse(rawUrl)

    if u.IsAbs() == false || err != nil {
        return errors.New("Invalid url : " + rawUrl)
    }

    if _, exists := c.sites[rawUrl]; exists == true {
        return errors.New("Checker already has : " + rawUrl)
    }

    s := newSite(u, delay, c.out)

    go s.start()

    c.sites[u.String()] = s

    return nil
}

func (c *Checker) ResultChan() <-chan *CheckResult {
    return c.out
}

func (c *Checker) StopCheckingUrl(rawUrl string) error {

    if s, ok := c.sites[rawUrl]; ok == true {
        log("Stopping checks for site %s\n", s.url.String())
        s.stop <- true
        delete(c.sites, rawUrl)
        return nil
    }

    return errors.New("Site " + rawUrl + " not found.")
}
