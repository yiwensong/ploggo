package postgres

import (
	context "context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	testing "testing"

	errors "github.com/pkg/errors"
	testcontainers "github.com/testcontainers/testcontainers-go"
	wait "github.com/testcontainers/testcontainers-go/wait"
)

type Logger struct{}

func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

var testConnectionString string

var testPassword = "testpassword"
var testUser = "testuser"
var testDb = "testdb"

func SetupPostgres() (
	connectionString string,
	cleanup func(),
	err error,
) {
	ctx := context.Background()

	// Register the image
	cmd := exec.Command("../../../containers/postgres.executable")
	output, err := cmd.Output()
	if err != nil {
		return "", func() {}, errors.Wrapf(err, "cmd.Output(%s)", cmd.String())
	}
	fmt.Println(string(output))

	req := testcontainers.ContainerRequest{
		Image:        "bazel/containers:postgres",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": testPassword,
			"POSTGRES_USER":     testUser,
			"POSTGRES_DB":       testDb,
		},
		WaitingFor: wait.ForExposedPort(),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", func() {}, errors.Wrapf(err, "testcontainers.GenericContainer")
	}

	cleanup = func() {
		if err := postgresC.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate container: %s\n", err.Error())
			os.Exit(1)
		}
	}

	endpoint, err := postgresC.Endpoint(ctx, "")
	if err != nil {
		return "", func() {}, errors.Wrapf(err, "postgresC.Endpoint()")
	}

	connectionString = fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		testUser,
		testPassword,
		endpoint,
		testDb,
	)

	storage, cleanup, err := NewAvalonPostgresStorage(ctx, connectionString)
	defer cleanup()

	err = storage.RunMigrations(ctx, "./schema")
	if err != nil {
		return "", func() {}, errors.Wrapf(err, "RunMigrations(./schema)")
	}

	return connectionString, cleanup, nil
}

func TestMain(m *testing.M) {
	// Using glog, make sure we have flag.Parse
	flag.Parse()

	// Set up the database and set the connection string for the test
	connectionString, cleanup, err := SetupPostgres()
	if err != nil {
		fmt.Printf("Postgres setup failed: %s\n", err.Error())
		os.Exit(1)
	}
	defer cleanup()
	testConnectionString = connectionString

	os.Exit(m.Run())
}
