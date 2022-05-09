package src

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var cf *Configuration

type Configuration struct {
	ServerAddress string
	ServerPort    int
	RetryTimes    int
	//database
	DBHost       string
	DBName       string
	DBCollection string
	DBUsername   string
	DBPassword   string

	//storage
	StorageRegion          string
	StorageUrl             string
	StorageBucketName      string
	StorageAccessKeyID     string
	StorageAccessKeySecret string
	//profile
	PostImageProfiles   []int
	AvatarImageProfiles []int
}

func GetConfig() *Configuration {
	return cf
}

func InitFromFile(filePathStr string) {

	env := os.Getenv("GO_ENV")
	if env == "" {
		os.Setenv("GO_ENV", "dev")
		env = "dev"
	}

	viper.SetConfigFile(filePathStr)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config file not found: %v", err)
	} else {
		cf = &Configuration{
			ServerAddress:          viper.GetString(env + ".server_address"),
			ServerPort:             viper.GetInt(env + ".server_port"),
			RetryTimes:             viper.GetInt(env + ".retry_times"),
			DBHost:                 viper.GetString(env + ".db_host"),
			DBName:                 viper.GetString(env + ".db_name"),
			DBCollection:           viper.GetString(env + ".db_collection"),
			DBUsername:             viper.GetString(env + ".db_username"),
			DBPassword:             viper.GetString(env + ".db_password"),
			StorageRegion:          viper.GetString(env + ".region"),
			StorageUrl:             viper.GetString(env + ".storage_url"),
			StorageBucketName:      viper.GetString(env + ".bucket_name"),
			StorageAccessKeyID:     viper.GetString(env + ".access_key_id"),
			StorageAccessKeySecret: viper.GetString(env + ".access_key_secret"),
		}

		imagePosts := strings.Split(viper.GetString(env+".post_image_profiles"), ",")
		for _, size := range imagePosts {
			iSize, _ := strconv.Atoi(size)
			cf.PostImageProfiles = append(cf.PostImageProfiles, iSize)
		}
		imageAvatars := strings.Split(viper.GetString(env+".avatar_image_profiles"), ",")
		for _, size := range imageAvatars {
			iSize, _ := strconv.Atoi(size)
			cf.AvatarImageProfiles = append(cf.AvatarImageProfiles, iSize)
		}
		log.Println(viper.ConfigFileUsed())
		log.Printf("Config %+v", *cf)
	}
}
