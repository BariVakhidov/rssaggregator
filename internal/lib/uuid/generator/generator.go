package uuidgenerator

import "github.com/google/uuid"

type RealUUIDGenerator struct{}

func New() *RealUUIDGenerator {
	return &RealUUIDGenerator{}
}

func (r *RealUUIDGenerator) Generate() uuid.UUID {
	return uuid.New()
}

func (r *RealUUIDGenerator) NilUUID() uuid.UUID {
	return uuid.Nil
}
