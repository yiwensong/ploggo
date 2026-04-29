package postgres

import (
	context "context"
	flag "flag"
	fmt "fmt"
	os "os"
	exec "os/exec"
	filepath "path/filepath"
	strings "strings"
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
	// rules_oci generates a load script. We try to find it in runfiles.
	wd, _ := os.Getwd()
	parent := filepath.Dir(wd)
	manifestPath := filepath.Join(parent, "MANIFEST")
	
	runfilesMap := make(map[string]string)
	if _, err := os.Stat(manifestPath); err == nil {
		data, _ := os.ReadFile(manifestPath)
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				runfilesMap[parts[0]] = strings.TrimSpace(parts[1])
			}
		}
	}

	resolvePath := func(path string) string {
		if p, ok := runfilesMap["_main/"+path]; ok {
			return p
		}
		if p, ok := runfilesMap[path]; ok {
			return p
		}
		return path
	}

	possiblePaths := []string{
		resolvePath("containers/wrap_postgres.bat"),
		resolvePath("containers/postgres.bat"),
		resolvePath("containers/postgres.sh"),
		resolvePath("containers/postgres"),
	}

	var loadCmd *exec.Cmd
	for _, path := range possiblePaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			if strings.HasSuffix(path, ".sh") && os.PathSeparator == '\\' {
				loadCmd = exec.Command("bash", path)
			} else {
				loadCmd = exec.Command(path)
			}
			break
		}
	}

	if loadCmd != nil {
		// Attempt to load the image, but don't fail hard if Docker is missing
		// since we know it might not be available in this environment.
		output, _ := loadCmd.CombinedOutput()
		fmt.Printf("Load attempt output: %s\n", string(output))
	}

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
