package schema

type Serializable interface {
	Serialize() string
}
