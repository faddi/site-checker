package checker

import (
    "testing"
    "time"
)

var curl string = "http://www.example.com"
var delay time.Duration = 2

func Test_New(t *testing.T){
    c := New()

    if c == nil {
        t.Error("nil returned from new")
    }

    if c.sites == nil {
        t.Error("urls not initiated")
    }
}

func Test_AddUrl(t *testing.T){

    c := New()

    if err := c.AddUrl("i am not a valid url", delay); err == nil {
        t.Error("Did not fail on invalid url")
    }

    if err := c.AddUrl("/ddfsd/das", delay); err == nil {
        t.Error("Did not fail on relative url")
    }

    if err := c.AddUrl(curl, delay); err != nil {
        t.Error("Failed to add proper url")
    }

    if _, ok := c.sites[curl]; ok != true {
        t.Error("A site should exist in the sites map if it is created without an error")
    }

}

func Test_StopUrl(t *testing.T){

    c := New()

    if err := c.AddUrl(curl, delay); err != nil {
        t.Error("Failed to add proper url")
    }

    time.Sleep(3 * time.Second)
    log("%v\n", <-c.ResultChan())
    err := c.StopCheckingUrl(curl)

    if err != nil {
        t.Fatal(err)
    }

    if _, ok := c.sites[curl]; ok == true {
        t.Error("A site should not exist in checker.sites after stop")
    }
}

func Test_Multiple(t *testing.T){

    c := New()

    urls := []string{curl, "http://www.google.com", "http://www.dn.se", "http://www.aftonbladet.se"}

    for _, u := range urls {
        if err := c.AddUrl(u, delay); err != nil {
            t.Error("Failed to add proper url")
        }
    }

    go func () {
        out := c.out
        for {
            d := <-out
            log("%s -> %s", d.Url, d.Resp.Status)

        }
    }()

    time.Sleep(5 * time.Second)

    for _, u := range urls {
        err := c.StopCheckingUrl(u)
        if err != nil {
            t.Fatal(err)
        }
    }

    if len(c.sites) > 0 {
        t.Error("No c.sites should be empty after all urls have been stopped")
    }

}

/*func Test_blah(t *testing.T){
    t.Error("dasd")
}*/
