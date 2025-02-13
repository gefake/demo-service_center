const
	website = require('./core/website')

const port = 3000

website.listen(port, function() {
    console.log(`App started on port ${port}`)  
})