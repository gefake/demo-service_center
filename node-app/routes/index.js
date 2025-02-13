const
    express = require('express'),
    rest_api = require('./_rest-api'),
    _ = require('lodash'),
    tabulate = require('tabulate'),
    mw = require('./_mw'),
    axios = require('axios')
const moment = require('moment')

const router = new express.Router()

router.use('/user', require('./user'))

router.get('/', (req, res) => {
    const user = req.session.user

    let date = new Date;
    let hour = date.getHours()
    // астраханское время, сервер в Москве
    hour = hour + 1

    // hour = 5

    let traffic_lights_stage

    if (hour >= 7 && hour < 22) {
        traffic_lights_stage = 1
    } else if (hour >= 22 || hour < 4) {
        traffic_lights_stage = 2
    } else {
        traffic_lights_stage = 3
    }

    rest_api.get('api/traffic-lights', {
        headers: {
            // Authorization: `Bearer ${user.token}`
        },
    })
        .then(data => {
            console.log(data)
            if (!data) return;
            res.render('main_old', { title: "Главная", user: user, traffic_lights_stage: data })
        })
        .catch(error => {
            res.render('main_old', { title: "Главная", user: user, traffic_lights_stage: traffic_lights_stage })
        });

})

router.get('/old', (req, res) => {
    const user = req.session.user

    res.render('main', { title: "Главная", user: user })
})

router.get('/mywot1d0b059f1b895a831f51.html', (req, res) => {
    res.sendFile(`${__dirname}/mywot1d0b059f1b895a831f51.html`);
})

router.get('/robots.txt', (req, res) => {
    res.sendFile(`${__dirname}/robots.txt`);
})

router.get('/sitemap.xml', (req, res) => {
    res.sendFile(`${__dirname}/sitemap.xml`);
})

router.post('/submit', (req, res) => {
    const user = req.session.user

    const name = req.body.name;
    let tel = req.body.tel;

    tel = tel.replace(/_/g, '');
    tel = tel.replace(/-/g, '');
    tel = tel.replace(/\(/g, '');
    tel = tel.replace(/\)/g, '');

    // if (req.session.submittedTask) {
    //   req.session.submittedTask = false;
    //   res.json({ success: true, message: 'Заявка успешно отправлена!' });
    //   return;
    // }

    // req.session.submittedTask = true;

    // res.render('main_old', { title: "Главная", user: user, submitted: true, traffic_lights_stage: 4 })

    rest_api.post(`api/task/`, {
      name: name,
      phoneNumber: tel
    })
      .then(data => {
        console.log(data);
        res.setHeader('Content-Type', 'application/json');
        res.json({ success: true, message: 'Заявка успешно отправлена!' });
        // Дальнейшая работа с данными
      })
      .catch(error => {
        res.setHeader('Content-Type', 'application/json');
        res.json({ success: false, message: 'Заявка не была отправлена!' });
        console.error('Ошибка запроса:', error);
        // Обработка ошибки здесь
      });
  });

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MTg1Mjk2MzEsImlhdCI6MTcxODQ0MzIzMX0.QKNXPGNgNzQCVT8fus7jKntttqqVIApSk3lyxSCL_uI"
const config = {
    headers: { Authorization: `Bearer ${token}` }
};

router.get('/admin', mw.Authorization, async(req, res) => {
    // const user = req.session.user
    res.redirect('/admin/requests')
})

function sliceObject(obj, start, end) {
    // Создаем новый пустой для хранения результата
    const result = {};

    // Перебираем все ключи исходного объекта
    for (let key in obj) {
        // Преобразуем ключ в число, если это возможно
        let index = parseInt(key, 10);

        if (!isNaN(index)) {
        // Проверяем, находится ли ключ в пределах указанного диапазона
        if (index >= start && index <= end) {
            // Добавляем ключ и значение в результат
            result[index] = obj[key];
        }
        }
    }

    return result;
}

router.get('/admin/requests', (req, res) => {res.redirect('/admin/requests/0-20')})

router.get('/admin/requests/:min-:max', mw.Authorization, async(req, res) => {
    const user = req.session.user

    let min = req.params.min || 0
    let max = req.params.max || 20

    rest_api.get('api/admin/task-manage', {
        headers: {
            Authorization: `Bearer ${user.token}`
        },
    })
        .then(data => {
            if (!data) return;

            for(var obj in data) {
                if (data[obj].date) {
                    let dateObj = moment.unix(data[obj].date)
                    data[obj].date = dateObj.locale('ru').format('DD MMMM YYYY, в HH:mm')
                }
            }

            data.sort((a,b) => b.id - a.id);

            let slicedData = sliceObject(data, min, max)

            res.render('admin/requests', { title: "Админ-панель", tasks: slicedData, dataCount: Object.keys(data).length, activeTab: "Вызовы" })
        })
        .catch(error => {
            res.render('message', { text: error, type: 'danger', back: '/user/login', title:"Error", btext:'Продолжить' })
        });
})

