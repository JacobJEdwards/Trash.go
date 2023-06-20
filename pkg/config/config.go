package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	configPath = filepath.Join(os.Getenv("HOME"), "/.config/trash/config.json")
	configDir  = filepath.Join(os.Getenv("HOME"), "/.config/trash")
)

type Config struct {
	TrashDir string `json:"trash_dir"`
	Logfile  string `json:"logfile"`
}

func LoadConfig() (*Config, error) {
	var c Config

	file, err := os.Open(configPath)

	if err != nil {

		err := c.createConfig()

		if err != nil {
			fmt.Println(err)
			return &c, err
		}

		return &c, nil
	}

	defer file.Close()
	byteConfig, err := io.ReadAll(file)

	if err != nil {
		return &c, err
	}

	json.Unmarshal(byteConfig, &c)

	if c.TrashDir == "" {
		c.TrashDir = filepath.Join(os.Getenv("HOME"), "/.go-trash")
		c.SaveConfig()
	}

	if c.Logfile == "" {
		c.Logfile = filepath.Join(os.Getenv("HOME"), "/.go-trash.log")
		c.SaveConfig()
	}

	err = c.ValidateConfig()

	if err != nil {
		return c.ResetConfig()
	}

	return &c, nil
}

func (c *Config) createConfig() error {
	err := createConfigDirectory()
	if err != nil {
		return err
	}

	if c.TrashDir == "" {
		c.TrashDir = filepath.Join(os.Getenv("HOME"), "/.go-trash")
		c.SaveConfig()
	}

	if c.Logfile == "" {
		c.Logfile = filepath.Join(os.Getenv("HOME"), "/.go-trash.log")
		c.SaveConfig()
	}

	err = os.MkdirAll(c.TrashDir, 0755)
	if err != nil {
		return fmt.Errorf("Error creating trash directory: %v", err)
	}

	err = createConfigFile(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SaveConfig() error {
	file, err := os.Open(configPath)

	if err != nil {
		return err
	}

	defer file.Close()

	byteConfig, err := json.MarshalIndent(c, "", "  ")

	if err != nil {
		return fmt.Errorf("Error creating config file: %v", err)
	}

	err = ioutil.WriteFile(configPath, byteConfig, 0644)

	if err != nil {
		return fmt.Errorf("Error writing config file: %v", err)
	}

	return nil
}

func createConfigDirectory() error {
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("Error creating config directory: %v", err)
	}
	return nil
}

func createConfigFile(c *Config) error {
	byteConfig, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("Error creating config file: %v", err)
	}
	err = ioutil.WriteFile(configPath, byteConfig, 0644)
	if err != nil {
		return fmt.Errorf("Error writing config file: %v", err)
	}
	return nil
}

func (c *Config) ResetConfig() (*Config, error) {
	err := os.Remove(configPath)

	if err != nil {
		return c, fmt.Errorf("Error removing config file: %v", err)
	}
	err = c.createConfig()
	if err != nil {
		return c, err
	}

	return LoadConfig()
}

// ValidateConfig checks if the configuration values are valid.
func (c *Config) ValidateConfig() error {
	// Check if the TrashDir is a valid directory path
	if err := validateDirectory(c.TrashDir); err != nil {
		return err
	}

	// Check if the Logfile is a valid file path
	if err := validateFile(c.Logfile); err != nil {
		return err
	}

	return nil
}

// validateDirectory checks if the provided directory path exists and is accessible.
func validateDirectory(directory string) error {
	stat, err := os.Stat(directory)
	if err != nil {
		return errors.New("Invalid directory path: " + directory)
	}
	if !stat.IsDir() {
		return errors.New("Path is not a directory: " + directory)
	}
	return nil
}

// validateFile checks if the provided file path exists and is accessible.
func validateFile(filePath string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return errors.New("Invalid file path: " + filePath)
	}
	if stat.IsDir() {
		return errors.New("Path is a directory, expected a file: " + filePath)
	}
	return nil
}

// SanitizeConfig sanitizes the configuration values, if necessary.
func (c *Config) SanitizeConfig() {
	// Sanitize the TrashDir by removing any trailing slashes
	c.TrashDir = strings.TrimSuffix(c.TrashDir, "/")
}
