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

package smarty_test

import (
    "testing"

    "github.com/NEXXTLAB/go-smarty-reader/smarty"
)

// Test if the decryption is possible with the provided key and pre-recorded telegram (found in Smarty_test.go)
func TestDecryption(t *testing.T) {
    smartyObj := smarty.NewDecryptor(string(key))
    iv := append(systemTitle, frameCounter...)
    cipher := append(payload, gcmTag...)
    plainText, ok := smartyObj.Decrypt(iv, cipher)

    if ok {
        t.Logf("Decryption success: \n%s\n", plainText)
    } else {
        t.Error("Decryption failed!")
    }
}
