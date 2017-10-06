package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	
	"github.com/leoloobeek/GoGreen/lib"
)

func main() {
	conf := flag.String("config", "", "Config file to read. Default ./config.json.")
	flag.Parse()
	config := parseConfig(*conf)
	if config == nil {
		return
	}
	payloadWriter := setupPayloadWriter(config.Language)

	if payloadWriter == nil {
		fmt.Println("[!] Specified language in config.json not supported")
		return
	}

	// get payload from config or file
	payload := config.Payload
	if payload == "" {
		pFile, err := readFile(config.PayloadPath)
		if err != nil {
			fmt.Printf("[!] Error reading payload from %s\n", config.PayloadPath)
			return
		}
		payload = string(pFile)
		fmt.Println(payload)
		if payload == "" {
			fmt.Printf("[!] Payload file %s was empty", config.PayloadPath)
			return
		}
	}

	// generate payload hash and key
	payloadHash := lib.GenerateSHA512(payload)
	if config.StartDir == "" {
		config.PathKey = ""
	}
	key, envVarOrder := buildKey(config.EnvVars, config.PathKey)

	// read in "master" template
	result, err := readFile(payloadWriter.MasterTemplate)
	if err != nil {
		fmt.Printf("[!] Unable to read %s\n", (payloadWriter.MasterTemplate))
		return
	}

	// if a StartDir was specified in config.json, add in directory walking
	result, err = directoryCode(result, config.StartDir, payloadWriter)
	if err != nil {
		return
	}

	// set up env vars code
	result, err = envVarCode(result, envVarOrder, payloadWriter)
	if err != nil {
		return
	}

	// TODO: handle multiple SHA512 iterations
	keyHash := lib.GenerateSHA512(strings.ToLower(key))[:32]
	text, iv, err := lib.AESEncrypt([]byte(keyHash), []byte(payload))

	if err != nil {
		fmt.Printf("Error received encrypting: %s", err)
		return
	}

	result = bytes.Replace(result, []byte("~AESIVBASE64~"), []byte(iv), 1)
	result = bytes.Replace(result, []byte("~ENCRYPTEDBASE64~"), []byte(text), 1)
	result = bytes.Replace(result, []byte("~PAYLOADHASH~"), []byte(payloadHash), 1)

	writeFile(payloadWriter.Outfile, result)
	fmt.Println("[*] Output File:        " + payloadWriter.Outfile)
	fmt.Println("[*] Environmental Keys: " + key)
	fmt.Println("[*] Decryption Key:     " + keyHash)

}

// Config defines structure to be read from config.json
type Config struct {
	Language    string
	StartDir    string
	PathKey     string
	EnvVars     map[string]string
	Payload     string
	PayloadPath string
}

// PayloadWriter for handling what to write
type PayloadWriter struct {
	Language       string
	MasterTemplate string
	DirTemplate    string
	EnvKeyTemplate string
	EnvKeyCode     string
	Outfile        string
}

func setupPayloadWriter(lang string) *PayloadWriter {
	switch lang {
	case "vbscript":
		return &PayloadWriter{
			Language:       "vbscript",
			MasterTemplate: "data/vbscript/base.vbs",
			DirTemplate:    "data/vbscript/directory.vbs",
			EnvKeyTemplate: "data/vbscript/envkey.vbs",
			EnvKeyCode:     "oEnv(\"%s\")",
			Outfile:        "payload.vbs"}
	case "jscript":
		return &PayloadWriter{
			Language:       "jscript",
			MasterTemplate: "data/jscript/base.js",
			DirTemplate:    "data/jscript/directory.js",
			EnvKeyTemplate: "data/jscript/envkey.js",
			EnvKeyCode:     "oEnv(\"%s\")",
			Outfile:        "payload.js"}
	case "powershell":
		return &PayloadWriter{
			Language:       "powershell",
			MasterTemplate: "data/powershell/base.ps1",
			DirTemplate:    "data/powershell/directory.ps1",
			EnvKeyTemplate: "data/powershell/envkey.ps1",
			EnvKeyCode:     "$oEnv.Invoke(\"%s\")",
			Outfile:        "payload.ps1"}
	default:
		return nil
	}
}

func parseConfig(configPath string) *Config {
	if configPath == "" {
		configPath = "config.json"
	}
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("[!] Error reading config.json")
		return nil
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("[!] Error parsing config.json:", err)
		return nil
	}
	return &config
}

func readFile(path string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("[!] readFile: Error reading file")
	}
	return fileBytes, nil
}

func writeFile(filename string, contents []byte) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.Write(contents)
	if err != nil {
		panic(err)
	}
}

func getFilenames(lang string) (string, string) {
	switch lang {
	case "vbscript":
		return "data/keying.vbs", "payload.vbs"
	case "jscript":
		return "data/keying.js", "payload.js"
	default:
		return "", ""
	}
}

func directoryCode(text []byte, startDir string, pw *PayloadWriter) ([]byte, error) {
	if startDir != "" {
		contents, err := readFile(pw.DirTemplate)
		if err != nil {
			fmt.Printf("[!] Error reading %s\n", pw.DirTemplate)
			return nil, err
		}

		// JScript needs escaped slashes, VBS does not
		if pw.Language == "jscript" {
			startDir = strings.Replace(startDir, "\\", "\\\\", -1)
		}

		result := bytes.Replace(text, []byte("~WALKOS~"), contents, 1)
		result = bytes.Replace(result, []byte("~STARTDIR~"), []byte(startDir), 1)

		return result, nil
	}
	result := bytes.Replace(text, []byte("~WALKOS~"), []byte(""), 1)
	return result, nil
}

// Set up the environmental vars code within the script
// TODO: This got overly complex, need to simplify
func envVarCode(text []byte, envVarCode []string, pw *PayloadWriter) ([]byte, error) {
	if len(envVarCode) > 0 {
		contents, err := readFile(pw.EnvKeyTemplate)
		if err != nil {
			fmt.Printf("[!] Error reading %s\n", pw.EnvKeyTemplate)
			return nil, err
		}

		result := bytes.Replace(text, []byte("~ENVVAR~"), contents, 1)

		envs := fmt.Sprintf(pw.EnvKeyCode, envVarCode[0])
		if len(envVarCode) > 1 {
			for _, envVar := range envVarCode[1:] {
				envs += fmt.Sprintf(", "+pw.EnvKeyCode, envVar)
			}
		}
		result = bytes.Replace(result, []byte("~ENVVARS~"), []byte(envs), 1)

		return result, nil
	}
	result := bytes.Replace(text, []byte("~ENVVAR~"), []byte(""), 1)
	return result, nil
}

// buildKey builds the final key with env vars, path, etc. and also returns
// the order of the env vars (maps don't always come out in same order)
// key = [env var 1][env var 2][env var N][file path]
func buildKey(envVars map[string]string, filePath string) (string, []string) {
	result := ""
	envVarOrder := make([]string, len(envVars))

	i := 0
	for k, v := range envVars {
		result += v
		envVarOrder[i] = k
		i++
	}
	result += filePath

	return result, envVarOrder
}
