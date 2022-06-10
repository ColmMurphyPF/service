package dbtest

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/colmmurphy91/go-service/business/sys/auth"
	"github.com/colmmurphy91/go-service/business/sys/database/mongodb"
	"github.com/colmmurphy91/go-service/foundation/docker"
	"github.com/colmmurphy91/go-service/foundation/keystore"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func StartMongo() (*docker.Container, error) {
	image := "mongo:latest"
	port := "27017"
	args := []string{"-e", "MONGO_INITDB_ROOT_USERNAME=admin", "-e", "MONGO_INITDB_ROOT_PASSWORD=admin"}
	fmt.Println("starting mongo")
	return docker.StartContainer(image, port, args...)
}

func NewMongoUnit(t *testing.T, c *docker.Container, collectionName string) (*zap.SugaredLogger, *mongo.Database, func()) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbM, err := mongodb.Open(mongodb.Config{
		User:     "admin",
		Password: "admin",
		DBName:   "test",
		Host:     c.Host,
	})

	if err != nil {
		fmt.Println("here")
	}

	var buf bytes.Buffer
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buf)
	log := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)).
		Sugar()

	teardown := func() {
		t.Helper()

		log.Sync()

		writer.Flush()
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return log, dbM, teardown
}

type MongoTest struct {
	DB       *mongo.Database
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	Teardown func()

	t *testing.T
}

func NewMongoIntegration(t *testing.T, c *docker.Container, dbName string) *MongoTest {
	log, db, teardown := NewMongoUnit(t, c, dbName)

	keyID := "4754d86b-7a6d-4df5-9c65-224741361492"
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Build an authenticator using this private key and id for the key store.
	auth1, err := auth.New(keyID, keystore.NewMap(map[string]*rsa.PrivateKey{keyID: privateKey}))

	if err != nil {
		t.Fatal(err)
	}

	test := MongoTest{
		DB:       db,
		Log:      log,
		Auth:     auth1,
		t:        t,
		Teardown: teardown,
	}

	return &test
}

func (t *MongoTest) AdminToken() string {
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: pq.StringArray{"admin"},
	}

	token, err := t.Auth.GenerateToken(claims)
	if err != nil {
		t.t.Fatal(err)
	}

	return token
}

func (t *MongoTest) OwnerToken(userID string) string {
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: pq.StringArray{"admin"},
	}

	token, err := t.Auth.GenerateToken(claims)
	if err != nil {
		t.t.Fatal(err)
	}

	return token
}
