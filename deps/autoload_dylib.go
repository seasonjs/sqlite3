// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package deps

import (
	"fmt"
	"os"
)

func DumpSqliteLibrary() (string, error) {
	file, err := os.CreateTemp("", libName)
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}

	defer file.Close()

	if err := os.WriteFile(file.Name(), getDylib(), os.ModePerm); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}

	return file.Name(), nil
}
