# graphql-dal

http://localhost:8080/graphql?query={product(id:%2290804%22){id}}
http://localhost:8080/graphql?query=mutation{product(id:%2290804%22,variant:{status:5,upc:%22blah%22,externalId:%222%22,images:[%22234%22]},topCategory:10){id,variant{status+upc+externalId+images},topCategory}}
