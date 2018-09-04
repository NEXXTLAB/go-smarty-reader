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
   Example on how to use the OnlineDecryptor
*/

package main

import (
	"github.com/NEXXTLAB/go-smarty-reader/cmd/util"
	"github.com/NEXXTLAB/go-smarty-reader/smarty"
)

func main() {

	// Function defined in cmd/util/CommonFlagParsing.go
	deviceFlag, keyFlag := util.StartupFlagParsing()

	// Create a new smarty reader which will decrypt the telegrams after reading them
	// The serial connection is established right away
	// smartyObj is the object you may invoke methods on
	smartyObj := smarty.NewOnlineDecryptor(*deviceFlag, *keyFlag)

	// Read until 100 telegrams could be successfully decrypted
	for telegramCounter := 0; telegramCounter < 100; {
		// Wait, get and decrypt the next telegram
		plainText, ok := smartyObj.GetTelegram()
		// If the decryption was successful, print the payload to the console
		if ok {
			println(string(plainText))
			println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			telegramCounter++
		}
	}
	// After use, remember to close to serial port!
	smartyObj.Disconnect()
}
