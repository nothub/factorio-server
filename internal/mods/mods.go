package mods

import (
	"encoding/json"
	"fmt"
	"github.com/nothub/factorio-server/internal/config"
	factorioCom "github.com/nothub/factorio-server/internal/factorio.com"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func Sync() error {

	err := os.MkdirAll("mods", 0755)
	if err != nil {
		return fmt.Errorf("failed to create mods directory: %w", err)
	}

	latestFileNames, err := downloadMods()
	if err != nil {
		return err
	}

	err = cleanupMods(latestFileNames)
	if err != nil {
		return err
	}

	err = generateListFile()
	if err != nil {
		return err
	}

	return nil
}

func downloadMods() ([]string, error) {

	var latestFileNames []string

	for _, modId := range config.Loaded.Mods {

		if modId == "base" {
			continue
		}

		mod, err := factorioCom.GetModInfo(modId)
		if err != nil {
			return nil, fmt.Errorf("failed to get mod info for %s: %w", modId, err)
		}

		fileName, err := factorioCom.DownloadMod(mod)
		if err != nil {
			return nil, fmt.Errorf("failed to download mod %s: %w", mod.Name, err)
		}

		latestFileNames = append(latestFileNames, fileName)
	}

	return latestFileNames, nil
}

func cleanupMods(latestFileNames []string) error {

	entries, err := os.ReadDir("mods")
	if err != nil {
		return fmt.Errorf("failed reading mods dir: %w", err)
	}

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		if entry.Name() == "mod-list.json" {
			continue
		}

		if entry.Name() == "mod-settings.dat" {
			continue
		}

		if slices.Contains(latestFileNames, entry.Name()) {
			continue
		}

		log.Printf("Removing %s (old)\n", entry.Name())
		err = os.Remove(filepath.Join("mods", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %w", entry.Name(), err)
		}

	}

	return nil
}

func generateListFile() error {

	type ModListData struct {
		Mods []struct {
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
		} `json:"mods"`
	}

	var modList ModListData
	modList.Mods = append(modList.Mods, struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}{
		Name:    "base",
		Enabled: true,
	})

	for _, mod := range config.Loaded.Mods {
		modList.Mods = append(modList.Mods, struct {
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
		}{
			Name:    mod,
			Enabled: true,
		})
	}

	modListData, err := json.MarshalIndent(modList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mod-list.json: %w", err)
	}

	err = os.WriteFile(filepath.Join("mods", "mod-list.json"), modListData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write mod-list.json: %w", err)
	}

	return nil
}
