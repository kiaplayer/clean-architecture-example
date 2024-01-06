package generators

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kiaplayer/clean-architecture-example/internal/domain/entity/reference"
)

type NumberGenerator struct{}

func NewNumberGenerator() *NumberGenerator {
	return &NumberGenerator{}
}

func (n *NumberGenerator) GenerateNumber(date time.Time, company reference.Company) string {
	return fmt.Sprintf("%s-%d-%d", date.Format("20060102"), company.ID, rand.Intn(10000))
}
