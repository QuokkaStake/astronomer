package database

func (d *Database) GetAllChainBinds(chatID string) ([]string, error) {
	chains := make([]string, 0)

	rows, err := d.client.Query(
		"SELECT chain FROM chain_binds WHERE chat_id = $1",
		chatID,
	)
	if err != nil {
		d.logger.Error().Str("chat", chatID).Err(err).Msg("Error getting chain binds")
		return chains, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()

	for rows.Next() {
		var chain string

		err = rows.Scan(&chain)
		if err != nil {
			d.logger.Error().Str("chat", chatID).Err(err).Msg("Error getting chain bind")
			return chains, err
		}

		chains = append(chains, chain)
	}

	return chains, nil
}

func (d *Database) InsertChainBind(
	reporter string,
	chatID string,
	chatName string,
	chain string,
) error {
	_, err := d.client.Exec(
		"INSERT INTO chain_binds (reporter, chat_id, chat_name, chain) VALUES ($1, $2, $3, $4)",
		reporter,
		chatID,
		chatName,
		chain,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert chain bind")
		return err
	}

	return nil
}

func (d *Database) DeleteChainBind(
	reporter string,
	chatID string,
	chain string,
) (bool, error) {
	result, err := d.client.Exec(
		"DELETE FROM chain_binds WHERE reporter = $1 AND chat_id = $2 AND chain = $3",
		reporter,
		chatID,
		chain,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not delete chain bind")
		return false, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
