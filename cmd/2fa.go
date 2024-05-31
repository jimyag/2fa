package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const secretKey = "018fcf3d2bf3745ca38d0360d6816e68"

type Key struct {
	Key string `json:"key"`
}
type TwoFactor struct {
	Keys    map[string]Key `json:"keys"`
	cfgFile string         `json:"-"`
}

func New() (*TwoFactor, error) {
	tf := &TwoFactor{
		Keys: map[string]Key{
			"": {
				Key: secretKey,
			},
		},
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	configDir := fmt.Sprintf("%s/.config/2fa", homeDir)
	if err = os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create config directory %s : %w", configDir, err)
	}

	tf.cfgFile = fmt.Sprintf("%s/.2fa.json", configDir)
	_, err = os.Stat(tf.cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			if err = tf.Write(); err != nil {
				return nil, fmt.Errorf("could not write config file %s : %w", tf.cfgFile, err)
			}
		} else {
			return nil, fmt.Errorf("could not stat config file %s : %w", tf.cfgFile, err)
		}
	}
	err = tf.Load()
	return tf, err

}

func (tf *TwoFactor) Add(name, key string) error {
	key = strings.ToUpper(key)
	tf.Keys[name] = Key{key}
	return tf.Write()
}

func (tf *TwoFactor) List() map[string]Key {
	list := make(map[string]Key)
	for k, v := range tf.Keys {
		list[k] = v
	}
	return list
}

func (tf *TwoFactor) Get(name string) string {
	if key, ok := tf.Keys[name]; ok {
		return key.Key
	}
	return ""
}

func (tf *TwoFactor) Remove(name string) error {
	delete(tf.Keys, name)
	return tf.Write()
}

func (tf *TwoFactor) Write() error {
	file, err := os.OpenFile(tf.cfgFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("could not open config file", err)
		return err
	}
	defer file.Close()
	plaintext, err := json.Marshal(tf.Keys)
	if err != nil {
		return err
	}

	cipherText, err := tf.EncryptMessage(plaintext)
	if err != nil {
		return err
	}
	_, err = file.WriteString(cipherText)
	return err
}

func (tf *TwoFactor) Load() error {
	file, err := os.Open(tf.cfgFile)
	if err != nil {
		return err
	}
	defer file.Close()

	cipherText, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	data, err := tf.DecryptMessage(cipherText)
	if err != nil {
		fmt.Println(err)
		return json.Unmarshal(cipherText, &tf.Keys)
	}
	return json.Unmarshal([]byte(data), &tf.Keys)
}

func (tf *TwoFactor) EncryptMessage(message []byte) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(message))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], message)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (tf *TwoFactor) DecryptMessage(message []byte) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(string(message))
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
