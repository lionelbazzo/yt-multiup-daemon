/*
 * Copyright (c) 2022 Michaël COLL.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package banner

import (
	"fmt"

	"github.com/fatih/color"
)

type Mode int

const (
	banner = `
          ________
         /\       \
        /  \       \
       /    \       \
      /      \_______\
      \      /       /
    ___\    /   ____/___
   /\   \  /   /\       \
  /  \   \/___/  \       \
 /    \       \   \       \
/      \_______\   \_______\
\      /       /   /       /   %s ---- %s
 \    /       /   /       /    =====<< %s >>=====
  \  /       /\  /       /
   \/_______/  \/_______/      %s

`

	Serve Mode = 0
	Index      = 1
)

func Print(version string, owner string, mode Mode) {
	var modeStr string

	switch mode {
	case Serve:
		modeStr = "serve mode"
	case Index:
		modeStr = "index mode"
	}

	var ownerStr string
	if owner != "" {
		ownerStr = fmt.Sprintf("%s%s%s",
			color.WhiteString("--|"),
			color.HiWhiteString(owner),
			color.WhiteString("|--"))
	}

	fmt.Printf(banner,
		color.BlueString("gallery daemon"),
		color.WhiteString(version),
		color.CyanString(modeStr),
		ownerStr)
}
