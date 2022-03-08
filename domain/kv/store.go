package kv

import (
	"fmt"
	"os"
	"strings"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/pkg/bolt"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Store struct {
	log  *logger.CLILogger
	bolt *bolt.Bolt
}

func NewStore(log *logger.CLILogger) *Store {
	kv := &Store{log: log}

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

func (s *Store) Close() error {
	return s.bolt.Close()
}

func (s *Store) Set(path, value string) error {
	bucket, key, err := s.splitBucketAndKey(path)
	if err != nil {
		return err
	}

	err = s.bolt.CreateBucket([]byte(bucket))
	if err != nil {
		return fmt.Errorf("error creating bucket %s: %v", bucket, err)
	}

	if value == "rand32" {
		value, err = localutils.Base64EncodedRandomNumber(32)
		if err != nil {
			return err
		}
	}

	err = s.bolt.Put([]byte(bucket), []byte(key), []byte(value))
	if err != nil {
		return fmt.Errorf("error adding key %s to bucket %s: %v", key, bucket, err)
	}

	return nil
}

func (s *Store) Get(path string) (string, error) {
	bucket, key, err := s.splitBucketAndKey(path)
	if err != nil {
		return "", err
	}

	value, err := s.bolt.Get([]byte(bucket), []byte(key))
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (s *Store) Delete(path string) error {
	bucket, key, err := s.splitBucketAndKey(path)
	if err != nil {
		return err
	}

	return s.bolt.Delete([]byte(bucket), []byte(key))
}

func (s *Store) GetAll() (map[string]string, error) {
	store := map[string]string{}
	buckets, err := s.bolt.ListBuckets()
	if err != nil {
		return store, err
	}

	for _, bucket := range buckets {
		bucketStore, err := s.bolt.FindAll([]byte(bucket))
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

func (s *Store) splitBucketAndKey(path string) (string, string, error) {
	splits := strings.Split(path, ".")
	if len(splits) != 2 {
		return "", "", fmt.Errorf("was expecting bucket_name.key_name, but recieved %s", path)
	}
	return splits[0], splits[1], nil
}

func (s *Store) Display(log *logger.CLILogger, action string) {
	store, err := s.GetAll()
	if err != nil {
		log.Error("kv").Msgf("could not read from the database: %v", err)
		_ = s.Close()
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
	_ = s.Close()
}

func (s *Store) SetUserValue(path string) {
	userValue := localutils.Question(s.log, "enter", "please enter a value")
	userValue = strings.TrimSuffix(userValue, "\n")
	err := s.Set(path, userValue)
	if err != nil {
		s.log.Error("kv").Msgf("could not set the value: %v", err)
	}
}

func (s *Store) SetValue(path, value string) {
	err := s.Set(path, value)
	if err != nil {
		s.log.Error("kv").Msgf("could not set the value: %v", err)
	}

	value, err = s.Get(path)
	if err != nil {
		s.log.Error("kv").Msgf("could not get value: %v", err)
	}

	if err == nil {
		s.log.Info("kv").Msgf("%s: %s", logger.Cyan(path), logger.GreenBold(value))
	}
}

func (s *Store) DeleteValue(path string) {
	err := s.Delete(path)
	if err != nil {
		s.log.Error("kv").Msgf("could not delete the value: %v", err)
	}
}

func (s *Store) GetValue(path string) {
	value, err := s.Get(path)
	if err != nil {
		s.log.Error("kv").Msgf("could not get value: %v", err)
	}

	if err == nil {
		s.log.Info("kv").Msgf("%s: %s", logger.Cyan(path), logger.GreenBold(value))
	}
}
