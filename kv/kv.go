package kv

import (
	"fmt"
	"os"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/pkg/bolt"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type KV struct {
	log  *logger.CLILogger
	bolt *bolt.Bolt
}

func NewKV(log *logger.CLILogger) *KV {
	kv := &KV{log: log}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("KV").Msgf("could not find user's home directory: %v", err)
		os.Exit(1)
	}
	golemDBPath := fmt.Sprintf("%s/.golem/golem.db", homeDir)

	kv.bolt, err = bolt.NewBolt(golemDBPath)
	if err != nil {
		log.Fatal("KV").Msgf("coudl not open database: %v", err)
		os.Exit(1)
	}

	return kv
}

func (kv *KV) Close() error {
	return kv.bolt.Close()
}

func (kv *KV) Set(path, value string) error {
	bucket, key, err := kv.splitBucketAndKey(path)
	if err != nil {
		return err
	}

	err = kv.bolt.CreateBucket([]byte(bucket))
	if err != nil {
		return fmt.Errorf("error creating bucket %s: %v", bucket, err)
	}

	if value == "rand32" {
		value, err = localutils.Base64EncodedRandomNumber(32)
		if err != nil {
			return err
		}
	}

	err = kv.bolt.Put([]byte(bucket), []byte(key), []byte(value))
	if err != nil {
		return fmt.Errorf("error adding key %s to bucket %s: %v", key, bucket, err)
	}

	return nil
}

func (kv *KV) Get(path string) (string, error) {
	bucket, key, err := kv.splitBucketAndKey(path)
	if err != nil {
		return "", err
	}

	value, err := kv.bolt.Get([]byte(bucket), []byte(key))
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (kv *KV) Delete(path string) error {
	bucket, key, err := kv.splitBucketAndKey(path)
	if err != nil {
		return err
	}

	return kv.bolt.Delete([]byte(bucket), []byte(key))
}

func (kv *KV) GetAll() (map[string]string, error) {
	store := map[string]string{}
	buckets, err := kv.bolt.ListBuckets()
	if err != nil {
		return store, err
	}

	for _, bucket := range buckets {
		bucketStore, err := kv.bolt.FindAll([]byte(bucket))
		if err != nil {
			return store, err
		}

		for key, value := range bucketStore {
			storeKey := fmt.Sprintf("%s.%s", bucket, key)
			store[storeKey] = string(value)
		}
	}
	return store, nil
}

func (kv *KV) splitBucketAndKey(path string) (string, string, error) {
	splits := strings.Split(path, ".")
	if len(splits) != 2 {
		return "", "", fmt.Errorf("was expecting bucket_name.key_name, but recieved %s", path)
	}
	return splits[0], splits[1], nil
}

func (kv *KV) Display(log *logger.CLILogger, action string) {
	store, err := kv.GetAll()
	if err != nil {
		log.Error("kv").Msgf("could not read from the database: %v", err)
		_ = kv.Close()
		return
	}

	kvLog := func(log *logger.CLILogger, key, value string) {
		if strings.Contains(key, "secret") || strings.Contains(key, "password") {
			value = value[:2] + "************"
		}
		log.Info("kv").Msgf("%s: %s", logger.Cyan(key), logger.GreenBold(value))
	}

	for key, value := range store {
		switch action {
		case "":
			kvLog(log, key, value)
		default:
			if strings.Contains(key, action) {
				kvLog(log, key, value)
			}
		}
	}
	_ = kv.Close()
}

func (kv *KV) SetUserValue(path string) {
	userValue := localutils.Question(kv.log, "enter", "please enter a value")
	userValue = strings.TrimSuffix(userValue, "\n")
	err := kv.Set(path, userValue)
	if err != nil {
		kv.log.Error("kv").Msgf("could not set the value: %v", err)
	}
}

func (kv *KV) SetValue(path, value string) {
	err := kv.Set(path, value)
	if err != nil {
		kv.log.Error("kv").Msgf("could not set the value: %v", err)
	}

	value, err = kv.Get(path)
	if err != nil {
		kv.log.Error("kv").Msgf("could not get value: %v", err)
	}

	if err == nil {
		kv.log.Info("kv").Msgf("%s: %s", logger.Cyan(path), logger.GreenBold(value))
	}
}

func (kv *KV) DeleteValue(path string) {
	err := kv.Delete(path)
	if err != nil {
		kv.log.Error("kv").Msgf("could not delete the value: %v", err)
	}
}

func (kv *KV) GetValue(path string) {
	value, err := kv.Get(path)
	if err != nil {
		kv.log.Error("kv").Msgf("could not get value: %v", err)
	}

	if err == nil {
		kv.log.Info("kv").Msgf("%s: %s", logger.Cyan(path), logger.GreenBold(value))
	}
}

// golem kv bucket1.key_name set
// -> ask for user input

// golem kv bucket1.key_name rand32 (or rand64 or rand16)
// -> set a random value and display it

// golem kv bucket.key_name
// -> show the value

// golem kv bucket.key_name delete
// -> deletes the value

// vars {
//     GOLEM_CONFIG_NAME = "@golem.kv.bucket1.key_name"
// }

// recipe "name" "local" {
//     // adds it to the env before running
//     env {
//         GOLEM_CONFIG_NAME = "@golem.kv.bucket1.key_name"
//     }
// }
