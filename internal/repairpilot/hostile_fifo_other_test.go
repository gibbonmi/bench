//go:build !linux && !darwin

package repairpilot

func platformStoredHostileFixtures() []storedHostileFixture {
	return nil
}
