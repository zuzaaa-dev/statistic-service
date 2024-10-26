package database

type DatabaseConnectable interface {
	Connect() (interface{}, error)
	CloseConnection(interface{}) error
}
