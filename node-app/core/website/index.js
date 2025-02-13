'use strict';

const
    express = require('express'),
    session = require('express-session'),
    bodyParser = require('body-parser'),
    routes = require('../../routes')

const app = express()
const sessionMW = session({
    secret: 'testsecret',
    resave: false,
    saveUninitialized: true,
    cookie: {
        secure: false,
        maxAge: 7 * 24 * 60 * 1000,
    },
})

app.set('view engine', 'pug')
app.use(sessionMW)
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }))
app.use(routes)
app.use('/static', express.static(__dirname + '/../../static'))

module.exports = app