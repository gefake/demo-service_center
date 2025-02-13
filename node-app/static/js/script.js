let form = document.getElementById('request_form')

if (form) {
    form.addEventListener('submit', async function (event) {
        event.preventDefault();
        const formData = new FormData(this);
        const response = await fetch('/submit', {
            method: 'POST',
            body: formData
        });
        const data = await response.json();
        console.log(data);
    });
}