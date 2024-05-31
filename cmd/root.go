package cmd

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/spf13/cobra"
)

var tfa *TwoFactor

var rootCmd = &cobra.Command{
	Use:   "2fa",
	Short: "2fa",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(delCmd)
	var err error
	tfa, err = New()
	if err != nil {
		log.Fatalln("could not create 2fa", err)
	}

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func GenTOTP(key string, t time.Time, digits int, timeStep int64) (string, error) {
	decode, err := base32.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("failed to decode token %s: %v", key, err)
	}

	t1 := t.Unix()
	C := t1 / timeStep

	cBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(cBuf, uint64(C))
	mac := hmac.New(sha1.New, decode)
	_, _ = mac.Write(cBuf)
	H := mac.Sum(nil)

	offset := H[len(H)-1] & 0xf
	value := int64(((int(H[offset]) & 0x7f) << 24) |
		(int(H[offset+1]&0xff) << 16) |
		(int(H[offset+2]&0xff) << 8) |
		(int(H[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(digits)))

	f := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(f, mod), nil
}
