package config

const (
	envPort                 = "PORT"
	envDatabaseEndpoint     = "DATABASE_HOST"
	envDatabaseKind         = "DATABASE_KIND"
	envDatabaseLogLevel     = "DATABASE_LOG_LEVEL"
	envDatabaseMaxIdleConns = "DATABASE_MAX_IDLE_CONNS"
	envDatabaseMaxOpenConns = "DATABASE_MAX_OPEN_CONNS"
	envStage                = "STAGE"
	envJwtSecret            = "JWT_SECRET"
	envNearEnv              = "NEAR_ENV"
	envNearID               = "NEAR_ID"
	envNearPrivateKey       = "NEAR_PRIVATE_KEY"
	envNearFtContract       = "NEAR_FT_CONTRACT"
	envNearFtStorageDeposit = "NEAR_FT_STORAGE_DEPOSIT"
)

type Stage string

const (
	LOCAL Stage = "local"
	DEV   Stage = "development"
	STAG  Stage = "staging"
	PROD  Stage = "production"
)

type config struct {
	port                 string
	database             []database
	stage                Stage
	jwtSecret            string
	aws                  aws
	nearEnv              string
	nearID               string
	nearPrivateKey       string
	nearFtContract       string
	nearFtStorageDeposit string
}

func (c config) GetPort() string {
	return c.port
}

func (c config) GetDatabases() []database {
	return c.database
}

func (c config) GetStage() Stage {
	return c.stage
}

func (c config) GetJWTSecret() string {
	return c.jwtSecret
}

func (c config) GetAws() aws {
	return c.aws
}

func (c config) GetNearEnv() string {
	return c.nearEnv
}

func (c config) GetNearID() string {
	return c.nearID
}

func (c config) GetNearPrivateKey() string {
	return c.nearPrivateKey
}

func (c config) GetNearFtContract() string {
	return c.nearFtContract
}

func (c config) GetNearFtStorageDeposit() string {
	return c.nearFtStorageDeposit
}

// database configuration
type database struct {
	kind         int
	endpoint     string // database endpoint which is in one line. e.g., {user}:{password}@({uri})/{scheme_name}?charset=utf8mb4&parseTime=True&loc=Europe%2FLondon&collation=utf8mb4_unicode_ci
	logLevel     int    // 1. silent 2. warn 3. error 4. info
	maxIdleConns int
	maxOpenConns int
}

func (d database) GetEndpoint() string {
	return d.endpoint
}

func (d database) GetKind() int {
	return d.kind
}

func (d database) GetLogLevel() int {
	return d.logLevel
}

func (d database) GetMaxIdleConns() int {
	if d.maxIdleConns == 0 {
		d.maxIdleConns = 10
	}
	return d.maxIdleConns
}

func (d database) GetMaxOpenConns() int {
	if d.maxOpenConns == 0 {
		d.maxOpenConns = 100
	}
	return d.maxOpenConns
}

type aws struct {
	region string
	ses    ses
}

func (a aws) GetRegion() string {
	return a.region
}

func (a aws) GetSesAccessID() string {
	return a.ses.accessID
}

func (a aws) GetSesSecretKey() string {
	return a.ses.secretKey
}

type ses struct {
	accessID, secretKey string
}
