package internal

var SelectGames = `
			SELECT 
			    COALESCE(name,'NULL'),
			    COALESCE(img,'NULL'),
			    COALESCE(description,'NULL'),
			    COALESCE(rating,0),
			    COALESCE(developer_id,0),
			    COALESCE(publisher_id,0),
			    COALESCE(steam_id,0)
			FROM games
		`
var SelectGamesSortName = `
			SELECT 
			    COALESCE(name,'NULL'),
			    COALESCE(img,'NULL'),
			    COALESCE(rating,0),
			    COALESCE(description,'NULL')
			FROM games
			ORDER BY name
		`
var SelectGameById = `
			SELECT
			    COALESCE(name,'NULL'),
			    COALESCE(img,'NULL'),
			    COALESCE(description,'NULL'),
			    COALESCE(rating,0),
			    COALESCE(developer_id,0),
			    COALESCE(publisher_id,0),
			    COALESCE(steam_id,0)
			FROM games
			WHERE steam_id = $1
		`
var DeleteGameById = `
			DELETE
			FROM games
			WHERE steam_id = $1
		`
var UpdateGameById = `
			UPDATE games
			SET 
			    name = $1,
			    img = $2,
			    description = $3,
			    rating = $4,
			    developer_id = $5,
			    publisher_id = $6
			WHERE steam_id = $7
		`
var SelectPublishers = `
			SELECT
				COALESCE(name,'NULL'),
				COALESCE(country,'NULL')
			FROM publishers
		`
var SelectPublisherById = `
			SELECT
				COALESCE(name,'NULL'),
				COALESCE(country,'NULL')
			FROM publishers
			WHERE steam_id = $1
		`
var DeletePublisherById = `
			DELETE
			FROM publishers
			WHERE steam_id = $1
		`
var UpdatePublisherById = `
			UPDATE publishers
			SET
			    name = $1,
				country = $2
			WHERE steam_id = $3
		`
