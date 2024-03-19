package internal

var SelectGames = `
			SELECT 
			    COALESCE(name,'NULL'),
			    COALESCE(img,'NULL'),
			    COALESCE(rating,0),
			    COALESCE(description,'NULL')
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
			    COALESCE(rating,0),
			    COALESCE(description,'NULL')
			FROM games
			WHERE steam_id = $1
`
