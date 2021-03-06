// Copyright 2014 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)

package main

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestPitDisplay(t *testing.T) {
	clearDb()
	defer clearDb()
	var err error
	db, err = OpenDatabase(testDbPath)
	assert.Nil(t, err)
	defer db.Close()
	eventSettings, _ = db.GetEventSettings()

	recorder := getHttpResponse("/displays/pit")
	assert.Equal(t, 200, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Pit Display - Untitled Event - Cheesy Arena")
}

func TestPitDisplayWebsocket(t *testing.T) {
	clearDb()
	defer clearDb()
	var err error
	db, err = OpenDatabase(testDbPath)
	assert.Nil(t, err)
	defer db.Close()
	eventSettings, _ = db.GetEventSettings()
	mainArena.Setup()

	server, wsUrl := startTestServer()
	defer server.Close()
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl+"/displays/pit/websocket", nil)
	assert.Nil(t, err)
	defer conn.Close()
	ws := &Websocket{conn, new(sync.Mutex)}

	// Check forced reloading as that is the only purpose the pit websocket serves.
	recorder := getHttpResponse("/setup/field/reload_displays")
	assert.Equal(t, 302, recorder.Code)
	readWebsocketType(t, ws, "reload")
}