router.get('/admin/telegram', mw.Authorization, async(req, res) => {
    const user = req.session.user

    rest_api.get('api/admin/telegram-trust', {
        headers: {
            Authorization: `Bearer ${user.token}`
        },
    })
        .then(data => {
            if (!data) return;

            console.log(data)

            data.sort((a,b) => b.id - a.id);

            res.render('admin/telegram', { title: "Админ-панель", users: data, activeTab: 'Telegram' })
        })
        .catch(error => {
            res.render('message', { text: error, type: 'danger', back: '/user/login', title:"Error", btext:'Продолжить' })
        });
})

router.get('/admin/settings', mw.Authorization, async(req, res) => {
    const user = req.session.user

    rest_api.get('api/traffic-lights', {
        // headers: {
        //     Authorization: `Bearer ${user.token}`
        // },
    })
        .then(data => {
            if (!data) return;

            res.render('admin/settings', { title: "Админ-панель", users: data, activeTab: 'Настройки', traffic_lights: data })
        })
        .catch(error => {
            res.render('message', { text: error, type: 'danger', back: '/user/login', title:"Error", btext:'Продолжить' })
        });
})

router.post('/admin/traffic-lights/:id', mw.Authorization, (req, res) => {
    const user = req.session.user
    const name = req.body.name;
    const id = req.params.id;
    if (!id) {
        return res.render('message', { text: 'Неизвестный ID для светофора', type: 'danger', back: '/admin', title:"Ошибка", btext:'Назад' })
    };

    rest_api.put(`api/admin/traffic-lights/${id}`, {}, {
        headers: {
            Authorization: `Bearer ${user.token}`
        },
    })
        .then(data => {
            console.log(data);
            res.redirect('/admin/settings')
        })
        .catch(error => {
            console.error('Ошибка запроса:', error);
        });
});

router.post('/admin/telegram-trust', mw.Authorization, (req, res) => {
    const user = req.session.user
    const username = req.body.username;
    if (!username || username == '') {
        return res.render('message', { text: 'Неизвестный username пользователя', type: 'danger', back: '/admin', title:"Ошибка", btext:'Назад' })
    };

    console.log(username)

    rest_api.post(`api/admin/telegram-trust/${username}`, {}, {
        headers: {
            Authorization: `Bearer ${user.token}`
        }
    })
        .then(data => {
            res.redirect('/admin/telegram')
        })
        .catch(error => {
            console.log(error.response.data)

            if (error.response && error.response.data && error.response.data.message) {
                res.render('message', { text: error.response.data.message, type: 'danger', back: '/admin/telegram', title:"Ошибка", btext:'Назад' })
            } else {
                res.render('message', { text: error, type: 'danger', back: '/admin/telegram', title:"Ошибка", btext:'Назад' })
            }
        });
});

router.post('/admin/telegram-trust/delete/:id', mw.Authorization, (req, res) => {
    const user = req.session.user
    const name = req.body.name;
    const id = req.params.id;
    if (!id) {
        return res.render('message', { text: 'Неизвестный ID пользователя', type: 'danger', back: '/admin', title:"Ошибка", btext:'Назад' })
    };

    rest_api.delete(`api/admin/telegram-trust/${id}`, {
        headers: {
            Authorization: `Bearer ${user.token}`
        },
    })
        .then(data => {
            console.log(data);
            res.redirect('/admin/telegram')
        })
        .catch(error => {
            console.error('Ошибка запроса:', error);
        });
});

router.post('/admin/delete/:id', mw.Authorization, (req, res) => {
    const user = req.session.user
    const name = req.body.name;
    const id = req.params.id;
    if (!id) {
        return res.render('message', { text: 'Неизвестный ID заявки', type: 'danger', back: '/admin', title:"Ошибка", btext:'Назад' })
    };

    rest_api.delete(`api/admin/task-manage/${id}`, {
        headers: {
            Authorization: `Bearer ${user.token}`
        },
    })
        .then(data => {
            console.log(data);
            res.redirect('/admin')
        })
        .catch(error => {
            console.error('Ошибка запроса:', error);
        });
});

module.exports = router
