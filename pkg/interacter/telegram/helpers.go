package telegram

func (interacter *Interacter) ChainNotFound() (string, error) {
	chains, err := interacter.Database.GetAllChains()
	if err != nil {
		return "Could not get chains list!", err
	}
	return interacter.TemplateManager.Render("chain_not_found", chains)
}
