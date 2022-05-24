# Elasticsearch with ktranslate
This is a Docker compose setup that will create an Elasticsearch with Kibana stack and setup
ktranslate to listen to a NetFlow generator.

# Run

`docker compose up`

This should launch all services. To view Kibana, visit http://localhost:5601. You will need to go to
http://localhost:5601/app/management/kibana/indexPatterns/create and create an index pattern of
`kentik*`. You can then visit http://localhost:5601/app/discover and you should see data.
