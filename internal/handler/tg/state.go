package tg

import (
	"tmballNews/internal/domain"
)

func (a *API) setUserState(userID string, state domain.UserState) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.userStates[userID] = state
}

func (a *API) getUserState(userID string) domain.UserState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.userStates[userID]
}

func (a *API) clearUserState(userID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.userStates, userID)
}
