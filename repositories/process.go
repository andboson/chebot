package repositories

import (
	"github.com/labstack/gommon/log"
	"strings"
	"time"
)

const (
	DEFAULT_CONTEXT_LIFETIME = 20
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

func ProcessMessage(proc Processer) bool {
	text := proc.GetText()
	uid := proc.GetUid()
	ctx, _ := userContexts[uid]

	log.Printf("[process] text: %s context: %s", text, ctx)
	text = strings.Replace(text, "CherkassyBot", "", -1)
	text = strings.Replace(text, "CherkasyBot", "", -1)
	text = strings.TrimSpace(text)

	// help
	if text == "/?" || strings.ToLower(text) == "/help" || strings.ToLower(text) == "help" {
		proc.ShowHelp()
		return true
	}

	// Taxi
	if strings.ToLower(text) == "taxi" || strings.ToLower(text) == "/taxi" {
		proc.ShowTaxiList()
		return true
	}

	// process text with context
	if ctx != "" {
		switch ctx {
		case CONTEXT_KINO:
			proc.ShowFilms(text)
			setUserContext(uid, CONTEXT_KINO)
		}

		return true
	}

	// catch commands if empty context
	if ctx == "" && (strings.ToLower(text) == "kino" || strings.ToLower(text) == "/kino") {
		setUserContext(uid, CONTEXT_KINO)
		proc.ShowKinoPlaces()
		return true
	}

	return false
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
