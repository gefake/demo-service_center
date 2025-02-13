module.exports = {
    Authorization(req, res, next) {
        const user = req.session.user
        if (!user) return res.render('message', { text: 'Чтобы продолжить, необходимо войти в аккаунт', error: '403', type: 'danger', btext: "Вход", back: '/user/login', error: '403' })
        const token = user.token
        if (!token) return res.render('message', { text: 'Чтобы продолжить, необходимо войти в аккаунт', error: '403', type: 'danger', btext: "Вход", back: '/user/login', error: '403' })
        next()
    },
}