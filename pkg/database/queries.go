package database

import "main/pkg/types"

func (d *Database) InsertQuery(query *types.Query) error {
	_, err := d.client.Exec(
		"INSERT INTO queries (reporter, user_id, user_name, chat_id, command, query) VALUES ($1, $2, $3, $4, $5, $6)",
		query.Reporter,
		query.UserID,
		query.Username,
		query.ChatID,
		query.Command,
		query.Query,
	)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not insert query")
		return err
	}

	return nil
}
