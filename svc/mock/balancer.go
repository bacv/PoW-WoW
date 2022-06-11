package mock

type MockBalancer struct{}

func (b *MockBalancer) GetChallengeBits(load int) uint {
	return 0
}
