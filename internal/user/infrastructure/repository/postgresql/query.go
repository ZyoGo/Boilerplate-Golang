package postgresql

var (
	queryInsertUser = `
		INSERT INTO user_account (
			id,
			email,
			password,
			created_at
		)
		VALUES (
			@id,
			@email,
			@password,
			@created_at
		)
	`

	queryFindUserByEmail = `
		SELECT
			email
		FROM
			user_account
		WHERE
			email = @email
		AND
			deleted_at IS NULL
	`

	queryUpdateUser = `
		UPDATE
			user_account
		SET
			password = @password
		WHERE
			id = @id
	`
)
