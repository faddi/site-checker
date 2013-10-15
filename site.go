package checker

import (
    "errors"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
)

type site struct {
    url    *url.URL
    out    chan *CheckResult
    client *http.Client
    stop   chan bool
    delay  time.Duration
}

func newSite(u *url.URL, delay time.Duration, out chan *CheckResult) *site {
    s := new(site)

    s.client = &http.Client{CheckRedirect: checkRedirect}
    s.out = out
    s.stop = make(chan bool)
    s.url = u
    s.delay = delay

    return s
}

func checkRedirect(req *http.Request, via []*http.Request) error {
    return errors.New("")
}

func (s *site) start() {

    log("Checking site %s every %f seconds\n", s.url.String(), s.delay.Seconds())
    t := time.Tick(s.delay)

    for {
        select {
        case _ = <-t:
            s.check()
        case _ = <-s.stop:
            return
        }
    }
}

func (s *site) check() {

    log("Getting %s \n", s.url.String())

    result := new(CheckResult)

    start := time.Now()
    resp, err := s.client.Get(s.url.String())
    end := time.Now()

    if err != nil {
        if _, ok := err.(*url.Error); ok == false || resp == nil {
            result.Error = err
            log("err : " + err.Error())

            s.out <- result
            return
        }
    }

    err = nil

    var data []byte

    if resp.StatusCode % 300 > 99 {
        data, err = ioutil.ReadAll(resp.Body)
        defer resp.Body.Close()
    }


    if err != nil {
        log("error when reading response body : %s", err.Error())
        return
    }

    result.Resp = resp
    result.Body = data
    result.ResponseTime = end.Sub(start)
    result.Timestamp = start
    result.Url = s.url.String()

    s.out <- result
}
