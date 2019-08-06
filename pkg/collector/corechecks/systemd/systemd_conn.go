// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build systemd

package systemd

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	godbus "github.com/godbus/dbus"
	"os"
	"strconv"
)

//
//// New establishes a connection to any available bus and authenticates.
//// Callers should call Close() when done with the connection.
//func New() (*Conn, error) {
//	conn, err := dbus.NewSystemConnection()
//	if err != nil && os.Geteuid() == 0 {
//		return NewSystemdConnection()
//	}
//	return conn, err
//}

// NewSystemdConnection establishes a private, direct connection to systemd.
// This can be used for communicating with systemd without a dbus daemon.
// Callers should call Close() when done with the connection.
func NewSystemdConnection(address string) (*dbus.Conn, error) {
	return dbus.NewConnection(func() (*godbus.Conn, error) {
		// We skip Hello when talking directly to systemd.
		return dbusAuthConnection(func() (*godbus.Conn, error) {
			return godbus.Dial(fmt.Sprintf("unix:path=%s", address))
		})
	})
}

//"unix:path=/tmp/run/systemd/private"

func dbusAuthConnection(createBus func() (*godbus.Conn, error)) (*godbus.Conn, error) {
	conn, err := createBus()
	if err != nil {
		return nil, err
	}

	// Only use EXTERNAL method, and hardcode the uid (not username)
	// to avoid a username lookup (which requires a dynamically linked
	// libc)
	methods := []godbus.Auth{godbus.AuthExternal(strconv.Itoa(os.Getuid()))}

	err = conn.Auth(methods)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
