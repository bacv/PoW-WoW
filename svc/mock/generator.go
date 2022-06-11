package mock

type MockGenerator struct{}

var MockID = "testid"

func (g *MockGenerator) GenID() string {
	return MockID
}
