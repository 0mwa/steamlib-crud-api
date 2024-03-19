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
