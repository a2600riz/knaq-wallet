package groups

import "knaq-wallet/controller/handlers/wallet"

func registerHandler() {
	secureRegister(wallet.SecureWalletHandler{})
	internalRegister(wallet.InternalWalletHandler{})
}

func globalRegister(handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(globalPath)
	}
}
func secureRegister(handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(securePath)
	}
}
func internalRegister(handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(internalPath)
	}
}
func adminRegister(handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(adminPath)
	}
}
