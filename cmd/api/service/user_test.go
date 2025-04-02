package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/iamsorryprincess/go-project-layout/cmd/api/model"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/configutils"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type testConfig struct {
	RedisConfig redis.Config `mapstructure:"redis"`
}

func TestFillRedisUsers(t *testing.T) {
	cfg, err := configutils.Parse[testConfig](nil, "..")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := redis.NewConnection(log.New("debug", "test"), cfg.RedisConfig)
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	count := 100000

	users := make([]interface{}, count)
	for i := 0; i < count; i++ {
		user := model.User{
			Name: fmt.Sprintf("user-%d", i),
		}
		data, cErr := json.Marshal(user)
		if cErr != nil {
			t.Fatal(cErr)
		}
		users[i] = data
	}

	if err = conn.RPush(context.Background(), "test", users...).Err(); err != nil {
		t.Fatal(err)
	}
}
