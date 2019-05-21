package agents

func NewDSNPreparer() *DSNPreparer{
	return &DSNPreparer{}
}

type DSNPreparer struct {}

func (p *DSNPreparer) MySQLDSN(host string, port uint16, username, password string) string {
	return mysqlDSN(host, port, username, password)
}

func (p *DSNPreparer) PostgreSQLDSN(host string, port uint16, username, password string) string {
	return postgresqlDSN(host, port, username, password)
}

func (p *DSNPreparer) MongoDBDSN(host string, port uint16, username, password string) string {
	return mongoDSN(host, port, username, password)
}