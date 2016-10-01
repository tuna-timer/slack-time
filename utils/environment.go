package utils

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/olebedev/config"
)

const (
	// ProductionEnv - a value that indicates about production env
	ProductionEnv = "production"
	// DevelopmentEnv - a value that indicates about development env
	DevelopmentEnv = "development"
	// TestEnv - a value that indicates about test env
	TestEnv = "test"
)

const (
	// ConfigFile - path to YML config file
	ConfigFile string = "config.yml"

	// MigrationsFolder - the folder to look migration SQLs in
	MigrationsFolder string = "models/migrations/"

	// gormLogSQL - Whether GORM SQL logging is enabled or not
	gormLogSQL bool = false
)

// Environment is a thing that holds env. specific stuff
type Environment struct {
	Config     *config.Config
	AppVersion string
	Name       string
	CreatedAt  time.Time
}

// NewEnvironment creates a new environment
func NewEnvironment(environment string, appVersion string) *Environment {
	cfg, err := readConfig(environment)
	if err != nil {
		log.Fatal(err) //no way to launch the app without an Environment, fatal!
	}

	cfg, err = cfg.Get(environment)
	cfg.Env()

	return &Environment{Name: environment, AppVersion: appVersion, CreatedAt: time.Now(), Config: cfg}
}

// MigrateDatabase - performs database migrations
func (env *Environment) MigrateDatabase(session *mgo.Session) error {
	log.Println("Migrating database...")

	teams := session.DB("").C("teams")

	teams.Create(&mgo.CollectionInfo{})
	teams.EnsureIndex(mgo.Index{
		Unique: true,
		Key:    []string{"ext_id"},
	})

	timers := session.DB("").C("timers")
	timers.Create(&mgo.CollectionInfo{})
	timers.EnsureIndex(mgo.Index{Key: []string{"team_id"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"project_id"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"team_user_id"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"hash"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"created_at"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"created_at"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"finished_at"}})
	timers.EnsureIndex(mgo.Index{Key: []string{"deleted_at"}})

	users := session.DB("").C("team_users")
	users.Create(&mgo.CollectionInfo{})
	users.EnsureIndex(mgo.Index{
		Unique: true,
		Key:    []string{"ext_id"},
	})

	log.Println("Database migrated!")
	return nil
}

// ConnectToDatabase todo
func ConnectToDatabase(cfg *config.Config) (*mgo.Session, error) {
	session, err := mgo.Dial(cfg.UString("database.url"))
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func readConfig(environmentName string) (*config.Config, error) {
	cfg, err := config.ParseYamlFile(adjustPath(environmentName, ConfigFile))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// a hack to walk around this issue:
// http://stackoverflow.com/questions/23847003/golang-tests-and-working-directory
// does it have a nicer solution?
func adjustPath(environmentName string, resource string) string {
	if environmentName == TestEnv {
		return "../" + resource
	}
	return resource
}
