const
    express = require('express'),
    rest_api = require('./_rest-api'),
    mw = require('./_mw')

const router = new express.Router()

router.post('/login', async(req, res) => {
    var user = req.body.username
    var password = req.body.password

    console.log(user, password)

    const body = {
        name: user,
        password: password
    }

    rest_api.post('auth/admin/sign-in', body)
        .then(data => {
            console.log(`Token: ${data.token}`);

            req.session.user = {
                token: data.token
            }

            res.redirect('/admin/requests')
        })
        .catch(error => {
            res.redirect('/user/login')
            console.error('Ошибка причении данных:', error);
        });
})

router.get('/login', (req, res) => {
    res.render('login', {title: "Вход"})
})

module.exports = router