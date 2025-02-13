const axios = require('axios')

module.exports = {
    async post(method, data, headers) {
        try {
            const response = await axios.post(`http://127.0.0.1:8080/${method}`, data || {}, headers || {});

            return response.data;
        } catch (error) {
            // console.error('Ошибка при выполнении POST запроса:', error);
            throw error;
        }
    },

    async delete(method, data) {
        try {
            const response = await axios.delete(`http://127.0.0.1:8080/${method}`, data || {});

            return response;
        } catch (error) {
            // console.error('Ошибка при выполнении DELETE запроса:', error);
            throw error;
        }
    },

    async get(method, data) {
        try {
            const response = await axios.get(`http://127.0.0.1:8080/${method}`, data || {});

            return response.data;
        } catch (error) {
            // console.error('Ошибка при выполнении GET запроса:', error);
            throw error;
        }
    },

    async put(method, data, headers) {
        try {
            const response = await axios.put(`http://127.0.0.1:8080/${method}`, data || {}, headers || {});

            return response.data;
        } catch (error) {
            // console.error('Ошибка при выполнении GET запроса:', error);
            throw error;
        }
    }
}