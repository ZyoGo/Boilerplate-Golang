package postgresql

var (
	queryInsertUser = `
		INSERT INTO user (
			id,
			email,
			password
		)
		VALUES (
			@id,
			@email,
			@password,
		)
	`

	queryFindUserByEmail = `
		SELECT
			email
		FROM
			user
		WHERE
			email = @email
	`

	queryUpdateUser = `
		UPDATE
			users
		SET
			password = @password
		WHERE
			id = @id
	`
)
