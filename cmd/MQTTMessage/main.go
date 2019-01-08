/*
 * This file is part of go-smarty-reader
 *
 * go-smarty-reader is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * go-smarty-reader is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with go-smarty-reader. If not, see <https://www.gnu.org/licenses/>.
 */

/*
   Example on how to use the MQTT wrapper
*/

package main

import (
	"github.com/NEXXTLAB/go-smarty-reader/cmd/util"
)

func main() {

	// Function defined in cmd/util/CommonFlagParsing.go
	util.StartupFlagParsing()

	// MQTT Setup extracted in a separate function.
	// Functions defined in cmd/util/CommonMqttSetup.go
	client := util.MqttSetup(util.GetHostname())

	// Publish "Hello" to the ""nexxtlab/dev/smarty/go/<hostname>/World" topic, without unit.
	client.Publish("World", "Hello", "", false, false)

	// Close the MQTT connection
	client.Disconnect(250)
}
