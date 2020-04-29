/*
Copyright © 2020 XP-1000 <xp-1000@hotmail.fr>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"fmt"
    "github.com/xp-1000/gowck/pkg/httpcheck"
)

func run(msg, loghost string) error {
	fmt.Printf("My message is '%s' and I'm logging it to '%s'\n", msg, loghost)
	httpcheck.monitor()
	return nil
}
