package interfaces

type DatabaseSource interface {
	FindAll(fieldName, fieldValue string) ([]interface{}, error)
	FindOne(key, value string) (interface{}, error)
	Save(data interface{}) (id interface{}, err error)
	Update(identifier string, data interface{}) (id interface{}, err error)
	Delete(identifier string) (id interface{}, err error)
}
