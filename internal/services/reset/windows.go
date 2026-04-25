package reset

import (
	"fmt"
	"strings"

	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func resetWindows(h string, products []string) {
	cleanDir(filepath.Join(h, "AppData", "Roaming", "JetBrains"), products)
	cleanDir(filepath.Join(h, "AppData", "Local", "JetBrains"), products)
	cleanRegistry(products)
}

func cleanRegistry(products []string) {
	const basePath = `Software\JavaSoft\Prefs\jetbrains`

	for _, product := range products {
		keyPath := basePath + `\` + strings.ToLower(product)
		err := deleteKeyRecursive(registry.CURRENT_USER, keyPath)
		if err != nil && err != registry.ErrNotExist {
			fmt.Printf("⚠️  Registry: could not clean %s: %v\n", product, err)
		} else if err == nil {
			fmt.Printf("🧹 Registry: cleaned %s\n", product)
		}
	}
}

// deleteKeyRecursive рекурсивно удаляет ключ реестра вместе со всеми подключами.
func deleteKeyRecursive(root registry.Key, path string) error {
	key, err := registry.OpenKey(root, path, registry.READ)
	if err != nil {
		return err
	}

	// Получаем все подключи и удаляем их рекурсивно
	subkeys, err := key.ReadSubKeyNames(-1)
	key.Close()
	if err != nil {
		return err
	}

	for _, sub := range subkeys {
		if err := deleteKeyRecursive(root, path+`\`+sub); err != nil {
			return err
		}
	}

	return registry.DeleteKey(root, path)
}
