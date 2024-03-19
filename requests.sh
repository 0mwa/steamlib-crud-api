#curl -X GET http://localhost:3333/games/440

curl -X POST http://localhost:3333/games/add/440 --data '{
                                                           "name": "TF2",
                                                           "img": "img.png",
                                                           "description": "cool",
                                                           "rating": "10",
                                                           "developer_id": "1",
                                                           "publisher_id": "15",
                                                           "steam_id": "440"
                                                         }'