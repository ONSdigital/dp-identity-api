package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GroupsFilename       string `envconfig:"GROUPS_FILENAME" required:"true"`
	GroupUsersFilename   string `envconfig:"GROUPUSERS_FILENAME" required:"true"`
	UserFileName         string `envconfig:"FILENAME" required:"true"`
	S3Bucket             string `envconfig:"S3_BUCKET" required:"true"`
	S3BaseDir            string `envconfig:"S3_BASE_DIR" required:"true"`
	S3Region             string `envconfig:"S3_REGION" required:"true"`
	AWSCognitoUserPoolID string `envconfig:"USER_POOL_ID" required:"true"`
}

func (c Config) GetS3UsersFilePath() string {
	return fmt.Sprintf("%s%s", c.S3BaseDir, c.UserFileName)
}

func GetConfig() *Config {
	conf := &Config{}
	envconfig.Process("", conf)
	return conf
}

func (c Config) GetS3GroupsFilePath() string {
	return fmt.Sprintf("%s%s", c.S3BaseDir, c.GroupsFilename)
}

func (c Config) GetS3GroupUsersFilePath() string {
	return fmt.Sprintf("%s%s", c.S3BaseDir, c.GroupUsersFilename)
}
