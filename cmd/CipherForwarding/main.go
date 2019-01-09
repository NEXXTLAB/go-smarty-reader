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
   Example on how to use the CipherForwarder
*/

package main

import (
    "github.com/NEXXTLAB/go-smarty-reader/cmd/util"
    "github.com/NEXXTLAB/go-smarty-reader/smarty"
)

func main() {

    // Function defined in cmd/util/CommonFlagParsing.go
    flags := util.StartupFlagParsing()

    // Create a new smarty reader which will not decrypt the telegrams,
    // but return the initial value and the cipher text
    // The serial connection is established right away
    // smartyObj is the object you may invoke methods on
    smartyObj := smarty.NewCipherForwarder(*flags.Device)

    // Print 100 initial value / cipher tuples
    for counter := 0; counter < 100; counter++ {
        // Wait and get the next telegram,
        // return the initial value and cipher text
        iv, cipher, gcmTag := smartyObj.GetTelegram()
        // Print as console output
        println(string(iv))
        println(string(cipher))
        println(string(gcmTag))
        println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
    }
    // After use, remember to close to serial port!
    smartyObj.Disconnect()
}
