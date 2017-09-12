package repositories

import (
	"log"
	"strings"
	"time"
)

const (
	DEFAULT_CONTEXT_LIFETIME = 10
	CONTEXT_KINO             = "kino"
)

var userContextsUpdated map[string]chan bool
var userContexts map[string]string

func init() {
	mu.Lock()
	mu.Unlock()
	userContexts = make(map[string]string)
	userContextsUpdated = make(map[string]chan bool)
}

type Processer interface {
	ShowHelp()
	ShowKinoPlaces()
	ShowFilms(location string)
	ShowTaxiList()
	GetText() string
	GetUid() string
}

func ProcessMessage(proc Processer){
	text := proc.GetText()
	uid := proc.GetUid()
	ctx, _ := userContexts[uid]

	log.Printf("[skype] text: %s context:", text, ctx)
	text = strings.Replace(text, "CherkassyBot", "", -1)
	text = strings.TrimSpace(text)

	// help
	if text == "/?" {
		proc.ShowHelp()
		return
	}

	// Taxi
	if strings.ToLower(text) == "taxi" {
		proc.ShowTaxiList()
		return
	}


	// process text with context
	if ctx != "" {
		switch ctx {
		case CONTEXT_KINO:
			proc.ShowFilms(text)
			setUserContext(uid, "")
		}

		return
	}

	// catch commands if empty context
	if ctx == "" && (strings.ToLower(text) == "kino" || strings.ToLower(text) == "films") {
		setUserContext(uid, CONTEXT_KINO)
		proc.ShowKinoPlaces()
	}
}

// set context
func setUserContext(id string, ctx string) {
	// clear context
	if ctx == "" {
		userContexts[id] = ""
		return
	}
	userContextsUpdated[id] = make(chan bool)

	// check and hold
	holded, ok := userContexts[id]
	if !ok || holded == "" {
		userContexts[id] = ctx
		go holdUserContext(id)
	}

	userContextsUpdated[id] <- true
}

// hold, update, forget context
func holdUserContext(id string) {
	defer func() {
		delete(userContexts, id)
	}()
	forgetContext := time.After(DEFAULT_CONTEXT_LIFETIME * time.Second)

	for {
		select {
		case <-userContextsUpdated[id]:
			forgetContext = time.After(DEFAULT_CONTEXT_LIFETIME * time.Second)
		case <-forgetContext:
			userContexts[id] = ""
			return
		}
	}
}
