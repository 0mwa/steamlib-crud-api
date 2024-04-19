package repository

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
		`
var SelectPublishers = `
			SELECT
				COALESCE(name,'NULL'),
				COALESCE(country,'NULL'),
				COALESCE(steam_id,0)
			FROM publishers
		`
var SelectPublisherById = `
			SELECT
				COALESCE(name,'NULL'),
				COALESCE(country,'NULL'),
				COALESCE(steam_id,0)
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
		`
var SelectDevelopers = `
			SELECT
				COALESCE(name,'NULL'),
				COALESCE(country,'NULL'),
				COALESCE(steam_id,0)
			FROM developers
		`
var SelectDeveloperById = `
			SELECT
				COALESCE(name,'NULL'),
				COALESCE(country,'NULL'),
				COALESCE(steam_id,0)
			FROM developers
			WHERE steam_id = $1
		`
var DeleteDeveloperById = `
			DELETE
			FROM developers
			WHERE steam_id = $1
		`
var UpdateDeveloperById = `
			UPDATE developers
			SET
		`
var GetGamesCount = `
			SELECT COUNT (*)
			FROM games
	`
var GetDevelopersCount = `
			SELECT COUNT (*)
			FROM developers
	`
var GetPublishersCount = `
			SELECT COUNT (*)
			FROM publishers
	`
