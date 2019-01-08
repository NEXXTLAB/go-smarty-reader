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
   CommonFlagParsing holds the common command-line flag operations in one file to avoid code duplicates.
*/

package util

import (
	"flag"

	"github.com/golang/glog"
)

func StartupFlagParsing() (deviceFlag, keyFlag *string) {

	const VERSION = "1.0.1"

	// Define flags
	deviceFlag = flag.String("device", "", "Serial device to read P1 data from.")
	keyFlag = flag.String("key", "", "Decryption Key to use.")

	flag.Parse()

	// Print version info and warnings if either the device- or keyFlag is missing
	glog.Infoln("Smarty Reader " + VERSION)
	if *deviceFlag == "" {
		glog.Warningln("Serial device parameter missing.\n\t" +
			"This program instance will not be able to access any serial devices.")
	}
	if *keyFlag == "" {
		glog.Warningln("No decryption key found.\n\t" +
			"Telegrams can not be decrypted in this program instance.")
	}
	return deviceFlag, keyFlag
}
